package models

import "time"

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
