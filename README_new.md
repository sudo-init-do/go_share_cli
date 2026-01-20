# GoShare

GoShare is a lightweight, high-performance file sharing utility designed for instant local and public sharing. It combines the raw power of **Go** for file streaming with a modern **React-based** user interface to provide a professional sharing experience without the complexity of cloud storage.

## Key Concepts

Unlike traditional file servers, GoShare focuses on **zero-configuration** and **ephemeral sharing**.

- **Ephemeral Security**: Every session generates a unique cryptographic token. When you stop the server, all access is wiped.
- **Single Binary Footprint**: The entire React frontend is embedded directly into the Go binary. You can move the executable to any machine and it "just works."
- **Dual-Mode Access**:
  - **Local**: Ultra-fast sharing over Wi-Fi/LAN using your internal IP.
  - **Public**: Secure sharing over the internet via an automated ngrok tunnel, bypassing NAT and firewall hurdles.

## Features

- **Professional React UI**: A responsive, high-performance interface built with React 18 and TypeScript.
- **Intelligent Theming**: Full dark and light mode support that respects your system settings and persists preferences.
- **Enhanced Security**:
  - **Randomized Sessions**: 32-byte cryptographically secure session tokens generated on every launch.
  - **HTTP-Only Cookies**: Protection against XSS and session hijacking.
  - **Traversal Protection**: Strict path validation prevents access outside the designated shared folder.
- **Dynamic File Management**:
  - **Drag & Drop**: Native-feeling upload experience with real-time progress indicators.
  - **Smart Search**: Instant, debounced filtering for large directories.
  - **Zip Downloads**: Download entire directories as compressed archives on-the-fly.
- **Insights & Tracking**:
  - **Real-time Stats**: Track download counts and last-access times for every file served.
- **Connectivity**:
  - **Instant QR**: Automatic terminal-based QR code generation for mobile devices.
  - **ngrok Integration**: One-flag internet exposure for global sharing.

## Technology Stack

- **Backend**: Go 1.21+ (Standard library, Cobra CLI, go-qrcode)
- **Frontend**: React 18, Tailwind CSS, Framer Motion, Heroicons, Vite
- **Storage**: Stateless (in-memory caching for statistics)

## Installation and Setup

### 1. Build from Source (Recommended)

To build a standalone production binary:

```bash
# Clone the repository
git clone https://github.com/sudo-init-do/goshare.git
cd goshare

# Build the frontend assets
cd frontend && npm install && npm run build && cd ..

# Compile the Go binary with embedded assets
go build -o goshare main.go
```

### 2. Quick Execution

If you have Go installed, you can run it immediately without a build step:

```bash
go run main.go [directory] --password [secret]
```

## Detailed Usage Guide

### Basic Sharing

Share the current folder publicly on your local network:

```bash
./goshare
```

### Secured Sharing

Protect your files with a password and set it on a specific port:

```bash
./goshare --dir ~/Photos --password my-secret-123 --port 9000
```

### High-Volume Uploading

If you expect to receive large files (e.g., 2GB videos), increase the upload limit:

```bash
./goshare --max-size 2048
```

### Global Internet Sharing

Expose your local folder to anyone in the world via ngrok:

```bash
./goshare --ngrok --password share-globally
```

## Command Line Reference

| Flag         | Short | Description                                   | Default       |
| :----------- | :---- | :-------------------------------------------- | :------------ |
| `--dir`      | `-d`  | The directory you wish to share               | `.` (current) |
| `--password` | `-p`  | Access password required for all visitors     | None (Open)   |
| `--port`     |       | The port the web server will listen on        | `8080`        |
| `--max-size` |       | Maximum allowed upload size in Megabytes (MB) | `100`         |
| `--ngrok`    |       | Automatically starts an ngrok tunnel          | `false`       |

## Architecture and Workflow

GoShare acts as a middleware between your file system and the web.

1. **Initialization**: The server scans the target directory and validates permissions.
2. **Session Gen**: A random token is generated. This token must be present in a cookie (`auth_session`) for any API request to succeed.
3. **Frontend Serving**: The browser requests the UI. If a password is set, the server interrupts with a custom Login page.
4. **API Interaction**: The React app communicates with the backend via a RESTful JSON API (`/api/files`, `/api/stats`).
5. **Streaming**: Files are streamed using `http.ServeContent`, allowing for "range requests" (perfect for seeking in videos).

## Troubleshooting

- **Access Denied (403)**: Ensure the user running GoShare has read permissions for the target directory.
- **Cannot Scan QR**: Some mobile browsers require the IP to be explicitly trusted if not using HTTPS. Using `--ngrok` provides a valid HTTPS link which fixes this.
- **Port Conflict**: If port 8080 is taken, use `--port [number]` to specify a free one.

---

**GoShare** - The easiest way to move files between devices. Built with code quality and security in mind.
