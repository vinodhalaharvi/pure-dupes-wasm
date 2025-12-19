package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"path/filepath"
	"strings"
)

// Check if file is an image
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// DCT (Discrete Cosine Transform) for pHash
func dct2D(matrix [][]float64, size int) [][]float64 {
	result := make([][]float64, size)
	for i := range result {
		result[i] = make([]float64, size)
	}

	for u := 0; u < size; u++ {
		for v := 0; v < size; v++ {
			sum := 0.0
			for x := 0; x < size; x++ {
				for y := 0; y < size; y++ {
					sum += matrix[x][y] *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/(2*float64(size))) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/(2*float64(size)))
				}
			}

			cu := 1.0
			if u == 0 {
				cu = 1.0 / math.Sqrt(2)
			}
			cv := 1.0
			if v == 0 {
				cv = 1.0 / math.Sqrt(2)
			}

			result[u][v] = 0.25 * cu * cv * sum
		}
	}

	return result
}

// Simple nearest-neighbor resize
func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcW / width
			srcY := y * srcH / height
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

func calculateMedian(values []float64) float64 {
	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort (small array)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// Compute pHash for an image
func computePHash(imageData []byte) (uint64, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return 0, err
	}

	// Step 1: Resize to 32x32
	const size = 32
	resized := resizeImage(img, size, size)

	// Step 2: Convert to grayscale
	gray := make([][]float64, size)
	for i := range gray {
		gray[i] = make([]float64, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Convert to grayscale (0-255)
			gray[y][x] = float64(r>>8)*0.299 + float64(g>>8)*0.587 + float64(b>>8)*0.114
		}
	}

	// Step 3: Compute DCT
	dct := dct2D(gray, size)

	// Step 4: Keep only top-left 8x8 (low frequencies)
	const hashSize = 8
	lowFreq := make([]float64, 0, hashSize*hashSize)
	for i := 0; i < hashSize; i++ {
		for j := 0; j < hashSize; j++ {
			lowFreq = append(lowFreq, dct[i][j])
		}
	}

	// Step 5: Calculate median
	median := calculateMedian(lowFreq)

	// Step 6: Create binary hash
	var hash uint64
	for i, val := range lowFreq {
		if val > median {
			hash |= (1 << uint(i))
		}
	}

	return hash, nil
}

// Calculate Hamming distance between two hashes
func hammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	distance := 0
	for xor != 0 {
		distance += int(xor & 1)
		xor >>= 1
	}
	return distance
}

// Convert Hamming distance to similarity percentage
func hashSimilarity(hash1, hash2 uint64) float64 {
	distance := hammingDistance(hash1, hash2)
	return 1.0 - (float64(distance) / 64.0)
}

// Find visually similar images
func findVisualDuplicates(files []FileTree, threshold float64) map[string][]DuplicateMatch {
	reportProgress(0, 100, "Finding visually similar images...", 80)

	imageFiles := Filter(files, func(ft FileTree) bool {
		return ft.IsImage && ft.PHash != 0
	})

	matches := make(map[string][]DuplicateMatch)

	for i, src := range imageFiles {
		for j, tgt := range imageFiles {
			if i >= j {
				continue
			}

			similarity := hashSimilarity(src.PHash, tgt.PHash)

			if similarity >= threshold {
				matches[src.Path] = append(matches[src.Path], DuplicateMatch{
					TargetPath: tgt.Path,
					Similarity: similarity,
					SharedSize: src.Size,
					MatchType:  "visual",
				})

				matches[tgt.Path] = append(matches[tgt.Path], DuplicateMatch{
					TargetPath: src.Path,
					Similarity: similarity,
					SharedSize: tgt.Size,
					MatchType:  "visual",
				})
			}
		}
	}

	return matches
}
