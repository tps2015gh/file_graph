package embedding

import (
	"encoding/hex"
	"file_graph/internal/models"
	"math"
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

	// 0. Size (log scale)
	if n.Size > 0 {
		v[0] = math.Log(float64(n.Size)+1) / 20.0 // Normalize
	}

	// 1. SizeLast3 (normalized 0-1)
	if n.Size > 0 {
		v[1] = float64(n.SizeLast3) / 999.0
	}

	// 2. Folder Depth (calculate from path)
	depth := strings.Count(n.Path, string(filepath.Separator))
	v[2] = float64(depth) / 20.0

	// 8-10. ModifiedAt components
	hour := n.ModifiedAt.Hour()
	weekday := int(n.ModifiedAt.Weekday())
	month := int(n.ModifiedAt.Month())
	v[8] = float64(hour) / 23.0
	v[9] = float64(weekday) / 6.0
	v[10] = float64(month) / 11.0

	// 11-15. First 5 chars of name
	for i := 0; i < 5; i++ {
		if i < len(n.Name) {
			v[11+i] = float64(n.Name[i]) / 255.0
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
	maxNeighbors := 5

	if len(nodes) == 0 {
		return links
	}

	// For large datasets, use hash-based spatial bucketing
	if len(nodes) > 1000 {
		return calculateRelationsBucketed(nodes, thresh, maxNeighbors)
	}

	// Original brute-force for small datasets
	neighborMatrix := make([][]models.Relation, len(nodes))

	for i := 0; i < len(nodes); i++ {
		var bestNeighbors []models.Relation

		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}

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
				sim += 0.4
			}

			// Filename prefix match bonus
			prefixLen := commonPrefixLen(nodes[i].Name, nodes[j].Name)
			if prefixLen >= 3 {
				sim += float64(prefixLen) * 0.05
			}

			// Filename suffix match
			extI := filepath.Ext(nodes[i].Name)
			extJ := filepath.Ext(nodes[j].Name)
			baseI := nodes[i].Name[:len(nodes[i].Name)-len(extI)]
			baseJ := nodes[j].Name[:len(nodes[j].Name)-len(extJ)]
			suffixLen := commonSuffixLen(baseI, baseJ)
			if suffixLen >= 3 {
				sim += float64(suffixLen) * 0.03
			}

			// Number proximity bonus
			sim += numberProximity(nodes[i].Name, nodes[j].Name)

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
				relation := models.Relation{
					Source:     nodes[i].ID,
					Target:     nodes[j].ID,
					Similarity: math.Min(sim, 1.0),
				}

				// Insert sorted by similarity
				bestNeighbors = insertSorted(bestNeighbors, relation, maxNeighbors)
			}
		}
		neighborMatrix[i] = bestNeighbors
	}

	// Collect all relations, deduplicate
	linkMap := make(map[string]models.Relation)
	for _, neighbors := range neighborMatrix {
		for _, rel := range neighbors {
			key := relationKey(rel.Source, rel.Target)
			linkMap[key] = rel
		}
	}

	// Convert map to slice
	for _, rel := range linkMap {
		links = append(links, rel)
	}
	return links
}

// Hash-based spatial bucketing for large datasets
func calculateRelationsBucketed(nodes []models.FileNode, thresh float64, maxNeighbors int) []models.Relation {
	if len(nodes) == 0 {
		return []models.Relation{}
	}

	// Create hash buckets - use first 4 bytes (8 hex chars) of SHA-256
	buckets := make(map[string][]int)
	for i, node := range nodes {
		if node.Hash != "" && len(node.Hash) >= 8 {
			bucketKey := node.Hash[:8] // First 4 bytes
			buckets[bucketKey] = append(buckets[bucketKey], i)
		}
	}

	// Also create folder-based buckets for folder proximity
	folderBuckets := make(map[string][]int)
	for i, node := range nodes {
		folder := filepath.Dir(node.Path)
		folderBuckets[folder] = append(folderBuckets[folder], i)
	}

	var links []models.Relation
	processedPairs := make(map[string]bool)

	// Compare within each bucket
	for _, bucket := range buckets {
		if len(bucket) > 1 {
			bucketLinks := compareBucket(nodes, bucket, thresh, maxNeighbors, processedPairs)
			links = append(links, bucketLinks...)
		}
	}

	// Compare within each folder bucket
	for _, bucket := range folderBuckets {
		if len(bucket) > 1 {
			bucketLinks := compareBucket(nodes, bucket, thresh, maxNeighbors, processedPairs)
			links = append(links, bucketLinks...)
		}
	}

	// For small buckets (<3 files), also check adjacent buckets
	for bucketKey, bucket := range buckets {
		if len(bucket) < 3 {
			adjacentLinks := compareToAdjacentBuckets(nodes, bucket, buckets, bucketKey, thresh, maxNeighbors, processedPairs)
			links = append(links, adjacentLinks...)
		}
	}

	return links
}

func compareBucket(nodes []models.FileNode, bucket []int, thresh float64, maxNeighbors int, processedPairs map[string]bool) []models.Relation {
	var links []models.Relation

	for _, i := range bucket {
		var bestNeighbors []models.Relation

		for _, j := range bucket {
			if i == j {
				continue
			}

			key := relationKey(nodes[i].ID, nodes[j].ID)
			if processedPairs[key] {
				continue
			}
			processedPairs[key] = true

			sim := cosineSimilarity(nodes[i].Vector, nodes[j].Vector)
			sim = applyBonuses(nodes[i], nodes[j], sim)

			if sim > thresh {
				relation := models.Relation{
					Source:     nodes[i].ID,
					Target:     nodes[j].ID,
					Similarity: math.Min(sim, 1.0),
				}
				bestNeighbors = insertSorted(bestNeighbors, relation, maxNeighbors)
			}
		}

		links = append(links, bestNeighbors...)
	}

	return links
}

func compareToAdjacentBuckets(nodes []models.FileNode, currentBucket []int, buckets map[string][]int, currentBucketKey string, thresh float64, maxNeighbors int, processedPairs map[string]bool) []models.Relation {
	var links []models.Relation

	// Get adjacent buckets by modifying the bucket key slightly
	adjacentBuckets := []string{
		currentBucketKey[:7] + "0", // Change last character
		currentBucketKey[:7] + "1",
		currentBucketKey[:6] + "00", // Change last 2 characters
	}

	for _, adjKey := range adjacentBuckets {
		if adjBucket, exists := buckets[adjKey]; exists {
			for _, i := range currentBucket {
				var bestNeighbors []models.Relation

				for _, j := range adjBucket {
					if i == j {
						continue
					}

					key := relationKey(nodes[i].ID, nodes[j].ID)
					if processedPairs[key] {
						continue
					}
					processedPairs[key] = true

					sim := cosineSimilarity(nodes[i].Vector, nodes[j].Vector)
					sim = applyBonuses(nodes[i], nodes[j], sim)

					if sim > thresh {
						relation := models.Relation{
							Source:     nodes[i].ID,
							Target:     nodes[j].ID,
							Similarity: math.Min(sim, 1.0),
						}
						bestNeighbors = insertSorted(bestNeighbors, relation, maxNeighbors)
					}
				}

				links = append(links, bestNeighbors...)
			}
		}
	}

	return links
}

func applyBonuses(a, b models.FileNode, similarity float64) float64 {
	sim := similarity

	// Allow folder-to-file links in same folder
	dirA := filepath.Dir(a.Path)
	dirB := filepath.Dir(b.Path)
	sameFolder := dirA == dirB

	// Skip cross-folder folder-to-file links
	if (a.IsFolder != b.IsFolder) && !sameFolder {
		return sim
	}

	// Folder proximity bonus - same folder gets strong bonus
	if sameFolder {
		sim += 0.4
	}

	// Filename prefix match bonus
	prefixLen := commonPrefixLen(a.Name, b.Name)
	if prefixLen >= 3 {
		sim += float64(prefixLen) * 0.05
	}

	// Filename suffix match
	extA := filepath.Ext(a.Name)
	extB := filepath.Ext(b.Name)
	baseA := a.Name[:len(a.Name)-len(extA)]
	baseB := b.Name[:len(b.Name)-len(extB)]
	suffixLen := commonSuffixLen(baseA, baseB)
	if suffixLen >= 3 {
		sim += float64(suffixLen) * 0.03
	}

	// Number proximity bonus
	sim += numberProximity(a.Name, b.Name)

	// Extension grouping bonus
	if extA == extB && extA != "" {
		sim += 0.15
	}

	// Same name last 4 chars bonus
	if a.NameLast4 == b.NameLast4 && a.NameLast4 != "" {
		sim += 0.1
	}

	// Same size last 3 digits bonus
	if a.SizeLast3 == b.SizeLast3 && a.Size > 0 {
		sim += 0.05
	}

	return sim
}

func relationKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}

func insertSorted(neighbors []models.Relation, rel models.Relation, maxSize int) []models.Relation {
	// Find insertion position
	idx := len(neighbors)
	for i, n := range neighbors {
		if rel.Similarity > n.Similarity {
			idx = i
			break
		}
	}

	// Insert
	if idx < maxSize {
		neighbors = append(neighbors, models.Relation{})
		copy(neighbors[idx+1:], neighbors[idx:])
		neighbors[idx] = rel

		// Keep only top maxSize
		if len(neighbors) > maxSize {
			neighbors = neighbors[:maxSize]
		}
	}
	return neighbors
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
	for i < len(s1) && i < len(s2) && s1[i] == s2[i] {
		i++
	}
	return i
}

func commonSuffixLen(s1, s2 string) int {
	i1, i2 := len(s1)-1, len(s2)-1
	count := 0
	for i1 >= 0 && i2 >= 0 && s1[i1] == s2[i2] {
		count++
		i1--
		i2--
	}
	return count
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0

	for i := range a {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	cosine := dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
	return math.Max(0, math.Min(cosine, 1.0))
}
