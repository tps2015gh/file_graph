package handlers

import (
	"encoding/json"
	"file_graph/internal/embedding"
	"file_graph/internal/logger"
	"file_graph/internal/models"
	"file_graph/internal/scanner"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func HandleScan(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "."
	}
	logger.Printf("SCAN: %s\n", dir)

	nodes, err := scanner.ScanDirectory(dir)
	if err != nil {
		fmt.Printf("SCAN ERROR for %s: %v\n", dir, err)
		logger.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Printf("RESULT: %d nodes\n", len(nodes))

	links := embedding.CalculateRelations(nodes)

	result := models.ScanResult{
		Nodes: nodes,
		Links: links,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func HandleOpen(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}
	logger.Printf("OPEN: %s\n", path)

	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}
	
	exec.Command("explorer", "/select,", path).Run()
	w.WriteHeader(http.StatusOK)
}

func HandleClientLog(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg != "" {
		logger.Printf("UI: %s\n", msg)
	}
	w.WriteHeader(http.StatusOK)
}

func HandleShutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shutdown requested...")
	logger.Println("EVENT: SHUTDOWN")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down..."))
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
