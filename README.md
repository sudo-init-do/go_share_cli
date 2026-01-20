# GoShare CLI

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB.svg)](https://reactjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6.svg)](https://www.typescriptlang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![One Command Setup](https://img.shields.io/badge/Setup-One%20Command-brightgreen.svg)](#quick-start)

GoShare is a lightweight, high-performance file sharing utility designed for instant local and public sharing. It transforms your local directory into a modern, secure web portal with one command, allowing you to share files with anyone on your network or across the globe.

## Why GoShare?

GoShare was built to solve the "last mile" problem of file sharing—moving files between local devices (Phone to PC, Mac to Linux) without relying on slow cloud services or untrusted third-party websites.

- **Stateless & Ephemeral**: No databases. No permanent sessions. When the server stops, everything is wiped.
- **Modern Performance**: Uses Go's concurrent architecture to handle multiple high-speed streams simultaneously.
- **Embedded Frontend**: The premium React UI is compiled into the Go binary. Just one file to carry around.

## Key Features

- **Professional UI**: A beautiful, responsive interface built with React 18, Tailwind CSS, and Framer Motion.
- **System-Aware Themes**: Automatic dark and light mode switching based on your OS preferences.
- **Hardened Security**:
  - **32-byte Random Sessions**: Cryptographic tokens generated fresh on every startup.
  - **XSS Protection**: Secure HTTP-only cookies prevent token theft.
  - **Directory Isolation**: Bulletproof path validation prevents directory traversal attacks.
- **Smart Directory Management**:
  - **Dynamic Zip Generation**: Download entire folders as a single compressed `.zip` archive on-the-fly.
  - **Real-time Search**: Instant, zero-latency filtering for large file collections.
  - **Drag & Drop Uploads**: Upload files directly to your device with real-time feedback.
- **Visibility**:
  - **QR Code Access**: Scan a terminal-generated QR code for instant mobile pairing.
  - **Download Tracking**: Monitor how many times each file has been accessed in the current session.
- **Public Tunnels**: Built-in support for ngrok to share files over the global internet securely.

## Quick Start

### Building from Source

```bash
# Clone the repository
git clone https://github.com/sudo-init-do/goshare.git
cd goshare

# 1. Build the React UI
cd frontend && npm install && npm run build && cd ..

# 2. Compile the Go binary
go build -o goshare main.go
```

### Running the Server

```bash
# Basic sharing (current directory)
./goshare

# Secure sharing with a password
./goshare --password secretPass --dir ~/Downloads

# Global internet sharing
./goshare --ngrok --password globalsync
```

## Detailed Command Reference

| Flag         | Short | Description                                        | Default |
| :----------- | :---- | :------------------------------------------------- | :------ |
| `--dir`      | `-d`  | Path to the directory you want to host.            | `.`     |
| `--password` | `-p`  | Sets a password required for all web access.       | None    |
| `--port`     |       | The network port to listen on.                     | `8080`  |
| `--max-size` |       | Max upload size in Megabytes (MB).                 | `100`   |
| `--ngrok`    |       | Automatically expose your server via ngrok tunnel. | `false` |

## How It Works: Under the Hood

GoShare acts as a high-speed bridge between your OS file system and the web browser.

1. **Server Launch**: GoShare scans your target directory and validates path integrity.
2. **Session Generation**: A unique 32-byte cryptographic token is generated and stored in RAM.
3. **Authentication**: If a password is provided, GoShare wraps the React App in an authentication layer.
4. **API Layer**: The frontend communicates via a high-performance JSON API.
5. **Byte Streaming**: Using Go's `http.ServeContent`, files are served with full support for video seeking and paused downloads.

## Development

If you'd like to contribute, you can run the backend and frontend in development mode for hot-reloading:

```bash
# terminal 1: Go Backend
go run main.go --password dev

# terminal 2: React Frontend
cd frontend
npm install
npm start
```

## Troubleshooting

- **Firewall Issues**: If using Local Mode, ensure your computer's firewall allows incoming connections on the specified port.
- **QR Code Recognition**: Some iOS/Android devices require a trusted connection. Using `--ngrok` provides a valid HTTPS certificate which improves mobile compatibility.
- **Large Directories**: For folders with 10k+ files, the search feature remains responsive due to debounced filtering in the React layer.

---

**Built with ❤️ by the GoShare Team.** MIT Licensed.
