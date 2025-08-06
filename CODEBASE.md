# GoShare Codebase Documentation

## 📋 Table of Contents
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Backend (Go)](#backend-go)
- [Frontend (React)](#frontend-react)
- [API Documentation](#api-documentation)
- [Authentication System](#authentication-system)
- [File Upload System](#file-upload-system)
- [Development Setup](#development-setup)
- [Building and Deployment](#building-and-deployment)
- [Contributing](#contributing)

## 🎯 Project Overview

GoShare is a modern, secure file sharing application built with Go backend and React frontend. It provides a beautiful web interface for browsing, uploading, and managing files over a local network with optional password protection.

## 🏗️ Architecture

```
┌─────────────────┐    HTTP/API    ┌─────────────────┐
│                 │   Requests     │                 │
│  React Frontend │ ◄─────────────► │   Go Backend    │
│   (TypeScript)  │    JSON/REST   │   (HTTP Server) │
│                 │                │                 │
└─────────────────┘                └─────────────────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────┐                ┌─────────────────┐
│  Static Files   │                │  File System    │
│  (Tailwind CSS) │                │   (Local Dir)   │
└─────────────────┘                └─────────────────┘
```

## 🛠️ Technology Stack

### Backend
- **Go 1.21+**: Core programming language
- **net/http**: Built-in HTTP server
- **html/template**: Template rendering (fallback)
- **github.com/skip2/go-qrcode**: QR code generation
- **github.com/spf13/cobra**: CLI framework

### Frontend
- **React 18**: UI framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling framework
- **Framer Motion**: Animations
- **React Hot Toast**: Notifications
- **Heroicons**: Icon library
- **React Dropzone**: File upload UI
- **Axios**: HTTP client

## 📁 Project Structure

```
goshare/
├── cmd/                    # CLI commands
│   └── root.go            # Cobra CLI setup
├── internal/              # Private application code
│   └── server/
│       └── share.go       # HTTP handlers and server logic
├── frontend/              # React application
│   ├── public/           # Static assets
│   ├── src/
│   │   ├── components/   # React components
│   │   │   ├── Login.tsx
│   │   │   └── FileBrowser.tsx
│   │   ├── services/     # API services
│   │   │   └── api.ts
│   │   ├── types/        # TypeScript types
│   │   │   └── index.ts
│   │   ├── App.tsx       # Main app component
│   │   └── index.tsx     # React entry point
│   ├── tailwind.config.js
│   └── package.json
├── main.go               # Application entry point
├── go.mod                # Go dependencies
└── README.md            # Project documentation
```

## 🔧 Backend (Go)

### Core Components

#### 1. File Handler (`internal/server/share.go`)
```go
type FileHandler struct {
    rootDir   string           // Directory being served
    template  *template.Template // HTML template (fallback)
    serverURL string           // Server URL for QR codes
}
```

**Key responsibilities:**
- HTTP request routing
- File serving and directory listing
- API endpoint handling
- Authentication middleware
- CORS configuration

#### 2. API Types
```go
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
```

#### 3. HTTP Routes
- `GET /api/auth/check` - Check authentication status
- `POST /login` - User authentication
- `GET /api/files` - File listing API
- `POST /upload` - File upload
- `GET /files/*` - Direct file access
- `GET /*` - React app (catch-all)

### Security Features

1. **Path Validation**: Prevents directory traversal attacks
2. **Authentication Middleware**: Session-based auth with HTTP-only cookies
3. **CORS Configuration**: Secure cross-origin requests
4. **File Type Validation**: Secure file serving

## ⚛️ Frontend (React)

### Component Architecture

#### 1. App Component (`src/App.tsx`)
- Main application entry point
- Authentication state management
- Route between Login and FileBrowser

#### 2. Login Component (`src/components/Login.tsx`)
- Beautiful authentication interface
- Form validation and error handling
- Smooth animations with Framer Motion

#### 3. FileBrowser Component (`src/components/FileBrowser.tsx`)
- Main file browsing interface
- Real-time search functionality
- Drag & drop file upload
- Dark/light theme toggle
- File operations (view, download, navigate)

### State Management

```typescript
// Authentication state
const [isAuthenticated, setIsAuthenticated] = useState(false);

// File browser state
const [pageData, setPageData] = useState<PageData | null>(null);
const [loading, setLoading] = useState(true);
const [searchTerm, setSearchTerm] = useState('');
const [darkMode, setDarkMode] = useState(true);
```

### API Service (`src/services/api.ts`)

Provides abstracted API calls:
- `authService.login(password)` - User authentication
- `authService.checkAuth()` - Check auth status
- `fileService.getFiles(path)` - Fetch file listing
- `fileService.uploadFiles(files, directory)` - Upload files

## 📡 API Documentation

### Authentication Endpoints

#### Check Authentication Status
```http
GET /api/auth/check
Content-Type: application/json

Response: 200 OK
{
  "authenticated": true
}
```

#### User Login
```http
POST /login
Content-Type: application/x-www-form-urlencoded

password=your_password&redirect=/

Response: 303 See Other
Location: /
Set-Cookie: auth_session=authenticated; Path=/; Max-Age=86400; HttpOnly
```

### File Management Endpoints

#### Get File Listing
```http
GET /api/files?path=/some/directory
Cookie: auth_session=authenticated

Response: 200 OK
{
  "title": "GoShare - File Browser",
  "currentPath": "/some/directory",
  "parentPath": "/some",
  "hasParent": true,
  "files": [
    {
      "name": "document.pdf",
      "path": "/some/directory/document.pdf",
      "size": 1024,
      "isDir": false,
      "modTime": "2025-08-06T10:30:00Z",
      "downloadCount": 5
    }
  ],
  "serverURL": "http://192.168.1.100:8081"
}
```

#### Upload Files
```http
POST /upload
Content-Type: multipart/form-data
Cookie: auth_session=authenticated

directory=/target/directory
files=<file1>
files=<file2>

Response: 303 See Other
Location: /target/directory?uploaded=2
```

## 🔐 Authentication System

### Session Management
- **HTTP-only cookies**: Secure session storage
- **24-hour expiry**: Automatic session timeout
- **Path-based**: Applies to entire application
- **Fallback support**: Basic Auth for API clients

### Middleware Flow
```go
func applyAuthMiddleware(h http.Handler, password string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Handle login form submission
        // 2. Check session cookie
        // 3. Fallback to basic auth
        // 4. Show login form if unauthenticated
    })
}
```

## 📤 File Upload System

### Frontend (React Dropzone)
```typescript
const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: async (acceptedFiles) => {
        const fileList = acceptedFiles as unknown as FileList;
        await fileService.uploadFiles(fileList, currentPath);
    },
    noClick: true,
    noKeyboard: true
});
```

### Backend Processing
1. **Parse multipart form**: Extract files and metadata
2. **Validate destination**: Ensure path is within served directory
3. **Create directories**: Auto-create target directories
4. **Stream files**: Efficient file copying
5. **Cleanup on error**: Remove partial uploads

## 🚀 Development Setup

### Prerequisites
- Go 1.21+ 
- Node.js 18+
- npm or yarn

### One Command Setup (Recommended)
```bash
# Clone repository
git clone https://github.com/yourusername/goshare.git
cd goshare

# Start both frontend and backend with one command
./start.sh --password demo123

# Or use make shortcuts
make start-with-password
```

### Manual Development Setup
```bash
# Backend development (terminal 1)
go run main.go --port 8081 --password yourpassword /path/to/serve

# Frontend development (terminal 2)
cd frontend
npm install
npm start  # Runs on http://localhost:3000
```

### Available Scripts
- `./start.sh` - Start both servers automatically
- `./start.sh --help` - Show all options
- `make start` - Start without password
- `make start-with-password` - Start with demo password
- `make build` - Build both frontend and backend
- `make clean` - Clean all build artifacts
- `make install` - Install all dependencies

### Development Workflow
1. **Backend changes**: Just rebuild and restart Go server
2. **Frontend changes**: Use `npm start` for hot reload during development
3. **Production build**: Run `npm run build` and restart Go server

## 🏗️ Building and Deployment

### Production Build
```bash
# Build React frontend
cd frontend && npm run build && cd ..

# Build Go binary
go build -o goshare

# Run production server
./goshare --port 8080 --password securepassword /path/to/serve
```

### Docker Deployment (Optional)
```dockerfile
FROM golang:1.21 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o goshare

FROM node:18 AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/goshare .
COPY --from=frontend-builder /app/frontend/build ./frontend/build
EXPOSE 8080
CMD ["./goshare", "--port", "8080", "/data"]
```

## 🤝 Contributing

### Code Style
- **Go**: Follow standard Go formatting (`gofmt`)
- **TypeScript**: Use Prettier for formatting
- **Git**: Conventional commit messages

### Adding Features

#### Backend Features
1. Add new endpoints in `internal/server/share.go`
2. Update API types if needed
3. Add authentication if required
4. Test with curl or Postman

#### Frontend Features
1. Create/modify React components
2. Update TypeScript types in `src/types/`
3. Add API calls in `src/services/api.ts`
4. Test in browser

### Security Considerations
- Always validate file paths (prevent directory traversal)
- Sanitize user inputs
- Use HTTPS in production
- Implement rate limiting for uploads
- Regular security audits

### Performance Tips
- Use React.memo for expensive components
- Implement virtual scrolling for large file lists
- Add file caching headers
- Compress static assets
- Monitor memory usage for large uploads

---

## 🔍 Key Design Decisions

1. **Why Go + React?**
   - Go: Excellent for file serving, simple deployment, great performance
   - React: Modern UI, component reusability, great developer experience

2. **Why Session Cookies over JWT?**
   - Simpler implementation
   - HTTP-only cookies prevent XSS
   - Automatic cleanup on browser close

3. **Why Tailwind CSS?**
   - Rapid prototyping
   - Consistent design system
   - Small bundle size with purging
   - Great responsive design utilities

4. **Why Single Binary Deployment?**
   - Easy deployment (just one file)
   - No complex build processes
   - Perfect for local/LAN usage

This codebase is designed to be simple, secure, and easily extensible. The modular architecture allows for easy feature additions while maintaining clean separation of concerns.
