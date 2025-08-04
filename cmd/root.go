package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"github.com/sudo-init-do/goshare/internal/server"
)

var (
	dir      string
	port     int
	password string
	useNgrok bool
)

var rootCmd = &cobra.Command{
	Use:   "goshare",
	Short: "Easily share local files over Wi‑Fi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting goshare on port %d serving directory: %s\n", port, dir)
		if useNgrok {
			startNgrokTunnel(dir, port, password)
			return
		}
		server.StartServer(dir, port, password)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "Directory to share")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "", "", "Optional password to protect access (Basic Auth)")
	rootCmd.PersistentFlags().BoolVar(&useNgrok, "ngrok", false, "Expose server to the internet using ngrok")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startNgrokTunnel(dir string, port int, password string) {
	// Start the local server concurrently (prints local IP + QR)
	go server.StartServer(dir, port, password)

	fmt.Println("📡 Launching ngrok tunnel...")

	// Run ngrok silently (no logs to stdout/stderr)
	cmd := exec.Command("ngrok", "http", fmt.Sprintf("%d", port))

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Failed to start ngrok:", err)
		os.Exit(1)
	}

	// Poll ngrok's local API for the public URL
	publicURL := waitForNgrokURL(30 * time.Second) // longer timeout for reliability
	if publicURL == "" {
		fmt.Println("⚠️  Could not detect ngrok public URL. Check http://127.0.0.1:4040")
	} else {
		fmt.Println("\n🌍 Public URL (ngrok):", publicURL)
		if qr, err := qrcode.New(publicURL, qrcode.Medium); err == nil {
			fmt.Println("\n📱 Scan this QR (ngrok):")
			fmt.Println(qr.ToSmallString(false))
		} else {
			fmt.Println("⚠️  Could not generate QR for ngrok URL:", err)
		}
	}

	// Keep ngrok process alive
	if err := cmd.Wait(); err != nil {
		fmt.Println("ngrok exited with error:", err)
	}
}

func waitForNgrokURL(timeout time.Duration) string {
	type tunnel struct {
		PublicURL string `json:"public_url"`
	}
	type tunnelsResp struct {
		Tunnels []tunnel `json:"tunnels"`
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://127.0.0.1:4040/api/tunnels")
		if err == nil && resp != nil && resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			var tr tunnelsResp
			if json.Unmarshal(body, &tr) == nil {
				// Prefer HTTPS
				for _, t := range tr.Tunnels {
					if strings.HasPrefix(t.PublicURL, "https://") {
						return t.PublicURL
					}
				}
				// Fallback: any URL
				for _, t := range tr.Tunnels {
					if t.PublicURL != "" {
						return t.PublicURL
					}
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ""
}
