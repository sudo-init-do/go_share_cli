import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  FolderIcon, 
  DocumentIcon, 
  ArrowUpIcon,
  CloudArrowUpIcon,
  MagnifyingGlassIcon,
  SunIcon,
  MoonIcon,
  ArrowDownTrayIcon,
  EyeIcon
} from '@heroicons/react/24/outline';
import { useDropzone } from 'react-dropzone';
import { fileService } from '../services/api';
import { FileItem, PageData } from '../types';
import toast from 'react-hot-toast';

interface FileBrowserProps {
  onLogout: () => void;
}

export const FileBrowser: React.FC<FileBrowserProps> = ({ onLogout }) => {
  const [pageData, setPageData] = useState<PageData | null>(null);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [darkMode, setDarkMode] = useState(() => {
    // Check localStorage for theme preference, default to dark
    const saved = localStorage.getItem('theme');
    return saved ? saved === 'dark' : true;
  });

  useEffect(() => {
    loadFiles('/');
  }, []);

  useEffect(() => {
    // Save theme preference
    localStorage.setItem('theme', darkMode ? 'dark' : 'light');
  }, [darkMode]);

  const loadFiles = async (path: string = '/') => {
    try {
      setLoading(true);
      const data = await fileService.getFiles(path);
      setPageData(data);
    } catch (error) {
      toast.error('Failed to load files');
    } finally {
      setLoading(false);
    }
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: async (acceptedFiles) => {
      if (!pageData) return;

      const fileList = acceptedFiles as unknown as FileList;
      
      try {
        toast.promise(
          fileService.uploadFiles(fileList, pageData.currentPath),
          {
            loading: 'Uploading files...',
            success: 'Files uploaded successfully!',
            error: 'Upload failed'
          }
        );
        
        // Reload files after upload
        setTimeout(() => loadFiles(pageData.currentPath), 1000);
      } catch (error) {
        console.error('Upload error:', error);
      }
    },
    noClick: true,
    noKeyboard: true
  });

  const filteredFiles = pageData?.files.filter(file =>
    file.name.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  const getFileIcon = (file: FileItem) => {
    if (file.isDir) {
      return <FolderIcon className="h-6 w-6 text-blue-500" />;
    }
    
    const ext = file.name.split('.').pop()?.toLowerCase();
    const iconClass = "h-6 w-6";
    
    switch (ext) {
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
      case 'webp':
        return <DocumentIcon className={`${iconClass} text-green-500`} />;
      case 'pdf':
        return <DocumentIcon className={`${iconClass} text-red-500`} />;
      case 'doc':
      case 'docx':
        return <DocumentIcon className={`${iconClass} text-blue-500`} />;
      case 'txt':
      case 'md':
        return <DocumentIcon className={`${iconClass} text-gray-500`} />;
      default:
        return <DocumentIcon className={`${iconClass} text-gray-400`} />;
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen transition-colors duration-300 ${darkMode ? 'dark' : ''}`} {...getRootProps()}>
      <input {...getInputProps()} />
      
      {/* Drag Overlay */}
      <AnimatePresence>
        {isDragActive && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-blue-500/20 backdrop-blur-sm z-50 flex items-center justify-center"
          >
            <div className="card p-8 text-center">
              <CloudArrowUpIcon className="h-16 w-16 text-blue-500 mx-auto mb-4" />
              <p className="text-xl font-semibold text-gray-900 dark:text-white">
                Drop files here to upload
              </p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="bg-gray-50 dark:bg-gray-900 min-h-screen">
        <div className="container mx-auto px-4 py-8 max-w-7xl">
          {/* Header */}
          <header className="mb-8">
            <div className="flex items-center justify-between mb-6">
              <motion.h1 
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                className="text-3xl font-bold text-gray-900 dark:text-white"
              >
                🗂️ GoShare File Browser
              </motion.h1>
              
              <div className="flex items-center space-x-3">
                <button
                  onClick={() => setDarkMode(!darkMode)}
                  className="p-2 rounded-lg bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                >
                  {darkMode ? (
                    <SunIcon className="h-5 w-5 text-yellow-500" />
                  ) : (
                    <MoonIcon className="h-5 w-5 text-gray-600" />
                  )}
                </button>
                
                <button
                  onClick={onLogout}
                  className="btn-secondary"
                >
                  Logout
                </button>
              </div>
            </div>

            {/* Breadcrumb */}
            <nav className="mb-6">
              <div className="text-sm text-gray-600 dark:text-gray-300">
                <span className="font-medium">Current Path:</span> {pageData?.currentPath || '/'}
              </div>
            </nav>

            {/* Search */}
            <div className="relative max-w-md">
              <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search files and folders..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 pr-4 py-2 w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </header>

          {/* Upload Zone */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="card p-6 mb-8"
          >
            <div className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center hover:border-blue-400 transition-colors">
              <CloudArrowUpIcon className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-lg text-gray-600 dark:text-gray-300 mb-2">
                Drag & drop files here, or click to select
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Upload files to {pageData?.currentPath || '/'}
              </p>
            </div>
          </motion.div>

          {/* File List */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="card overflow-hidden"
          >
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Files & Folders ({filteredFiles.length})
              </h2>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Size
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Modified
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {/* Parent Directory */}
                  {pageData?.hasParent && (
                    <motion.tr
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer"
                      onClick={() => loadFiles(pageData.parentPath)}
                    >
                      <td className="px-6 py-4">
                        <div className="flex items-center">
                          <ArrowUpIcon className="h-6 w-6 text-gray-400 mr-3" />
                          <span className="text-blue-600 dark:text-blue-400 font-medium">
                            .. (Parent Directory)
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">-</td>
                      <td className="px-6 py-4 text-sm text-gray-500">-</td>
                      <td className="px-6 py-4">-</td>
                    </motion.tr>
                  )}

                  {/* Files and Folders */}
                  {filteredFiles.map((file, index) => (
                    <motion.tr
                      key={file.path}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.05 }}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800"
                    >
                      <td className="px-6 py-4">
                        <div className="flex items-center">
                          {getFileIcon(file)}
                          <span
                            className={`ml-3 font-medium ${
                              file.isDir 
                                ? 'text-blue-600 dark:text-blue-400 cursor-pointer hover:underline' 
                                : 'text-gray-900 dark:text-white'
                            }`}
                            onClick={() => file.isDir && loadFiles(file.path)}
                          >
                            {file.name}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500 dark:text-gray-400">
                        {file.isDir ? '-' : formatFileSize(file.size)}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500 dark:text-gray-400">
                        {formatDate(file.modTime)}
                      </td>
                      <td className="px-6 py-4">
                        {!file.isDir && (
                          <div className="flex space-x-2">
                            <button
                              onClick={() => window.open(fileService.getPreviewUrl(file.path), '_blank')}
                              className="p-1 text-gray-400 hover:text-blue-600 transition-colors"
                              title="Preview"
                            >
                              <EyeIcon className="h-5 w-5" />
                            </button>
                            <button
                              onClick={() => window.open(fileService.getDownloadUrl(file.path), '_blank')}
                              className="p-1 text-gray-400 hover:text-green-600 transition-colors"
                              title="Download"
                            >
                              <ArrowDownTrayIcon className="h-5 w-5" />
                            </button>
                          </div>
                        )}
                      </td>
                    </motion.tr>
                  ))}

                  {filteredFiles.length === 0 && (
                    <tr>
                      <td colSpan={4} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                        No files found
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </motion.div>

          {/* Footer */}
          <footer className="mt-8 text-center text-gray-500 dark:text-gray-400 text-sm">
            <p>Powered by <strong>GoShare</strong> - Modern file sharing</p>
          </footer>
        </div>
      </div>
    </div>
  );
};
