package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

func StartServer(dir string, port int) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	fs := http.FileServer(http.Dir(absDir))
	http.Handle("/", fs)

	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	fmt.Printf("📂 Serving %s at:\n➡️  %s\n", absDir, url)

	// Generate and display QR code
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Fatalf("QR generation failed: %v", err)
	}
	fmt.Println("\n📱 Scan this QR to open:")
	fmt.Println(qr.ToSmallString(false))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
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
