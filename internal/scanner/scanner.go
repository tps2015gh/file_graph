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
)

func ScanDirectory(root string) ([]models.FileNode, error) {
	root = filepath.Clean(root)
	var nodes []models.FileNode
	
	// Check if root exists
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != root {
			return filepath.SkipDir
		}

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
