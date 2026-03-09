package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// FileNode represents a file or folder in the graph
type FileNode struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	IsFolder     bool      `json:"is_folder"`
	Size         int64     `json:"size"`
	SizeLast3    int       `json:"size_last_3"`
	NameLast4    string    `json:"name_last_4"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
	Hash         string    `json:"hash"`
	Vector       []float64 `json:"vector"` // Normalized vector for similarity
}

// Relation represents a link between nodes
type Relation struct {
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Similarity float64 `json:"similarity"`
}

type ScanResult struct {
	Nodes []FileNode `json:"nodes"`
	Links []Relation `json:"links"`
}

var appLogger *log.Logger

func init() {
	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not create log file: %v\n", err)
		return
	}
	// Concise format: date time message
	appLogger = log.New(logFile, "", log.Ldate|log.Ltime)
}

func main() {
	if appLogger != nil {
		appLogger.Println("--- START ---")
	}
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/scan", handleScan)
	http.HandleFunc("/api/open", handleOpen)
	http.HandleFunc("/api/log", handleClientLog)
	http.HandleFunc("/api/shutdown", handleShutdown)

	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	if appLogger != nil {
		appLogger.Printf("Server: http://localhost:%s\n", port)
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleClientLog(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg != "" && appLogger != nil {
		appLogger.Printf("UI: %s\n", msg)
	}
	w.WriteHeader(http.StatusOK)
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shutdown requested...")
	if appLogger != nil {
		appLogger.Println("EVENT: SHUTDOWN")
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down..."))
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func handleOpen(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}
	if appLogger != nil {
		appLogger.Printf("OPEN: %s\n", path)
	}

	// For Windows: open the path in Explorer and select the file
	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}
	
	exec.Command("explorer", "/select,", path).Run()
	w.WriteHeader(http.StatusOK)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func handleScan(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "."
	}
	if appLogger != nil {
		appLogger.Printf("SCAN: %s\n", dir)
	}

	nodes, err := scanDirectory(dir)
	if err != nil {
		fmt.Printf("SCAN ERROR for %s: %v\n", dir, err)
		if appLogger != nil {
			appLogger.Printf("ERROR: %v\n", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if appLogger != nil {
		appLogger.Printf("RESULT: %d nodes\n", len(nodes))
	}

	links := calculateRelations(nodes)

	result := ScanResult{
		Nodes: nodes,
		Links: links,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func scanDirectory(root string) ([]FileNode, error) {
	var nodes []FileNode
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip hidden folders
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != root {
			return filepath.SkipDir
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		node := FileNode{
			ID:         path,
			Name:       d.Name(),
			Path:       path,
			IsFolder:   d.IsDir(),
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
		}

		// Get creation time (os-specific, simplified here)
		node.CreatedAt = node.ModifiedAt // Fallback

		// Last 3 digits of size
		node.SizeLast3 = int(node.Size % 1000)

		// Last 4 chars of name (not extension)
		ext := filepath.Ext(node.Name)
		base := strings.TrimSuffix(node.Name, ext)
		if len(base) > 4 {
			node.NameLast4 = base[len(base)-4:]
		} else {
			node.NameLast4 = base
		}

		// Hash for files
		if !node.IsFolder {
			node.Hash = calculateHash(path)
		}

		// Build a simple "embedding" vector for similarity comparison
		node.Vector = buildVector(node)

		nodes = append(nodes, node)
		return nil
	})

	return nodes, err
}

func calculateHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func buildVector(n FileNode) []float64 {
	// Simple normalization/encoding for similarity
	// 1: Size (log scale)
	// 2: SizeLast3 (normalized 0-1)
	// 3: NameHash (first byte of name)
	// 4: ModTime (unix timestamp normalized)
	
	v := make([]float64, 4)
	if n.Size > 0 {
		v[0] = math.Log10(float64(n.Size))
	}
	v[1] = float64(n.SizeLast3) / 1000.0
	
	if len(n.Name) > 0 {
		v[2] = float64(n.Name[0]) / 255.0
	}
	
	v[3] = float64(n.ModifiedAt.Unix() % 86400) / 86400.0 // Daily cycle similarity
	
	return v
}

func calculateRelations(nodes []FileNode) []Relation {
	var links []Relation
	// To avoid O(N^2) explosion on large dirs, we'll only link strongest matches
	// or use a threshold.
	
	threshold := 0.85
	
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			sim := cosineSimilarity(nodes[i].Vector, nodes[j].Vector)
			
			// Also check explicit attributes
			if nodes[i].NameLast4 == nodes[j].NameLast4 && nodes[i].NameLast4 != "" {
				sim += 0.1
			}
			if nodes[i].SizeLast3 == nodes[j].SizeLast3 && nodes[i].Size > 0 {
				sim += 0.05
			}
			
			if sim > threshold {
				links = append(links, Relation{
					Source:     nodes[i].ID,
					Target:     nodes[j].ID,
					Similarity: sim,
				})
			}
		}
	}
	return links
}

func cosineSimilarity(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return 0
	}
	var dot, mag1, mag2 float64
	for i := range v1 {
		dot += v1[i] * v2[i]
		mag1 += v1[i] * v1[i]
		mag2 += v2[i] * v2[i]
	}
	if mag1 == 0 || mag2 == 0 {
		return 0
	}
	return dot / (math.Sqrt(mag1) * math.Sqrt(mag2))
}
