package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"file_graph/internal/embedding"
	"file_graph/internal/models"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	CurrentFile      string
	CurrentFileMu    sync.RWMutex
	FilesScanned     int
	FilesScannedMu   sync.RWMutex
	MaxFilesInMemory = 5000
	BatchSize        = 1000
)

func SetScannerLimits(maxFiles, batch int) {
	if maxFiles > 0 {
		MaxFilesInMemory = maxFiles
	}
	if batch > 0 {
		BatchSize = batch
	}
}

func GetScannerStats() (maxFiles, batch, scanned int) {
	FilesScannedMu.RLock()
	defer FilesScannedMu.RUnlock()
	return MaxFilesInMemory, BatchSize, FilesScanned
}

func ScanDirectory(root string) ([]models.FileNode, error) {
	// Fix common path issues
	// 1. Convert forward slashes to backslashes
	root = strings.ReplaceAll(root, "/", "\\")

	// 2. Fix duplicated drive letter (e.g., C:C:\path -> C:\path)
	if len(root) > 3 && root[1] == ':' && root[2] == ':' {
		root = root[1:]
	}

	// 3. Clean and make absolute
	root = filepath.Clean(root)
	if !filepath.IsAbs(root) {
		absRoot, err := filepath.Abs(root)
		if err == nil {
			root = absRoot
		}
	}

	var nodes []models.FileNode

	// Reset counter
	FilesScannedMu.Lock()
	FilesScanned = 0
	FilesScannedMu.Unlock()

	// Check if root exists
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip .git and other hidden folders
		if strings.HasPrefix(d.Name(), ".") && d.IsDir() && path != root {
			return filepath.SkipDir
		}

		CurrentFileMu.Lock()
		CurrentFile = path
		CurrentFileMu.Unlock()

		info, err := d.Info()
		if err != nil {
			return nil
		}

		node := models.FileNode{
			ID:         path,
			Name:       d.Name(),
			Path:       path,
			IsFolder:   d.IsDir(),
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
		}

		node.CreatedAt = node.ModifiedAt
		node.SizeLast3 = int(node.Size % 1000)

		ext := filepath.Ext(node.Name)
		base := strings.TrimSuffix(node.Name, ext)
		if len(base) > 4 {
			node.NameLast4 = base[len(base)-4:]
		} else {
			node.NameLast4 = base
		}

		if !node.IsFolder {
			node.Hash = calculateHash(path)
		}

		node.Vector = embedding.BuildVector(node)

		nodes = append(nodes, node)

		// Increment counter
		FilesScannedMu.Lock()
		FilesScanned++
		FilesScannedMu.Unlock()

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
