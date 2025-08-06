@echo off
REM GoShare - One Command Startup Script for Windows
REM This script starts both the React frontend and Go backend automatically

echo 🚀 Starting GoShare...
echo ================================================

REM Check if Node.js is installed
node --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Node.js is not installed. Please install Node.js 18+ first.
    pause
    exit /b 1
)

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go is not installed. Please install Go 1.21+ first.
    pause
    exit /b 1
)

REM Default values
set PORT=8081
set PASSWORD=
set DIRECTORY=.

REM Parse command line arguments (simplified)
if "%1"=="--help" (
    echo Usage: %0 [--password PASSWORD] [--port PORT] [DIRECTORY]
    echo.
    echo Examples:
    echo   %0                           # Start with no password
    echo   %0 --password demo123        # Start with password
    echo   %0 --port 9000              # Custom port
    echo   %0 C:\Users\%USERNAME%\Documents  # Share specific directory
    pause
    exit /b 0
)

if "%1"=="--password" (
    set PASSWORD=%2
    shift
    shift
)

if "%1"=="--port" (
    set PORT=%2
    shift
    shift
)

if not "%1"=="" (
    set DIRECTORY=%1
)

echo 🔨 Building Go backend...
go build -o goshare.exe .
if errorlevel 1 (
    echo ❌ Failed to build Go backend
    pause
    exit /b 1
)

echo 📦 Setting up React frontend...
cd frontend

REM Install dependencies if needed
if not exist "node_modules" (
    echo 📥 Installing frontend dependencies...
    npm install
    if errorlevel 1 (
        echo ❌ Failed to install frontend dependencies
        pause
        exit /b 1
    )
)

echo ⚛️  Starting React development server...
start /b npm start

echo 🟢 Starting Go backend...
cd ..
if "%PASSWORD%"=="" (
    start /b goshare.exe --port %PORT% "%DIRECTORY%"
) else (
    start /b goshare.exe --port %PORT% --password "%PASSWORD%" "%DIRECTORY%"
)

timeout /t 3 /nobreak >nul

echo.
echo 🎉 GoShare is now running!
echo ================================================
echo 📱 React Frontend: http://localhost:3000
echo 🔧 Go Backend API: http://localhost:%PORT%
echo.
if not "%PASSWORD%"=="" (
    echo 🔒 Password: %PASSWORD%
) else (
    echo 🔓 No password required
)
echo.
echo 📂 Serving directory: %DIRECTORY%
echo.
echo 🌐 Access your files at: http://localhost:3000
echo.
echo Press any key to stop all servers...
pause >nul

REM Kill processes
taskkill /f /im node.exe >nul 2>&1
taskkill /f /im goshare.exe >nul 2>&1
echo ✅ GoShare stopped.
