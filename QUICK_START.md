# GoShare - One Command Demo

## Quick Test

Run this to test your setup:

```bash
# Test the startup script
./start.sh --help

# Start with demo password
./start.sh --password demo123

# Or use make
make start-with-password
```

## What Happens When You Run It

1. ✅ **Checks Dependencies**: Verifies Node.js and Go are installed
2. ✅ **Builds Backend**: Compiles the Go server automatically  
3. ✅ **Installs Frontend**: Runs `npm install` if needed
4. ✅ **Starts React Dev Server**: Hot reload on `localhost:3000`
5. ✅ **Starts Go API Server**: Backend API on `localhost:8081`
6. ✅ **Shows URLs**: Clear instructions on where to access your app

## Access Your App

- **Main App**: http://localhost:3000 (React frontend)
- **Direct API**: http://localhost:8081 (Go backend)
- **Password**: Whatever you set (or no password)

## Stop Everything

Just press `Ctrl+C` and both servers will stop automatically!

## Super Easy Commands

```bash
make start                    # No password required
make start-with-password      # Demo password: demo123  
make start-port              # Runs on port 9000
make help                    # See all options
```

That's it! One command, everything works. 🚀
