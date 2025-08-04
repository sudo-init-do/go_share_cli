package server

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// PageData contains data for the HTML template
type PageData struct {
	Title       string
	CurrentPath string
	ParentPath  string
	Files       []FileInfo
	HasParent   bool
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
</head>
<body class="bg-gray-50 min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-6xl">
        <header class="mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">
                <i class="fas fa-share-alt text-blue-600 mr-2"></i>
                GoShare File Browser
            </h1>
            <p class="text-gray-600">Current directory: <code class="bg-gray-200 px-2 py-1 rounded">{{.CurrentPath}}</code></p>
        </header>

        <div class="bg-white rounded-lg shadow-md overflow-hidden">
            <div class="bg-gray-100 px-6 py-3 border-b">
                <h2 class="text-lg font-semibold text-gray-800">Files & Folders</h2>
            </div>
            
            <div class="overflow-x-auto">
                <table class="w-full">
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
                                        <span class="text-gray-900">{{.Name}}</span>
                                    {{end}}
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{.SizeStr}}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{.ModTime.Format "2006-01-02 15:04:05"}}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                {{if not .IsDir}}
                                    <a href="{{.Path}}?download=1" class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                                        <i class="fas fa-download mr-1"></i>
                                        Download
                                    </a>
                                {{else}}
                                    <span class="text-gray-400">Folder</span>
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
        
        <footer class="mt-8 text-center text-gray-500 text-sm">
            <p>Powered by <strong>GoShare</strong> - Easy file sharing over Wi-Fi</p>
        </footer>
    </div>
</body>
</html>
`

// FileHandler handles HTTP requests for file browsing and downloading
type FileHandler struct {
	rootDir  string
	template *template.Template
}

// ServeHTTP implements the http.Handler interface
func (fh *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Prepare template data
	data := PageData{
		Title:       "GoShare - File Browser",
		CurrentPath: urlPath,
		ParentPath:  parentPath,
		Files:       files,
		HasParent:   hasParent,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := fh.template.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
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

	// Custom file handler
	handler := &FileHandler{
		rootDir:  absDir,
		template: template.Must(template.New("index").Parse(htmlTemplate)),
	}

	// Wrap with auth middleware if password provided
	http.Handle("/", applyAuthMiddleware(handler, password))

	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	fmt.Printf("📂 Serving %s at:\n➡️  %s\n", absDir, url)

	// Generate and display local QR code
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Fatalf("QR generation failed: %v", err)
	}
	fmt.Println("\n📱 Scan this QR to open (local):")
	fmt.Println(qr.ToSmallString(false))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func applyAuthMiddleware(h http.Handler, password string) http.Handler {
	if password == "" {
		return h // no protection
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := r.BasicAuth()
		if !ok || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="go_share_cli"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
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
