# GoShare

A modern, secure file sharing application with a beautiful React frontend and Go backend. Share files over your local network with style and security.

## ✨ Features

- 🎨 **Modern React UI**: Beautiful, responsive interface built with React, TypeScript, and Tailwind CSS
- 🌓 **Dark/Light Mode**: Elegant theme switching with user preference persistence
- 🔒 **Secure Authentication**: Session-based authentication with HTTP-only cookies
- 📱 **QR Code Access**: Automatic QR code generation for easy mobile access
- 📤 **Drag & Drop Upload**: Intuitive file upload with progress feedback
- 🔍 **Real-time Search**: Instant file and folder search functionality
- 📊 **File Information**: Display file sizes, modification dates, and download counts
- 🎯 **Download Tracking**: Keep track of how many times files have been downloaded
- 🚀 **Single Binary**: Easy deployment with embedded React build
- 🌐 **RESTful API**: Clean JSON API for programmatic access
- 📱 **Mobile Responsive**: Perfect experience on all devices

## 🛠️ Technology Stack

### Frontend
- **React 18** with TypeScript
- **Tailwind CSS** for styling
- **Framer Motion** for smooth animations
- **React Hot Toast** for notifications
- **Heroicons** for beautiful icons
- **React Dropzone** for file uploads

### Backend
- **Go** with built-in HTTP server
- **Session-based authentication**
- **CORS support** for frontend communication
- **RESTful JSON API**

## 🚀 Quick Start

### Download & Run
1. Download the latest binary from [Releases](https://github.com/yourusername/goshare/releases)
2. Make it executable: `chmod +x goshare`
3. Run: `./goshare --password mypassword ~/Documents`
4. Open `http://localhost:8081` in your browser

### Build from Source
```bash
git clone https://github.com/yourusername/goshare.git
cd goshare

# Build frontend
cd frontend && npm install && npm run build && cd ..

# Build Go binary
go build -o goshare

# Run
./goshare --password mypassword /path/to/share
```

## 📖 Usage

### Command Line Options
```bash
./goshare [flags] [directory]

Flags:
  -h, --help              Help for goshare
  -p, --password string   Password for accessing files (optional)
      --port int          Port to run server on (default 8081)
```

### Examples
```bash
# Share current directory (no password)
./goshare

# Share with password protection
./goshare --password secretpass ~/Documents

# Custom port
./goshare --port 3000 ~/Downloads

# All options
./goshare --password mypass --port 9000 ~/Pictures
```

## 🖥️ Web Interface

The React frontend provides a modern, intuitive experience:

- **📁 File Browser**: Navigate directories with breadcrumb navigation
- **🔍 Search**: Real-time search across files and folders
- **📤 Upload**: Drag files directly onto the interface
- **🌓 Themes**: Toggle between beautiful dark and light modes
- **📱 Responsive**: Perfect on desktop, tablet, and mobile
- **⚡ Fast**: Optimized performance with React best practices

## 🔐 Security Features

- **Path Validation**: Prevents directory traversal attacks
- **Session Authentication**: Secure HTTP-only cookie sessions
- **CORS Protection**: Configured for safe cross-origin requests
- **Input Sanitization**: All user inputs are validated and sanitized
- **Local Network**: Accessible only within your local network

## 🔧 API Documentation

GoShare provides a RESTful JSON API:

### Authentication
```bash
# Check auth status
GET /api/auth/check

# Login
POST /login
Content-Type: application/x-www-form-urlencoded
Body: password=yourpassword
```

### Files
```bash
# Get file listing
GET /api/files?path=/some/directory

# Upload files
POST /upload
Content-Type: multipart/form-data
```

For detailed API documentation, see [CODEBASE.md](CODEBASE.md).

## 🏗️ Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Development Setup
```bash
# Backend (terminal 1)
go run main.go --port 8081 --password dev /path/to/serve

# Frontend development (terminal 2)
cd frontend
npm install
npm start  # Runs on http://localhost:3000
```

### Building for Production
```bash
# Build frontend
cd frontend && npm run build && cd ..

# Build backend with embedded frontend
go build -o goshare
```

## 📚 Documentation

- **[CODEBASE.md](CODEBASE.md)** - Comprehensive technical documentation
- **Architecture overview** - How frontend and backend work together
- **API reference** - Complete endpoint documentation
- **Development guide** - Setup and contribution guidelines

## 🤝 Contributing

We welcome contributions! Please see [CODEBASE.md](CODEBASE.md) for:
- Development setup instructions
- Code style guidelines
- Architecture explanations
- Security considerations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎯 Roadmap

- [ ] Multiple file selection for download
- [ ] File preview for common formats
- [ ] User management and permissions
- [ ] File sharing via links
- [ ] Docker image
- [ ] HTTPS support

---

**GoShare** - Share files beautifully. Built with ❤️ using Go and React.
