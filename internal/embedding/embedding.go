package embedding

import (
	"encoding/hex"
	"file_graph/internal/models"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func BuildVector(n models.FileNode) []float64 {
	// 25-Dimensional Vector Construction
	// 0: Size (log scale)
	// 1: SizeLast3 (normalized 0-1)
	// 2: Folder Depth
	// 3: Name Length (normalized)
	// 4: Extension Hash
	// 5: Count '.'
	// 6: Count '_'
	// 7: Count Digits
	// 8: ModTime: Hour (0-1)
	// 9: ModTime: Weekday (0-1)
	// 10: ModTime: Month (0-1)
	// 11-15: First 5 chars of name (normalized ascii)
	// 16-25: First 10 bytes of SHA256 Hash (normalized)

	v := make([]float64, 26) // 26 dimensions total

	// 0. Size Log10
	if n.Size > 0 {
		v[0] = math.Log10(float64(n.Size))
	}

	// 1. Size Mod 1000
	v[1] = float64(n.SizeLast3) / 1000.0

	// 2. Folder Depth (Approximate by separator count)
	v[2] = float64(strings.Count(n.Path, string(os.PathSeparator))) / 10.0

	// 3. Name Length (cap at 100)
	v[3] = math.Min(float64(len(n.Name)), 100.0) / 100.0

	// 4. Extension Hash
	ext := filepath.Ext(n.Name)
	if len(ext) > 0 {
		// Simple sum hash of extension chars
		sum := 0
		for _, c := range ext {
			sum += int(c)
		}
		v[4] = float64(sum%255) / 255.0
	}

	// 5, 6, 7. Char Stats
	v[5] = float64(strings.Count(n.Name, ".")) / 5.0
	v[6] = float64(strings.Count(n.Name, "_")) / 5.0
	digitCount := 0
	for _, c := range n.Name {
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}
	v[7] = float64(digitCount) / 10.0

	// 8, 9, 10. Time Attributes
	v[8] = float64(n.ModifiedAt.Hour()) / 24.0
	v[9] = float64(n.ModifiedAt.Weekday()) / 7.0
	v[10] = float64(n.ModifiedAt.Month()) / 12.0

	// 11-15. First 5 chars of name
	nameRunes := []rune(n.Name)
	for i := 0; i < 5; i++ {
		if i < len(nameRunes) {
			v[11+i] = float64(nameRunes[i]) / 255.0
		}
	}

	// 16-25. Hash Segments (10 dimensions)
	if n.Hash != "" {
		hashBytes, err := hex.DecodeString(n.Hash)
		if err == nil {
			for i := 0; i < 10; i++ {
				if i < len(hashBytes) {
					v[16+i] = float64(hashBytes[i]) / 255.0
				}
			}
		}
	}

	return v
}

func CalculateRelations(nodes []models.FileNode) []models.Relation {
	var links []models.Relation
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
				links = append(links, models.Relation{
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
