#!/bin/bash

# GoShare - One Command Startup Script
# This script starts both the React frontend and Go backend automatically

set -e  # Exit on any error

echo "🚀 Starting GoShare..."
echo "================================================"

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Default values
PORT=8081
PASSWORD=""
DIRECTORY="."

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        --password)
            PASSWORD="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS] [DIRECTORY]"
            echo ""
            echo "Options:"
            echo "  -p, --port PORT     Server port (default: 8081)"
            echo "  --password PASS     Password for access (optional)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Start with no password"
            echo "  $0 --password demo123                # Start with password"
            echo "  $0 --port 9000 --password secret     # Custom port and password"
            echo "  $0 ~/Documents                       # Share specific directory"
            exit 0
            ;;
        -*)
            echo "Unknown option $1"
            exit 1
            ;;
        *)
            DIRECTORY="$1"
            shift
            ;;
    esac
done

# Function to cleanup processes on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down GoShare..."
    if [[ ! -z $FRONTEND_PID ]]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    if [[ ! -z $BACKEND_PID ]]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    echo "✅ GoShare stopped."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Build the Go backend
echo "🔨 Building Go backend..."
go build -o goshare . || {
    echo "❌ Failed to build Go backend"
    exit 1
}

# Install and start React frontend
echo "📦 Setting up React frontend..."
cd frontend

# Install dependencies if node_modules doesn't exist or package.json is newer
if [[ ! -d "node_modules" ]] || [[ "package.json" -nt "node_modules" ]]; then
    echo "📥 Installing frontend dependencies..."
    npm install || {
        echo "❌ Failed to install frontend dependencies"
        exit 1
    }
fi

echo "⚛️  Starting React development server..."
npm start &
FRONTEND_PID=$!

# Give React time to start
sleep 3

# Go back to root directory
cd ..

# Start Go backend
echo "🟢 Starting Go backend..."
if [[ -z "$PASSWORD" ]]; then
    ./goshare --port $PORT "$DIRECTORY" &
else
    ./goshare --port $PORT --password "$PASSWORD" "$DIRECTORY" &
fi
BACKEND_PID=$!

# Wait a moment for servers to start
sleep 2

echo ""
echo "🎉 GoShare is now running!"
echo "================================================"
echo "📱 React Frontend: http://localhost:3000"
echo "🔧 Go Backend API: http://localhost:$PORT"
echo ""
if [[ ! -z "$PASSWORD" ]]; then
    echo "🔒 Password: $PASSWORD"
else
    echo "🔓 No password required"
fi
echo ""
echo "📂 Serving directory: $(realpath "$DIRECTORY")"
echo ""
echo "🌐 Access your files at: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all servers..."

# Wait for processes
wait $FRONTEND_PID $BACKEND_PID
