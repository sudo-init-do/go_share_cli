package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

func StartServer(dir string, port int, password string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// File server
	fs := http.FileServer(http.Dir(absDir))
	// Wrap with auth middleware if password provided
	http.Handle("/", applyAuthMiddleware(fs, password))

	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	fmt.Printf("📂 Serving %s at:\n➡️  %s\n", absDir, url)

	// Generate and display local QR code
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Fatalf("QR generation failed: %v", err)
	}
	fmt.Println("\n📱 Scan this QR to open (local):")
	fmt.Println(qr.ToSmallString(false))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func applyAuthMiddleware(h http.Handler, password string) http.Handler {
	if password == "" {
		return h // no protection
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := r.BasicAuth()
		if !ok || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="go_share_cli"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "localhost"
}
