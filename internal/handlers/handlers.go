package handlers

import (
	"encoding/json"
	"file_graph/internal/embedding"
	"file_graph/internal/logger"
	"file_graph/internal/models"
	"file_graph/internal/scanner"
	"file_graph/internal/state"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// HandleHighlightedNode checks for a node ID set by the MCP and returns it.
// The state is cleared after being read to prevent re-highlighting.
func HandleHighlightedNode(w http.ResponseWriter, r *http.Request) {
	nodeID, found := state.GlobalState.GetHighlightedNodeID()

	w.Header().Set("Content-Type", "application/json")
	if !found {
		// Return an empty object if no node is highlighted
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"highlightNodeId": nodeID,
	})
}

var StartPath string

func SetStartPath(path string) {
	StartPath = path
}

func GetStartPath() string {
	return StartPath
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func HandleGetStartPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": StartPath,
	})
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

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(result); err != nil {
		logger.Printf("JSON Encode Error: %v\n", err)
	}
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

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"time":   time.Now().Format("15:04:05"),
		"status": "alive",
	})
}

func HandleProgress(w http.ResponseWriter, r *http.Request) {
	scanner.CurrentFileMu.RLock()
	defer scanner.CurrentFileMu.RUnlock()

	maxFiles, batchSize, scanned := scanner.GetScannerStats()
	threshold := embedding.GetThreshold()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path":      scanner.CurrentFile,
		"max_files": maxFiles,
		"batch":     batchSize,
		"threshold": threshold,
		"scanned":   scanned,
	})
}

func HandleShutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shutdown requested (Restart)...")
	logger.Println("EVENT: RESTART")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down..."))
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0) // Exit code 0 signals RESTART to the batch file
	}()
}

func HandleHardExit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hard Exit requested (Kill)...")
	logger.Println("EVENT: KILL")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Killing server..."))
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(1) // Exit code 1 signals STOP to the batch file
	}()
}
