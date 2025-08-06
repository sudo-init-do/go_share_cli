package server

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
)

// FileInfo represents a file or directory for template rendering
type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
	IsDir   bool
	Icon    string
	SizeStr string
}

// API response types for React frontend
type APIFileItem struct {
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Size          int64     `json:"size"`
	IsDir         bool      `json:"isDir"`
	ModTime       time.Time `json:"modTime"`
	DownloadCount int       `json:"downloadCount"`
}

type APIPageData struct {
	Title       string        `json:"title"`
	CurrentPath string        `json:"currentPath"`
	ParentPath  string        `json:"parentPath"`
	Files       []APIFileItem `json:"files"`
	HasParent   bool          `json:"hasParent"`
	ServerURL   string        `json:"serverURL"`
}

// PageData contains data for the HTML template
type PageData struct {
	Title       string
	CurrentPath string
	ParentPath  string
	Files       []FileInfo
	HasParent   bool
	ServerURL   string
	QRCodeData  string
	HasAuth     bool
}

// FileStats tracks download counts and access logs
type FileStats struct {
	DownloadCount int       `json:"download_count"`
	LastAccessed  time.Time `json:"last_accessed"`
}

var (
	fileStatsMap = make(map[string]*FileStats)
	statsMapLock sync.RWMutex
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <script>
        // Theme Toggle
        function toggleTheme() {
            const body = document.body;
            const isDark = body.classList.contains('dark');
            if (isDark) {
                body.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            } else {
                body.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            }
        }
        
        // Load theme on page load
        document.addEventListener('DOMContentLoaded', function() {
            const theme = localStorage.getItem('theme') || 'light';
            if (theme === 'dark') {
                document.body.classList.add('dark');
            }
        });

        // File Search Functionality
        function searchFiles() {
            const searchTerm = document.getElementById('fileSearch').value.toLowerCase();
            const rows = document.querySelectorAll('#fileTable tbody tr');
            
            rows.forEach(row => {
                const fileName = row.querySelector('td:first-child')?.textContent?.toLowerCase() || '';
                if (fileName.includes(searchTerm) || searchTerm === '') {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            });
        }

        // Toggle QR Code visibility
        function toggleQR() {
            const qrSection = document.getElementById('qrSection');
            qrSection.classList.toggle('hidden');
        }
    </script>
    <style>
        .dark {
            --tw-bg-opacity: 1;
            background-color: rgb(17 24 39 / var(--tw-bg-opacity));
            color: white;
        }
        .dark .bg-gray-50 { background-color: rgb(17 24 39); }
        .dark .bg-white { background-color: rgb(31 41 55); }
        .dark .bg-gray-100 { background-color: rgb(55 65 81); }
        .dark .text-gray-800 { color: rgb(229 231 235); }
        .dark .text-gray-600 { color: rgb(156 163 175); }
        .dark .text-gray-500 { color: rgb(107 114 128); }
        .dark .border-gray-200 { border-color: rgb(55 65 81); }
    </style>
</head>
<body class="bg-gray-50 min-h-screen transition-colors duration-200">
    <div class="container mx-auto px-4 py-8 max-w-6xl">
        <header class="mb-8">
            <div class="flex items-center justify-between mb-4">
                <h1 class="text-3xl font-bold text-gray-800">
                    <i class="fas fa-share-alt text-blue-600 mr-2"></i>
                    GoShare File Browser
                </h1>
                <div class="flex items-center space-x-4">
                    <button onclick="toggleQR()" class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
                        <i class="fas fa-qrcode mr-2"></i>
                        QR Code
                    </button>
                    <button onclick="toggleTheme()" class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
                        <i class="fas fa-moon mr-2"></i>
                        Theme
                    </button>
                </div>
            </div>
            <p class="text-gray-600 mb-4">Current directory: <code class="bg-gray-200 px-2 py-1 rounded">{{.CurrentPath}}</code></p>
            
            <!-- QR Code Section -->
            <div id="qrSection" class="hidden bg-white rounded-lg shadow-md p-6 mb-6">
                <div class="flex flex-col md:flex-row items-center justify-between">
                    <div class="mb-4 md:mb-0">
                        <h3 class="text-lg font-semibold text-gray-800 mb-2">Quick Access</h3>
                        <p class="text-gray-600 mb-2">Scan this QR code with your mobile device:</p>
                        <p class="text-sm text-blue-600 font-mono break-all">{{.ServerURL}}</p>
                    </div>
                    {{if .QRCodeData}}
                    <div class="flex-shrink-0">
                        <img src="data:image/png;base64,{{.QRCodeData}}" alt="QR Code" class="w-32 h-32 border rounded-lg">
                    </div>
                    {{end}}
                </div>
            </div>
        </header>

        <!-- Search Bar -->
        <div class="mb-6">
            <div class="relative">
                <input 
                    type="text" 
                    id="fileSearch" 
                    placeholder="Search files and folders..." 
                    onkeyup="searchFiles()" 
                    class="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                <i class="fas fa-search absolute left-3 top-3 text-gray-400"></i>
            </div>
        </div>

        <!-- Upload Section -->
        <div class="mb-6 bg-white rounded-lg shadow-md overflow-hidden">
            <div class="bg-gray-100 px-6 py-3 border-b">
                <h3 class="text-lg font-semibold text-gray-800">
                    <i class="fas fa-cloud-upload-alt text-blue-600 mr-2"></i>
                    Upload Files
                </h3>
            </div>
            <div class="p-6">
                <form id="uploadForm" enctype="multipart/form-data" method="POST" action="/upload">
                    <input type="hidden" name="directory" value="{{.CurrentPath}}">
                    <div id="dropZone" class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-400 transition-colors duration-200">
                        <i class="fas fa-cloud-upload-alt text-4xl text-gray-400 mb-4"></i>
                        <p class="text-lg text-gray-600 mb-2">Drag & drop files here, or</p>
                        <label for="fileInput" class="inline-block bg-blue-600 text-white px-6 py-2 rounded-lg cursor-pointer hover:bg-blue-700 transition-colors duration-200">
                            <i class="fas fa-folder-open mr-2"></i>
                            Choose Files
                        </label>
                        <input type="file" id="fileInput" name="files" multiple style="display: none;">
                        <p class="text-sm text-gray-500 mt-2">Maximum 10MB per file</p>
                    </div>
                    <div id="uploadProgress" class="mt-4 hidden">
                        <div class="bg-gray-200 rounded-full h-2">
                            <div id="progressBar" class="bg-blue-600 h-2 rounded-full transition-all duration-300" style="width: 0%"></div>
                        </div>
                        <p id="uploadStatus" class="text-sm text-gray-600 mt-2">Uploading...</p>
                    </div>
                </form>
            </div>
        </div>

        <div class="bg-white rounded-lg shadow-md overflow-hidden">
            <div class="bg-gray-100 px-6 py-3 border-b">
                <h2 class="text-lg font-semibold text-gray-800">Files & Folders</h2>
            </div>
            
            <div class="overflow-x-auto">
                <table class="w-full" id="fileTable">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Size</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Modified</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        {{if .HasParent}}
                        <tr class="hover:bg-gray-50">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <i class="fas fa-level-up-alt text-gray-400 mr-3"></i>
                                    <a href="{{.ParentPath}}" class="text-blue-600 hover:text-blue-800 font-medium">.. (Parent Directory)</a>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">-</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">-</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">-</td>
                        </tr>
                        {{end}}
                        
                        {{range .Files}}
                        <tr class="hover:bg-gray-50">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <i class="{{.Icon}} mr-3"></i>
                                    {{if .IsDir}}
                                        <a href="{{.Path}}" class="text-blue-600 hover:text-blue-800 font-medium">{{.Name}}</a>
                                    {{else}}
                                        <span class="text-gray-900 cursor-pointer" onclick="previewFile('{{.Name}}', '{{.Path}}')">{{.Name}}</span>
                                    {{end}}
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{.SizeStr}}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{.ModTime.Format "2006-01-02 15:04:05"}}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                {{if not .IsDir}}
                                    <div class="flex space-x-2">
                                        <a href="{{.Path}}?download=1" class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                                            <i class="fas fa-download mr-1"></i>
                                            Download
                                        </a>
                                        <button onclick="previewFile('{{.Name}}', '{{.Path}}')" class="inline-flex items-center px-3 py-1 border border-gray-300 text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                                            <i class="fas fa-eye mr-1"></i>
                                            Preview
                                        </button>
                                    </div>
                                {{else}}
                                    <a href="{{.Path}}?download=zip" class="inline-flex items-center px-3 py-1 border border-gray-300 text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                                        <i class="fas fa-file-archive mr-1"></i>
                                        Zip Download
                                    </a>
                                {{end}}
                            </td>
                        </tr>
                        {{end}}
                        
                        {{if not .Files}}
                        <tr>
                            <td colspan="4" class="px-6 py-8 text-center text-gray-500">
                                <i class="fas fa-folder-open text-4xl mb-2 text-gray-300"></i>
                                <p>This directory is empty</p>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- File Preview Modal -->
        <div id="previewModal" class="fixed inset-0 bg-black bg-opacity-50 hidden z-50" onclick="closePreview()">
            <div class="flex items-center justify-center min-h-screen p-4">
                <div class="bg-white rounded-lg max-w-4xl max-h-screen overflow-auto" onclick="event.stopPropagation()">
                    <div class="p-4 border-b flex justify-between items-center">
                        <h3 id="previewTitle" class="text-lg font-semibold">File Preview</h3>
                        <button onclick="closePreview()" class="text-gray-500 hover:text-gray-700">
                            <i class="fas fa-times"></i>
                        </button>
                    </div>
                    <div id="previewContent" class="p-4">
                        <!-- Preview content will be loaded here -->
                    </div>
                </div>
            </div>
        </div>
        
        <footer class="mt-8 text-center text-gray-500 text-sm">
            <p>Powered by <strong>GoShare</strong> - Easy file sharing over Wi-Fi</p>
        </footer>
    </div>

    <script>
        function previewFile(fileName, filePath) {
            const modal = document.getElementById('previewModal');
            const title = document.getElementById('previewTitle');
            const content = document.getElementById('previewContent');
            
            title.textContent = fileName;
            
            const ext = fileName.toLowerCase().split('.').pop();
            
            if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(ext)) {
                content.innerHTML = '<img src="' + filePath + '" class="max-w-full h-auto rounded" alt="' + fileName + '">';
            } else if (['txt', 'md', 'json', 'css', 'js', 'html', 'xml', 'csv'].includes(ext)) {
                fetch(filePath)
                    .then(response => response.text())
                    .then(text => {
                        content.innerHTML = '<pre class="bg-gray-100 p-4 rounded overflow-auto max-h-96 text-sm"><code>' + 
                            text.replace(/</g, '&lt;').replace(/>/g, '&gt;') + '</code></pre>';
                    })
                    .catch(() => {
                        content.innerHTML = '<p class="text-red-500">Unable to preview this file.</p>';
                    });
            } else {
                content.innerHTML = '<p class="text-gray-500">Preview not available for this file type. <a href="' + filePath + '?download=1" class="text-blue-600 hover:underline">Download instead</a></p>';
            }
            
            modal.classList.remove('hidden');
        }
        
        function closePreview() {
            document.getElementById('previewModal').classList.add('hidden');
        }

        // Drag & Drop Upload Functionality
        const dropZone = document.getElementById('dropZone');
        const fileInput = document.getElementById('fileInput');
        const uploadForm = document.getElementById('uploadForm');
        const uploadProgress = document.getElementById('uploadProgress');
        const progressBar = document.getElementById('progressBar');
        const uploadStatus = document.getElementById('uploadStatus');

        // Prevent default drag behaviors
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, preventDefaults, false);
            document.body.addEventListener(eventName, preventDefaults, false);
        });

        // Highlight drop zone when item is dragged over it
        ['dragenter', 'dragover'].forEach(eventName => {
            dropZone.addEventListener(eventName, highlight, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, unhighlight, false);
        });

        // Handle dropped files
        dropZone.addEventListener('drop', handleDrop, false);

        // Handle file input change
        fileInput.addEventListener('change', function(e) {
            handleFiles(e.target.files);
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        function highlight(e) {
            dropZone.classList.add('border-blue-400', 'bg-blue-50');
        }

        function unhighlight(e) {
            dropZone.classList.remove('border-blue-400', 'bg-blue-50');
        }

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            handleFiles(files);
        }

        function handleFiles(files) {
            if (files.length === 0) return;

            // Create FormData object
            const formData = new FormData();
            formData.append('directory', document.querySelector('input[name="directory"]').value);

            // Add all files to form data
            Array.from(files).forEach(file => {
                formData.append('files', file);
            });

            // Show progress
            uploadProgress.classList.remove('hidden');
            uploadStatus.textContent = 'Uploading ' + files.length + ' file(s)...';
            progressBar.style.width = '0%';

            // Upload files
            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    progressBar.style.width = '100%';
                    uploadStatus.textContent = 'Upload completed successfully!';
                    setTimeout(() => {
                        window.location.reload();
                    }, 1000);
                } else {
                    throw new Error('Upload failed');
                }
            })
            .catch(error => {
                uploadStatus.textContent = 'Upload failed. Please try again.';
                uploadStatus.classList.add('text-red-600');
            });
        }
    </script>
</body>
</html>
`

// FileHandler handles HTTP requests for file browsing and downloading
type FileHandler struct {
	rootDir   string
	template  *template.Template
	serverURL string
	password  string
}

// ServeHTTP implements the http.Handler interface
func (fh *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for React frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle auth check endpoint (not protected by auth middleware but checks auth status)
	if r.URL.Path == "/api/auth/check" {
		w.Header().Set("Content-Type", "application/json")

		// Check if user is authenticated
		isAuthenticated := false

		// If no password is set, everyone is authenticated
		if fh.password == "" {
			isAuthenticated = true
		} else {
			// Check for valid session cookie
			if cookie, err := r.Cookie("auth_session"); err == nil && cookie.Value == "authenticated" {
				isAuthenticated = true
			} else {
				// Check basic auth as fallback
				_, pass, ok := r.BasicAuth()
				if ok && pass == fh.password {
					isAuthenticated = true
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		if isAuthenticated {
			w.Write([]byte(`{"authenticated": true}`))
		} else {
			w.Write([]byte(`{"authenticated": false}`))
		}
		return
	}

	// Handle API endpoints
	if strings.HasPrefix(r.URL.Path, "/api/") {
		fh.handleAPI(w, r)
		return
	}

	// Handle upload
	if r.Method == "POST" && r.URL.Path == "/upload" {
		fh.handleUpload(w, r)
		return
	}

	requestPath := r.URL.Path
	if requestPath == "" || requestPath == "/" {
		requestPath = "/"
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(requestPath)
	if cleanPath == "." {
		cleanPath = "/"
	}

	// Convert URL path to filesystem path
	fsPath := filepath.Join(fh.rootDir, strings.TrimPrefix(cleanPath, "/"))

	// Security check: ensure the path is within the root directory
	if !strings.HasPrefix(fsPath, fh.rootDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	stat, err := os.Stat(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check for zip download request for directories
	if stat.IsDir() && r.URL.Query().Get("download") == "zip" {
		fh.serveDirectoryAsZip(w, r, fsPath, stat.Name())
		return
	}

	// If it's a file, serve it for download
	if !stat.IsDir() {
		fh.serveFile(w, r, fsPath, stat)
		return
	}

	// If it's a directory, show the file listing
	fh.serveDirectory(w, r, fsPath, cleanPath)
}

// serveFile serves a file for download
func (fh *FileHandler) serveFile(w http.ResponseWriter, r *http.Request, fsPath string, stat os.FileInfo) {
	// Check if download is requested
	if r.URL.Query().Get("download") == "1" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	}

	// Set content type based on file extension
	w.Header().Set("Content-Type", getContentType(fsPath))

	file, err := os.Open(fsPath)
	if err != nil {
		http.Error(w, "Could not open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

// serveDirectory serves a directory listing
func (fh *FileHandler) serveDirectory(w http.ResponseWriter, r *http.Request, fsPath, urlPath string) {
	entries, err := os.ReadDir(fsPath)
	if err != nil {
		http.Error(w, "Could not read directory", http.StatusInternalServerError)
		return
	}

	// Convert entries to FileInfo
	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Name:    info.Name(),
			Path:    filepath.Join(urlPath, info.Name()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
			Icon:    getFileIcon(info.Name(), info.IsDir()),
			SizeStr: formatFileSize(info.Size(), info.IsDir()),
		}
		files = append(files, fileInfo)
	}

	// Sort files: directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	// Determine parent path
	var parentPath string
	hasParent := urlPath != "/" && urlPath != ""
	if hasParent {
		parentPath = filepath.Dir(urlPath)
		if parentPath == "." {
			parentPath = "/"
		}
	}

	// Generate QR code for server URL
	var qrCodeData string
	if fh.serverURL != "" {
		qr, err := qrcode.New(fh.serverURL, qrcode.Medium)
		if err == nil {
			qrBytes, err := qr.PNG(256)
			if err == nil {
				qrCodeData = base64.StdEncoding.EncodeToString(qrBytes)
			}
		}
	}

	// Prepare template data
	data := PageData{
		Title:       "GoShare - File Browser",
		CurrentPath: urlPath,
		ParentPath:  parentPath,
		Files:       files,
		HasParent:   hasParent,
		ServerURL:   fh.serverURL,
		QRCodeData:  qrCodeData,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := fh.template.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

// serveDirectoryAsZip serves a directory as a zip file
func (fh *FileHandler) serveDirectoryAsZip(w http.ResponseWriter, r *http.Request, fsPath, dirName string) {
	// Set headers for zip download
	zipFilename := dirName + ".zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFilename))

	// Create zip writer
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Walk through directory and add files to zip
	err := filepath.Walk(fsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == fsPath {
			return nil
		}

		// Get relative path for zip entry
		relPath, err := filepath.Rel(fsPath, path)
		if err != nil {
			return err
		}

		// Create zip entry
		if info.IsDir() {
			// Create directory entry
			_, err := zipWriter.Create(relPath + "/")
			return err
		} else {
			// Create file entry
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			// Open source file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy file contents to zip
			_, err = io.Copy(zipFile, file)
			return err
		}
	})

	if err != nil {
		log.Printf("Error creating zip: %v", err)
		// Since we've already started writing to response, we can't send a proper error
		return
	}
}

// getFileIcon returns the appropriate Font Awesome icon for a file
func getFileIcon(filename string, isDir bool) string {
	if isDir {
		return "fas fa-folder text-blue-500"
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt", ".md", ".readme":
		return "fas fa-file-alt text-gray-600"
	case ".pdf":
		return "fas fa-file-pdf text-red-600"
	case ".doc", ".docx":
		return "fas fa-file-word text-blue-600"
	case ".xls", ".xlsx":
		return "fas fa-file-excel text-green-600"
	case ".ppt", ".pptx":
		return "fas fa-file-powerpoint text-orange-600"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "fas fa-file-archive text-purple-600"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "fas fa-file-image text-pink-600"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return "fas fa-file-audio text-green-600"
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv":
		return "fas fa-file-video text-red-600"
	case ".html", ".htm", ".css", ".js", ".json", ".xml":
		return "fas fa-file-code text-blue-600"
	case ".go", ".py", ".java", ".cpp", ".c", ".h", ".php", ".rb", ".rs":
		return "fas fa-file-code text-green-600"
	default:
		return "fas fa-file text-gray-600"
	}
}

// formatFileSize formats file size in human-readable format
func formatFileSize(size int64, isDir bool) string {
	if isDir {
		return "-"
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// getContentType returns the MIME type for a file
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

func StartServer(dir string, port int, password string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)

	// Custom file handler for API and file serving
	handler := &FileHandler{
		rootDir:   absDir,
		template:  template.Must(template.New("index").Parse(htmlTemplate)),
		serverURL: url,
		password:  password,
	}

	// Set up routes
	mux := http.NewServeMux()

	// We'll handle all routing in the main handler function below
	// No need for individual route handlers since we're using a custom dispatcher	// Serve React build files (check if frontend/build exists)
	frontendPath := filepath.Join(absDir, "frontend", "build")
	if _, err := os.Stat(frontendPath); err == nil {
		// Create a static file server for React
		reactFS := http.FileServer(http.Dir(frontendPath))

		// Custom handler that routes correctly
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Check if this is an API route that should be handled by our handlers
			switch {
			case strings.HasPrefix(r.URL.Path, "/api/"):
				handler.ServeHTTP(w, r)
			case r.URL.Path == "/login":
				// Login should go through auth middleware to handle the login logic
				applyAuthMiddleware(handler, password).ServeHTTP(w, r)
			case r.URL.Path == "/upload":
				applyAuthMiddleware(handler, password).ServeHTTP(w, r)
			case strings.HasPrefix(r.URL.Path, "/files/"):
				applyAuthMiddleware(handler, password).ServeHTTP(w, r)
			default:
				// Serve React app - if file doesn't exist, serve index.html for React Router
				if _, err := os.Stat(filepath.Join(frontendPath, r.URL.Path)); os.IsNotExist(err) && r.URL.Path != "/" {
					http.ServeFile(w, r, filepath.Join(frontendPath, "index.html"))
				} else {
					reactFS.ServeHTTP(w, r)
				}
			}
		})
		fmt.Printf("🚀 Serving React frontend from: %s\n", frontendPath)
	} else {
		// Fallback to original file browser
		mux.Handle("/", applyAuthMiddleware(handler, password))
		fmt.Printf("📂 Serving original file browser\n")
	}

	fmt.Printf("📂 Serving %s at:\n➡️  %s\n", absDir, url)

	// Generate and display local QR code
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Fatalf("QR generation failed: %v", err)
	}
	fmt.Println("\n📱 Scan this QR to open (local):")
	fmt.Println(qr.ToSmallString(false))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// handleUpload handles file uploads via drag & drop or file selection
func (fh *FileHandler) handleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the target directory from form data
	targetDir := r.FormValue("directory")
	if targetDir == "" {
		targetDir = "/"
	}

	// Clean and validate the target directory path
	cleanDir := filepath.Clean(targetDir)
	if cleanDir == "." {
		cleanDir = "/"
	}

	// Convert to filesystem path
	fsDir := filepath.Join(fh.rootDir, strings.TrimPrefix(cleanDir, "/"))

	// Security check: ensure the path is within the root directory
	if !strings.HasPrefix(fsDir, fh.rootDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(fsDir, 0755)
	if err != nil {
		http.Error(w, "Unable to create directory", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	uploadedCount := 0

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Create the destination file
		destPath := filepath.Join(fsDir, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			continue
		}
		defer destFile.Close()

		// Copy file contents
		_, err = io.Copy(destFile, file)
		if err != nil {
			os.Remove(destPath) // Clean up on error
			continue
		}

		uploadedCount++
	}

	// Redirect back to the directory with a success message
	redirectURL := cleanDir
	if uploadedCount > 0 {
		if strings.Contains(redirectURL, "?") {
			redirectURL += "&uploaded=" + fmt.Sprintf("%d", uploadedCount)
		} else {
			redirectURL += "?uploaded=" + fmt.Sprintf("%d", uploadedCount)
		}
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleAPI handles API endpoints for the React frontend
func (fh *FileHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api")

	switch {
	case path == "/files" || strings.HasPrefix(path, "/files/"):
		fh.handleAPIFiles(w, r)
	case path == "/auth/check":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
	default:
		http.NotFound(w, r)
	}
}

// handleAPIFiles handles file listing API endpoints
func (fh *FileHandler) handleAPIFiles(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Query().Get("path")
	if requestPath == "" {
		requestPath = "/"
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(requestPath)
	if cleanPath == "." {
		cleanPath = "/"
	}

	// Convert URL path to filesystem path
	fsPath := filepath.Join(fh.rootDir, strings.TrimPrefix(cleanPath, "/"))

	// Security check: ensure the path is within the root directory
	if !strings.HasPrefix(fsPath, fh.rootDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	stat, err := os.Stat(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if !stat.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	// Read directory contents
	entries, err := os.ReadDir(fsPath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	// Create API response
	var files []APIFileItem
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip hidden files (starting with .)
		if strings.HasPrefix(info.Name(), ".") {
			continue
		}

		filePath := filepath.Join(cleanPath, info.Name())
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}

		apiFile := APIFileItem{
			Name:          info.Name(),
			Path:          filePath,
			Size:          info.Size(),
			IsDir:         info.IsDir(),
			ModTime:       info.ModTime(),
			DownloadCount: 0, // TODO: implement download tracking
		}

		files = append(files, apiFile)
	}

	// Sort files: directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	// Determine parent path
	parentPath := "/"
	hasParent := cleanPath != "/"
	if hasParent {
		parentPath = filepath.Dir(cleanPath)
		if parentPath == "." {
			parentPath = "/"
		}
	}

	pageData := APIPageData{
		Title:       "GoShare - File Browser",
		CurrentPath: cleanPath,
		ParentPath:  parentPath,
		Files:       files,
		HasParent:   hasParent,
		ServerURL:   fh.serverURL,
	}

	json.NewEncoder(w).Encode(pageData)
}

func applyAuthMiddleware(h http.Handler, password string) http.Handler {
	if password == "" {
		return h // no protection
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle login form submission
		if r.Method == "POST" && r.URL.Path == "/login" {
			r.ParseForm()
			submittedPassword := r.FormValue("password")
			if submittedPassword == password {
				// Set a session cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_session",
					Value:    "authenticated",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   86400, // 24 hours
				})
				redirectTo := r.FormValue("redirect")
				if redirectTo == "" {
					redirectTo = "/"
				}
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			} else {
				// Wrong password, show login form with error
				showLoginForm(w, r, "Invalid password. Please try again.")
				return
			}
		}

		// Check for valid session cookie
		if cookie, err := r.Cookie("auth_session"); err == nil && cookie.Value == "authenticated" {
			h.ServeHTTP(w, r)
			return
		}

		// Check basic auth as fallback
		_, pass, ok := r.BasicAuth()
		if ok && pass == password {
			h.ServeHTTP(w, r)
			return
		}

		// Show login form
		showLoginForm(w, r, "")
	})
}

func showLoginForm(w http.ResponseWriter, r *http.Request, errorMsg string) {
	loginHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoShare - Login</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
    <div class="max-w-md w-full space-y-8 p-8">
        <div class="text-center">
            <i class="fas fa-shield-alt text-4xl text-blue-600 mb-4"></i>
            <h2 class="text-3xl font-bold text-gray-900">Access Required</h2>
            <p class="mt-2 text-sm text-gray-600">Please enter the password to access GoShare</p>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
            <form method="POST" action="/login" class="space-y-6">
                <input type="hidden" name="redirect" value="` + r.URL.String() + `">
                
                ` + func() string {
		if errorMsg != "" {
			return `<div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
                        <i class="fas fa-exclamation-triangle mr-2"></i>
                        ` + errorMsg + `
                    </div>`
		}
		return ""
	}() + `
                
                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700 mb-2">Password</label>
                    <div class="relative">
                        <input 
                            type="password" 
                            id="password" 
                            name="password" 
                            required 
                            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 pl-12"
                            placeholder="Enter password"
                            autofocus
                        >
                        <i class="fas fa-lock absolute left-4 top-4 text-gray-400"></i>
                    </div>
                </div>
                
                <button 
                    type="submit" 
                    class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 font-medium"
                >
                    <i class="fas fa-sign-in-alt mr-2"></i>
                    Access GoShare
                </button>
            </form>
        </div>
        
        <div class="text-center text-sm text-gray-500">
            <p>Powered by <strong>GoShare</strong> - Secure file sharing</p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(loginHTML))
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "localhost"
}
