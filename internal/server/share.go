package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
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
