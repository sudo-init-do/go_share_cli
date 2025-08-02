package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sudo-init-do/goshare/internal/server"
)

var (
	dir  string
	port int
)

var rootCmd = &cobra.Command{
	Use:   "goshare",
	Short: "Easily share local files over Wi-Fi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting goshare on port %d serving directory: %s\n", port, dir)
		server.StartServer(dir, port)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "Directory to share")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
