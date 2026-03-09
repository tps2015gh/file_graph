package main

import (
	"file_graph/internal/handlers"
	"file_graph/internal/logger"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	logger.Println("--- START ---")

	startPath := flag.String("startpath", "", "Initial path to scan on startup (e.g. -startpath=C:\\folder)")
	port := flag.String("port", "8080", "Port to listen on (default: 8080)")
	flag.Usage = func() {
		fmt.Println("File Graph Visualizer - Server")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  file_graph_server.exe [options]")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  file_graph_server.exe                      Run with current folder")
		fmt.Println("  file_graph_server.exe -startpath=C:\\docs   Scan specific folder")
		fmt.Println("  file_graph_server.exe -port=9000           Use custom port")
	}
	flag.Parse()

	wd, _ := os.Getwd()
	logger.Printf("Working Dir: %s\n", wd)

	// Show help if no startpath provided - but don't exit, just warn
	if *startPath == "" {
		fmt.Println("")
		fmt.Println(">>> Tip: To scan a specific folder, use: -startpath=<folder>")
		fmt.Println(">>> Example: -startpath=" + wd)
		fmt.Println("")
	}

	if *startPath != "" {
		absPath, err := filepath.Abs(*startPath)
		if err == nil {
			*startPath = absPath
		}
		logger.Printf("StartPath: %s\n", *startPath)
		handlers.SetStartPath(*startPath)
	}

	http.HandleFunc("/", handlers.ServeIndex)
	http.HandleFunc("/favicon.ico", handlers.ServeFavicon)
	http.HandleFunc("/api/startpath", handlers.HandleGetStartPath)
	http.HandleFunc("/api/scan", handlers.HandleScan)
	http.HandleFunc("/api/open", handlers.HandleOpen)
	http.HandleFunc("/api/log", handlers.HandleClientLog)
	http.HandleFunc("/api/heartbeat", handlers.HandleHeartbeat)
	http.HandleFunc("/api/progress", handlers.HandleProgress)
	http.HandleFunc("/api/shutdown", handlers.HandleShutdown)
	http.HandleFunc("/api/kill", handlers.HandleHardExit)

	webURL := fmt.Sprintf("http://localhost:%s", *port)
	apiURL := fmt.Sprintf("http://localhost:%s/api", *port)

	fmt.Printf("========================================\n")
	fmt.Printf("  Web App : %s\n", webURL)
	fmt.Printf("  API     : %s\n", apiURL)
	fmt.Printf("========================================\n")
	logger.Printf("Web App: %s\n", webURL)
	logger.Printf("API: %s\n", apiURL)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
