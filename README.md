# GoShare CLI

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)](https://github.com/sudo-init-do/go_share_cli)
[![Release](https://img.shields.io/github/v/release/sudo-init-do/go_share_cli)](https://github.com/sudo-init-do/go_share_cli/releases)

**GoShare CLI** is a modern, lightweight command-line tool that transforms file sharing over Wi-Fi networks. Built with Go, it provides an elegant, responsive web interface for browsing, navigating, and downloading files from any device on your network—no cloud storage required.

## Quick Start

### Install GoShare (One Command)
```bash
go install github.com/sudo-init-do/goshare@latest
```

### Share Files Instantly
```bash
# Share current directory
goshare

# Share with password protection
goshare --password mysecret

# Share specific folder
goshare -d ~/Documents

# Share to the internet
goshare --ngrok
```

### Access from Any Device
- **Scan the QR code** displayed in your terminal
- **Or visit the URL** shown (e.g., `http://192.168.1.100:8080`)
- **Browse and download** files through the beautiful web interface

## Features

### Beautiful Web Interface
- **Modern Design**: Clean, responsive interface with Tailwind CSS
- **File Icons**: Visual indicators for documents, images, videos, code files
- **Smart Navigation**: Breadcrumb trails and intuitive folder browsing
- **Mobile-First**: Perfect experience on phones, tablets, and desktops
- **One-Click Downloads**: Download any file with a single click

### Instant Access
- **QR Code Generation**: Automatic QR codes for instant device connection
- **Local Network**: Share files instantly across Wi-Fi networks
- **No Setup Required**: Works immediately without configuration

### Enterprise-Grade Security
- **Password Protection**: Optional HTTP Basic Authentication
- **Path Security**: Bulletproof protection against directory traversal
- **Controlled Access**: Users can only access shared directories
- **MIME Detection**: Secure content type handling

### Global Sharing
- **Internet Access**: Expose files worldwide using ngrok integration
- **Public URLs**: Share with anyone, anywhere with secure tunneling
- **Temporary Sharing**: Perfect for one-time file distributions

## Installation Options

### Option 1: Direct Install (Recommended)
```bash
# Install latest version directly
go install github.com/sudo-init-do/goshare@latest

# Verify installation
goshare --help
```

### Option 2: Download Binary
Visit [Releases](https://github.com/sudo-init-do/go_share_cli/releases) and download the binary for your platform.

### Option 3: Build from Source
```bash
git clone https://github.com/sudo-init-do/go_share_cli.git
cd go_share_cli
go build -o goshare .
```

### Prerequisites
- **Go 1.24.4+** (for Option 1 & 3)
- **ngrok** (optional, for internet sharing)

## Usage Guide

### Basic Commands

#### Share Current Directory
```bash
goshare
```
- Shares all files in your current working directory
- Server runs on `http://localhost:8080`
- Displays QR code for instant mobile access

#### Share Specific Directory
```bash
goshare -d /path/to/your/files
goshare -d ~/Documents
goshare -d "C:\Users\John\Pictures"  # Windows
```

#### Custom Port
```bash
goshare -p 9000
```

### Advanced Features

#### Password Protection
```bash
goshare --password mysecretpassword
goshare --password "complex password 123" -d ~/Documents
```
- Adds HTTP Basic Authentication
- Users must enter the password to access files
- Password appears in username field (leave username empty)

#### Internet Sharing (ngrok)
```bash
goshare --ngrok
goshare --ngrok --password sharefiles -d ~/shared-files
```
- Exposes your files to the internet securely
- Generates public URL accessible from anywhere
- Combines with password protection for security

### Real-World Examples

#### Share Photos with Family
```bash
# Share vacation photos
goshare -d ~/Pictures/Vacation2024 --password family123
```

#### Business File Distribution
```bash
# Share presentation materials
goshare -d ~/Presentations --ngrok --password meeting2024
```

#### Developer File Sharing
```bash
# Share build artifacts
goshare -d ./dist -p 3000
```

## Command Reference

| Command | Short | Description | Example |
|---------|-------|-------------|---------|
| `--dir` | `-d` | Directory to share | `goshare -d ~/Downloads` |
| `--port` | `-p` | Server port | `goshare -p 9000` |
| `--password` | | Access password | `goshare --password secret123` |
| `--ngrok` | | Internet sharing | `goshare --ngrok` |
| `--help` | `-h` | Show help | `goshare --help` |

### Pro Tips

1. **Combine Options**: `goshare -d ~/Files -p 8080 --password secure --ngrok`
2. **Quick Access**: After running, just scan the QR code with your phone
3. **Security**: Always use passwords when sharing over the internet
4. **File Selection**: Navigate to the specific folder you want to share before running `goshare`

## Web Interface Features

GoShare's web interface provides a premium file browsing experience:

### Modern Design
- **Responsive Layout**: Automatically adapts to screen size
- **Clean Typography**: Easy-to-read fonts and proper spacing
- **Professional Styling**: Tailwind CSS for consistent, modern appearance
- **Mobile Optimized**: Touch-friendly buttons and navigation

### File Management
- **File Type Icons**: Instant visual recognition of file types
- **Size Information**: Human-readable file sizes (KB, MB, GB)
- **Modification Dates**: See when files were last changed
- **Folder Navigation**: Click folders to explore, use breadcrumbs to go back

### Supported File Types

| Category | File Types | Visual |
|----------|------------|--------|
| **Documents** | PDF, Word, Excel, PowerPoint | Document icons |
| **Images** | JPEG, PNG, GIF, SVG, WebP, BMP | Image icons |
| **Videos** | MP4, AVI, MKV, MOV, WMV | Video icons |
| **Audio** | MP3, WAV, FLAC, AAC, OGG | Audio icons |
| **Archives** | ZIP, RAR, 7Z, TAR, GZ | Archive icons |
| **Code** | Go, Python, JavaScript, C++, Java, PHP | Code icons |
| **Web** | HTML, CSS, JSON, XML | Web icons |

### Download Options
- **Direct Download**: Click file names to view/download
- **Force Download**: Use download buttons to force file download
- **Batch Access**: Navigate freely between folders
- **Secure Serving**: Proper MIME types and headers

## 🖥️ Web Interface Features

GoShare's web interface provides a premium file browsing experience:

### 🎨 **Modern Design**
- **Responsive Layout**: Automatically adapts to screen size
- **Clean Typography**: Easy-to-read fonts and proper spacing
- **Professional Styling**: Tailwind CSS for consistent, modern appearance
- **Mobile Optimized**: Touch-friendly buttons and navigation

### 📋 **File Management**
- **File Type Icons**: Instant visual recognition of file types
- **Size Information**: Human-readable file sizes (KB, MB, GB)
- **Modification Dates**: See when files were last changed
- **Folder Navigation**: Click folders to explore, use breadcrumbs to go back

### 🎯 **Supported File Types**

| Category | File Types | Icon |
|----------|------------|------|
| **Documents** | PDF, Word, Excel, PowerPoint | 📄 |
| **Images** | JPEG, PNG, GIF, SVG, WebP, BMP | 🖼️ |
| **Videos** | MP4, AVI, MKV, MOV, WMV | 🎥 |
| **Audio** | MP3, WAV, FLAC, AAC, OGG | 🎵 |
| **Archives** | ZIP, RAR, 7Z, TAR, GZ | 📦 |
| **Code** | Go, Python, JavaScript, C++, Java, PHP | 💻 |
| **Web** | HTML, CSS, JSON, XML | 🌐 |

### � **Download Options**
- **Direct Download**: Click file names to view/download
- **Force Download**: Use download buttons to force file download
- **Batch Access**: Navigate freely between folders
- **Secure Serving**: Proper MIME types and headers

## 🔧 How It Works

### Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Device │────│   GoShare CLI   │────│   File System  │
│   (Phone/PC)    │    │  HTTP Server    │    │   (Your Files)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        │ ◄─── QR Code ──────────┤                        │
        │                        │                        │
        │ ◄─── Web Interface ────┤ ◄─── File Reading ─────┤
        │                        │                        │
        │ ◄─── File Download ────┤                        │
```

### Security Model

1. **Path Sanitization**: All file paths are cleaned and validated
2. **Root Containment**: Users cannot access files outside the shared directory
3. **Authentication**: Optional password protection with HTTP Basic Auth
4. **MIME Security**: Proper content type detection prevents XSS attacks
5. **No Write Access**: Read-only file sharing (no uploads or modifications)

### Network Flow

1. **Server Start**: GoShare binds to specified port and discovers local IP
2. **QR Generation**: Creates QR code linking to `http://[LOCAL_IP]:[PORT]`
3. **Request Handling**: Custom HTTP handler processes all file requests
4. **Template Rendering**: Dynamic HTML generation for directory listings
5. **File Serving**: Secure file downloads with proper headers

### Performance Features

- **Efficient File Serving**: Uses Go's optimized `http.ServeContent`
- **Concurrent Connections**: Handles multiple users simultaneously
- **Memory Efficient**: Streams large files without loading into memory
- **Fast Directory Listing**: Optimized directory scanning and sorting

## Use Cases & Scenarios

### Business & Professional
```bash
# Share meeting presentations
goshare -d ~/Presentations --password meeting2024

# Distribute project files to team
goshare -d ./project-assets --ngrok --password team123

# Share reports with clients
goshare -d ~/Reports -p 9000 --password client2024
```

### Personal & Family
```bash
# Share vacation photos
goshare -d ~/Pictures/Vacation2024

# Transfer files between devices
goshare -d ~/Downloads

# Share documents with family
goshare -d ~/Documents --password family
```

### Development & Tech
```bash
# Share build artifacts
goshare -d ./dist -p 3000

# Distribute test files
goshare -d ./test-data --password testing

# Quick static file serving
goshare -d ./public -p 8080
```

### Education & Training
```bash
# Share course materials
goshare -d ~/Course-Materials --password students2024

# Distribute handouts during presentations
goshare -d ~/Handouts --ngrok

# Share resources in workshops
goshare -d ~/Workshop-Files
```

## Development & Contributing

### Project Structure
```
goshare/
├── cmd/
│   └── root.go          # CLI commands and flags
├── internal/
│   └── server/
│       └── share.go     # HTTP server & file handling
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
└── README.md            # Documentation
```

### Building from Source
```bash
# Clone repository
git clone https://github.com/sudo-init-do/go_share_cli.git
cd go_share_cli

# Install dependencies
go mod tidy

# Build for current platform
go build -o goshare .

# Build for multiple platforms
GOOS=windows GOARCH=amd64 go build -o goshare-windows.exe .
GOOS=darwin GOARCH=amd64 go build -o goshare-macos .
GOOS=linux GOARCH=amd64 go build -o goshare-linux .
```

### Testing
```bash
# Run tests
go test ./...

# Test with coverage
go test -cover ./...

# Test specific functionality
go test ./internal/server -v
```

### Contributing Guidelines

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature-name`
3. **Write** tests for new functionality
4. **Follow** Go best practices and `gofmt` standards
5. **Update** documentation for any changes
6. **Submit** a pull request with clear description

### Development Setup
```bash
# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Set up pre-commit hooks
go mod tidy
go fmt ./...
go vet ./...
staticcheck ./...
```

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Error: "bind: address already in use"
# Solution: Use a different port
goshare -p 8081
```

#### Can't Access from Other Devices
```bash
# Check firewall settings
# Ensure devices are on same Wi-Fi network
# Try different port: goshare -p 8080
```

#### ngrok Not Working
```bash
# Install ngrok first
# macOS: brew install ngrok
# Windows: Download from ngrok.com
# Linux: Download binary from ngrok.com
```

#### QR Code Not Scanning
- Ensure good lighting when scanning
- Try scanning from different angles
- Use the URL directly if QR doesn't work

### Performance Tips

- **Large Files**: GoShare streams files efficiently, no size limits
- **Many Files**: Directory listing is optimized for thousands of files
- **Concurrent Users**: Server handles multiple connections simultaneously
- **Network Speed**: Performance depends on your Wi-Fi network speed

## License & Legal

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) file for details.

### Dependencies
- [Cobra](https://github.com/spf13/cobra) - MIT License
- [go-qrcode](https://github.com/skip2/go-qrcode) - MIT License
- [Tailwind CSS](https://tailwindcss.com) - MIT License (CDN)
- [Font Awesome](https://fontawesome.com) - Font Awesome Free License (CDN)

## Acknowledgments

Special thanks to:
- **Go Team** for the excellent standard library
- **Open Source Community** for amazing libraries and tools
- **Contributors** who help improve GoShare
- **Users** who provide feedback and feature requests

## Support & Community

### Need Help?
1. **Check** the [FAQ](https://github.com/sudo-init-do/go_share_cli/wiki/FAQ)
2. **Search** existing [Issues](https://github.com/sudo-init-do/go_share_cli/issues)
3. **Create** a new issue with:
   - Operating system and version
   - Go version (`go version`)
   - Command you ran
   - Error message or unexpected behavior

### Community
- **GitHub Discussions**: Ask questions and share ideas
- **Issues**: Report bugs and request features
- **Pull Requests**: Contribute code improvements

### Statistics
- **Platform Support**: Windows, macOS, Linux
- **Go Version**: 1.24.4+
- **Dependencies**: Minimal (4 direct dependencies)
- **Binary Size**: ~10MB (statically compiled)

---

<div align="center">

**Star this repository if GoShare helped you!**

**Made with care by the GoShare team**

[Homepage](https://github.com/sudo-init-do/go_share_cli) • [Issues](https://github.com/sudo-init-do/go_share_cli/issues) • [Releases](https://github.com/sudo-init-do/go_share_cli/releases) • [Wiki](https://github.com/sudo-init-do/go_share_cli/wiki)

</div>
