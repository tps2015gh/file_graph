package embedding

import (
	"encoding/hex"
	"file_graph/internal/models"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	threshold   float64 = 0.75
	thresholdMu sync.RWMutex
	numberRegex = regexp.MustCompile(`\d+`)
)

func SetThreshold(t float64) {
	thresholdMu.Lock()
	defer thresholdMu.Unlock()
	if t > 0 && t <= 1.0 {
		threshold = t
	}
}

func GetThreshold() float64 {
	thresholdMu.RLock()
	defer thresholdMu.RUnlock()
	return threshold
}

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
	thresh := GetThreshold()

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {

			// Allow folder-to-file links in same folder
			dirI := filepath.Dir(nodes[i].Path)
			dirJ := filepath.Dir(nodes[j].Path)
			sameFolder := dirI == dirJ

			// Skip cross-folder folder-to-file links
			if (nodes[i].IsFolder != nodes[j].IsFolder) && !sameFolder {
				continue
			}

			sim := cosineSimilarity(nodes[i].Vector, nodes[j].Vector)

			// Folder proximity bonus - same folder gets strong bonus
			if sameFolder {
				sim += 0.4 // Strong bonus for same folder
			}

			// Filename prefix match bonus (works cross-folder)
			prefixLen := commonPrefixLen(nodes[i].Name, nodes[j].Name)
			if prefixLen >= 3 {
				sim += float64(prefixLen) * 0.05
			}

			// Filename suffix match (before extension)
			extI := filepath.Ext(nodes[i].Name)
			extJ := filepath.Ext(nodes[j].Name)
			baseI := nodes[i].Name[:len(nodes[i].Name)-len(extI)]
			baseJ := nodes[j].Name[:len(nodes[j].Name)-len(extJ)]
			suffixLen := commonSuffixLen(baseI, baseJ)
			if suffixLen >= 3 {
				sim += float64(suffixLen) * 0.03
			}

			// Number proximity bonus - files with similar numbers (e.g., file1.txt, file2.txt)
			numBonus := numberProximity(nodes[i].Name, nodes[j].Name)
			sim += numBonus

			// Extension grouping bonus
			if extI == extJ && extI != "" {
				sim += 0.15
			}

			// Same name last 4 chars bonus
			if nodes[i].NameLast4 == nodes[j].NameLast4 && nodes[i].NameLast4 != "" {
				sim += 0.1
			}

			// Same size last 3 digits bonus
			if nodes[i].SizeLast3 == nodes[j].SizeLast3 && nodes[i].Size > 0 {
				sim += 0.05
			}

			if sim > thresh {
				links = append(links, models.Relation{
					Source:     nodes[i].ID,
					Target:     nodes[j].ID,
					Similarity: math.Min(sim, 1.0), // Cap at 1.0
				})
			}
		}
	}
	return links
}

func numberProximity(name1, name2 string) float64 {
	nums1 := extractNumbers(name1)
	nums2 := extractNumbers(name2)

	if len(nums1) == 0 || len(nums2) == 0 {
		return 0
	}

	// Compare each number pair
	for _, n1 := range nums1 {
		for _, n2 := range nums2 {
			diff := float64(abs(n1 - n2))
			if diff == 0 && n1 > 0 {
				// Same number - reduced bonus
				return 0.15
			} else if diff == 1 {
				// Adjacent numbers (1,2) or (2,3) - small bonus
				return 0.08
			} else if diff <= 3 {
				// Very close numbers within 3 - tiny bonus
				return 0.03
			}
		}
	}
	return 0
}

func extractNumbers(s string) []int {
	matches := numberRegex.FindAllString(s, -1)
	var nums []int
	for _, m := range matches {
		if n, err := strconv.Atoi(m); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func commonPrefixLen(s1, s2 string) int {
	i := 0
	for ; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			break
		}
	}
	return i
}

func commonSuffixLen(s1, s2 string) int {
	i, j := len(s1)-1, len(s2)-1
	count := 0
	for ; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if s1[i] != s2[j] {
			break
		}
		count++
	}
	return count
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
