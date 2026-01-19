# GoShare

A modern, secure file sharing application with a beautiful React frontend and Go backend. Share files over your local network with style and security.

## ✨ Features

- 🎨 **Modern React UI**: Beautiful, responsive interface built with React, TypeScript, and Tailwind CSS
- 🌓 **Dark/Light Mode**: Elegant theme switching with user preference persistence
- 🔒 **Secure Authentication**: Session-based authentication with dynamic HTTP-only tokens and optional Basic Auth
- 📱 **QR Code Access**: Automatic QR code generation for easy mobile access (Local & Public)
- 📤 **Drag & Drop Upload**: Intuitive file upload with progress feedback and configurable size limits
- 🔍 **Real-time Search**: Instant file and folder search functionality
- 📊 **File Information**: Display file sizes, modification dates, and real-time download counts
- 🎯 **Download Tracking**: Keep track of how many times each file has been downloaded during the session
- 🚀 **Single Binary**: Easy deployment with embedded React build
- 🌐 **RESTful API**: Clean JSON API for programmatic access
- 🗺️ **Public Tunneling**: Built-in `ngrok` support to share files across the internet securely

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
- **Session-based authentication** with dynamic tokens
- **CORS support** for frontend communication
- **RESTful JSON API**
- **Zip Compression** for on-the-fly directory downloads

## 🚀 Quick Start

### Build & Run (Recommended)

1. **Prerequisites**: Go 1.21+ and Node.js 18+
2. **Build everything**:

   ```bash
   # Build the frontend
   cd frontend && npm install && npm run build && cd ..

   # Build the Go binary
   go build -o goshare
   ```

3. **Run**:
   ```bash
   ./goshare --password mysecurepass /path/to/share
   ```
4. **Access**: Scan the QR code in your terminal or open `http://localhost:8081`

## 📖 How it Works

GoShare is designed to be a "run-and-done" utility.

1. **Local Mode**: When run, it detects your local IP and starts a server. Anyone on your Wi-Fi (who knows your password) can scan the QR code and immediately access/upload files.
2. **Public Mode**: Use the `--ngrok` flag to generate a public URL. This allows you to share files with someone across the world without configuring port forwarding.
3. **Security**: Every time the server starts, a unique 32-byte session token is generated. This token is used for session cookies, ensuring that old sessions are invalidated on restart.

## 🖥️ Command Line Options

```bash
./goshare [flags] [directory]

Flags:
  -h, --help              Help for goshare
  -p, --password string   Password for accessing files (optional)
      --port int          Port to run server on (default 8080)
      --max-size int      Maximum upload size in MB (default 100)
      --ngrok             Expose server to the internet using ngrok
  -d, --dir string        Directory to share (default ".")
```

### Examples

```bash
# Share with 500MB upload limit
./goshare --password pass --max-size 500

# Share over the internet
./goshare --password pass --ngrok
```

## 🔐 Security Features

- **Dynamic Sessions**: Session tokens are randomized on every startup.
- **Path Validation**: Strict checks prevent directory traversal attacks.
- **Secure Cookies**: HTTP-only cookies prevent XSS-based token theft.
- **Sanitized Uploads**: Filenames and paths are cleaned before saving to disk.

## 🏗️ Development

### Development Setup

```bash
# Backend (terminal 1)
go run main.go --port 8080 --password dev /path/to/serve

# Frontend development (terminal 2)
cd frontend
npm install
npm start  # Proxies requests to backend on 8080
```

---

**GoShare** - Share files beautifully. Built with ❤️ using Go and React.
