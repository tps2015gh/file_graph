package main

import (
	"file_graph/internal/embedding"
	"file_graph/internal/handlers"
	"file_graph/internal/logger"
	"file_graph/internal/scanner"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DefaultThreshold = 0.75
	DefaultBatchSize = 1000
	MaxFilesDefault  = 5000
)

var (
	ScanThreshold float64
	BatchSize     int
	LowRAM        bool
	Ram8g         bool
	Ram16g        bool
	SkipHash      bool
	MaxFiles      int
)

func applyMemoryProfile() {
	if LowRAM {
		BatchSize = 300
		ScanThreshold = 0.6
	} else if Ram8g {
		BatchSize = 500
		ScanThreshold = 0.65
	} else if Ram16g {
		BatchSize = 800
		ScanThreshold = 0.7
	} else {
		BatchSize = DefaultBatchSize
		ScanThreshold = DefaultThreshold
	}

	scanner.SetScannerLimits(MaxFilesDefault, BatchSize)
	embedding.SetThreshold(ScanThreshold)
}

func main() {
	startPath := flag.String("startpath", "", "Initial path to scan on startup (e.g. -startpath=C:\\folder)")
	port := flag.String("port", "8080", "Port to listen on (default: 8080)")
	flag.Float64Var(&ScanThreshold, "threshold", DefaultThreshold, "Similarity threshold (0.0-1.0, lower=more links)")
	flag.IntVar(&BatchSize, "batch", DefaultBatchSize, "Batch size for scanning (lower=less memory)")
	flag.BoolVar(&LowRAM, "low_ram", false, "Low memory mode (~300 batch, 0.6 threshold)")
	flag.BoolVar(&Ram8g, "ram8g", false, "8GB RAM profile (~500 batch, 0.65 threshold)")
	flag.BoolVar(&Ram16g, "ram16g", false, "16GB RAM profile (~800 batch, 0.7 threshold)")
	flag.Usage = func() {
		fmt.Println("File Graph Visualizer - Server")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  file_graph_server.exe [options]")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Memory Profiles:")
		fmt.Println("  -low_ram     Low memory mode (for systems with <2GB free RAM)")
		fmt.Println("  -ram8g       8GB RAM profile (recommended for your system)")
		fmt.Println("  -ram16g      16GB RAM profile")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  file_graph_server.exe -startpath=C:\\docs   Scan specific folder")
		fmt.Println("  file_graph_server.exe -startpath=C:\\docs -ram8g")
		fmt.Println("  file_graph_server.exe -startpath=C:\\docs -low_ram")
		fmt.Println("  file_graph_server.exe -port=9000 -startpath=C:\\myproject")
	}
	flag.Parse()

	// Check if -h was passed or no arguments
	if flag.NFlag() == 0 || flag.NArg() > 0 || len(os.Args) == 1 {
		flag.Usage()
		return
	}

	logger.Println("--- START ---")

	// Apply memory profile settings
	applyMemoryProfile()

	wd, _ := os.Getwd()
	logger.Printf("Working Dir: %s\n", wd)

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
	http.HandleFunc("/api/filelog", handlers.HandleFileLog)

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
