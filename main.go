package main

import (
	"context"
	"file_graph/internal/embedding"
	"file_graph/internal/handlers"
	"file_graph/internal/logger"
	"file_graph/internal/mcp"
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
)

func applyMemoryProfile() {
	maxFiles := MaxFilesDefault
	batch := BatchSize
	threshold := ScanThreshold

	if LowRAM {
		batch = 300
		maxFiles = 2000
		threshold = 0.6
		fmt.Println("[Profile] Low RAM mode: batch=300, threshold=0.6")
	} else if Ram8g {
		batch = 500
		maxFiles = 3500
		threshold = 0.65
		fmt.Println("[Profile] 8GB RAM mode: batch=500, threshold=0.65")
	} else if Ram16g {
		batch = 800
		maxFiles = 4500
		threshold = 0.7
		fmt.Println("[Profile] 16GB RAM mode: batch=800, threshold=0.7")
	}

	if batch > 0 {
		BatchSize = batch
		scanner.SetScannerLimits(maxFiles, batch)
	}
	if threshold > 0 {
		ScanThreshold = threshold
		embedding.SetThreshold(threshold)
	}
}

func main() {
	startPath := flag.String("startpath", "", "Initial path to scan on startup (e.g. -startpath=C:\\folder)")
	port := flag.String("port", "8080", "Port to listen on (default: 8080)")
	mcpMode := flag.Bool("mcp", false, "Run in MCP server mode")
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
		fmt.Println("MCP Mode (Model Context Protocol):")
		fmt.Println("  -mcp         Run as an MCP server. This mode allows AI models to interact with the File Graph Visualizer via stdio using predefined tools.")
		fmt.Println("               Tools include querying node count, getting server heartbeat, retrieving node details, and highlighting nodes in the UI.")
		fmt.Println("               For detailed usage, refer to MCP_USAGE.md.")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  file_graph_server.exe -startpath=C:\\docs   Scan specific folder")
		fmt.Println("  file_graph_server.exe -mcp              Run as MCP server")
	}
	flag.Parse()

	// Check if -h was passed or no arguments
	if flag.NFlag() == 0 || flag.NArg() > 0 || len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if *mcpMode {
		logger.Println("--- START MCP ---")
		fmt.Printf("========================================\n")
		fmt.Printf("  MCP Server: Listening on Stdin/Stdout\n")
		fmt.Printf("========================================\n")
		if err := mcp.RunServer(context.Background()); err != nil {
			log.Fatalf("MCP server failed: %v", err)
		}
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
	http.HandleFunc("/api/highlighted-node", handlers.HandleHighlightedNode)
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
