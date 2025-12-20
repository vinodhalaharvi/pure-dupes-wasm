// main_wasm.go - Enhanced with progress reporting
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"
)

// ============================================================================
// PROGRESS REPORTING
// ============================================================================

var progressCallback js.Value

func reportProgress(current, total int, message string, percent float64) {
	if !progressCallback.IsUndefined() {
		progressCallback.Invoke(map[string]interface{}{
			"current": current,
			"total":   total,
			"message": message,
			"percent": percent,
		})
	}
}

// ============================================================================
// MONOID
// ============================================================================

type Monoid[A any] struct {
	Empty   func() A
	Combine func(A, A) A
}

func (m Monoid[A]) Fold(xs []A) A {
	return FoldLeft(xs, m.Empty(), m.Combine)
}

var SHA256Monoid = Monoid[[]byte]{
	Empty: func() []byte { return []byte{} },
	Combine: func(a, b []byte) []byte {
		h := sha256.New()
		h.Write(a)
		h.Write(b)
		return h.Sum(nil)
	},
}

func MapMonoid[K comparable, V any](vm Monoid[V]) Monoid[map[K]V] {
	return Monoid[map[K]V]{
		Empty: func() map[K]V { return make(map[K]V) },
		Combine: func(a, b map[K]V) map[K]V {
			result := make(map[K]V)
			for k, v := range a {
				result[k] = v
			}
			for k, v := range b {
				if existing, ok := result[k]; ok {
					result[k] = vm.Combine(existing, v)
				} else {
					result[k] = v
				}
			}
			return result
		},
	}
}

func SliceMonoid[A any]() Monoid[[]A] {
	return Monoid[[]A]{
		Empty: func() []A { return []A{} },
		Combine: func(a, b []A) []A {
			result := make([]A, 0, len(a)+len(b))
			result = append(result, a...)
			result = append(result, b...)
			return result
		},
	}
}

// ============================================================================
// FOLD OPERATIONS
// ============================================================================

func FoldLeft[A, B any](xs []A, zero B, f func(B, A) B) B {
	acc := zero
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

func FoldRight[A, B any](xs []A, zero B, f func(A, B) B) B {
	acc := zero
	for i := len(xs) - 1; i >= 0; i-- {
		acc = f(xs[i], acc)
	}
	return acc
}

func FoldMap[A, B any](xs []A, m Monoid[B], f func(A) B) B {
	return m.Fold(Map(xs, f))
}

// ============================================================================
// FUNCTOR INSTANCES
// ============================================================================

type Maybe[A any] struct {
	value   A
	present bool
}

func Just[A any](v A) Maybe[A] {
	return Maybe[A]{value: v, present: true}
}

func Nothing[A any]() Maybe[A] {
	return Maybe[A]{present: false}
}

func (m Maybe[A]) FMap(f func(A) A) Maybe[A] {
	if !m.present {
		return Nothing[A]()
	}
	return Just(f(m.value))
}

func (m Maybe[A]) IsPresent() bool {
	return m.present
}

func (m Maybe[A]) Get() A {
	return m.value
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

func Map[A, B any](xs []A, f func(A) B) []B {
	return FoldRight(xs, []B{}, func(a A, acc []B) []B {
		return append([]B{f(a)}, acc...)
	})
}

func Filter[A any](xs []A, pred func(A) bool) []A {
	return FoldRight(xs, []A{}, func(a A, acc []A) []A {
		if pred(a) {
			return append([]A{a}, acc...)
		}
		return acc
	})
}

func GroupBy[A any, K comparable](xs []A, key func(A) K) map[K][]A {
	return FoldLeft(xs, make(map[K][]A), func(acc map[K][]A, x A) map[K][]A {
		k := key(x)
		acc[k] = append(acc[k], x)
		return acc
	})
}

// ============================================================================
// DOMAIN TYPES
// ============================================================================

type MerkleNode struct {
	Hash     []byte
	Children []MerkleNode
	IsLeaf   bool
}

type FileTree struct {
	Path       string
	Root       []byte
	Tree       MerkleNode
	Size       int64
	ChunkCount int
	Leaves     []string
	ModTime    int64
	PHash      uint64   // Phase 2: Image perceptual hash
	IsImage    bool     // Phase 2: Is this an image file?
	VideoHash  []uint64 // Phase 2: Video frame hashes (array of pHashes)
	IsVideo    bool     // Phase 2: Is this a video file?
}

type DuplicateMatch struct {
	TargetPath string
	Similarity float64
	SharedSize int64
	MatchType  string // "exact", "partial", "content"
}

type FileNode struct {
	Path         string
	Name         string
	IsDir        bool
	Children     []FileNode
	Matches      []DuplicateMatch
	BestMatch    float64
	Size         int64
	RelativePath string
}

type DuplicateGroup struct {
	Files      []string
	Similarity float64
	Size       int64
	GroupType  string // "exact", "similar"
	Savings    int64
}

type DedupResult struct {
	RootTree        FileNode
	AllMatches      map[string][]DuplicateMatch
	DuplicateGroups []DuplicateGroup
	TotalFiles      int
	UniqueFiles     int
	FullDupCount    int
	PartialDupCount int
	VisualDupCount  int // Phase 2: Visual duplicate count
	SpaceSaved      int64
	ProcessingTime  float64
}

type JSFile struct {
	Name             string
	Path             string
	Size             int64
	Data             []byte
	ModTime          int64
	VideoFrameHashes []uint64 // Phase 2: Video frame hashes from JavaScript
}

// ============================================================================
// MERKLE TREE
// ============================================================================

func HashLeaf(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func BuildMerkleTree(hashes [][]byte, m Monoid[[]byte]) MerkleNode {
	if len(hashes) == 0 {
		return MerkleNode{Hash: m.Empty(), IsLeaf: true, Children: []MerkleNode{}}
	}
	if len(hashes) == 1 {
		return MerkleNode{Hash: hashes[0], IsLeaf: true, Children: []MerkleNode{}}
	}

	pairs := pairwiseFold(hashes, m)

	if len(pairs) == 1 {
		return pairs[0]
	}

	parentHashes := Map(pairs, func(n MerkleNode) []byte { return n.Hash })
	upperTree := BuildMerkleTree(parentHashes, m)
	upperTree.Children = pairs
	return upperTree
}

func pairwiseFold(hashes [][]byte, m Monoid[[]byte]) []MerkleNode {
	type Acc struct {
		nodes   []MerkleNode
		pending Maybe[[]byte]
	}

	result := FoldLeft(hashes, Acc{nodes: []MerkleNode{}, pending: Nothing[[]byte]()},
		func(acc Acc, hash []byte) Acc {
			if !acc.pending.IsPresent() {
				return Acc{nodes: acc.nodes, pending: Just(hash)}
			}

			left := MerkleNode{Hash: acc.pending.Get(), IsLeaf: true, Children: []MerkleNode{}}
			right := MerkleNode{Hash: hash, IsLeaf: true, Children: []MerkleNode{}}
			combined := m.Combine(acc.pending.Get(), hash)

			parent := MerkleNode{
				Hash:     combined,
				IsLeaf:   false,
				Children: []MerkleNode{left, right},
			}

			return Acc{nodes: append(acc.nodes, parent), pending: Nothing[[]byte]()}
		})

	if result.pending.IsPresent() {
		single := MerkleNode{Hash: result.pending.Get(), IsLeaf: true, Children: []MerkleNode{}}
		return append(result.nodes, single)
	}

	return result.nodes
}

func collectLeaves(node MerkleNode) [][]byte {
	if node.IsLeaf {
		return [][]byte{node.Hash}
	}

	leafMonoid := SliceMonoid[[]byte]()
	return FoldMap(node.Children, leafMonoid, collectLeaves)
}

// ============================================================================
// FILE PROCESSING
// ============================================================================

func ProcessFile(file JSFile, chunkSize int, index int, total int) FileTree {
	reportProgress(index, total, fmt.Sprintf("Processing %s", file.Name), float64(index)/float64(total)*100)

	data := file.Data
	chunks := [][]byte{}

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}

	hashes := Map(chunks, HashLeaf)
	tree := BuildMerkleTree(hashes, SHA256Monoid)
	root := tree.Hash

	leafBytes := collectLeaves(tree)
	leaves := Map(leafBytes, func(b []byte) string {
		return hex.EncodeToString(b)
	})

	// Phase 2: Compute pHash for images
	var pHash uint64
	isImage := isImageFile(file.Path)
	if isImage {
		hash, err := computePHash(data)
		if err == nil {
			pHash = hash
		}
	}

	// Phase 2: Video frame hashes (computed by JavaScript)
	var videoHash []uint64
	isVideo := isVideoFile(file.Path)
	if isVideo && len(file.VideoFrameHashes) > 0 {
		videoHash = file.VideoFrameHashes
	}

	return FileTree{
		Path:       file.Path,
		Root:       root,
		Tree:       tree,
		Size:       file.Size,
		ChunkCount: len(chunks),
		Leaves:     leaves,
		ModTime:    file.ModTime,
		PHash:      pHash,
		IsImage:    isImage,
		VideoHash:  videoHash,
		IsVideo:    isVideo,
	}
}

// ============================================================================
// DEDUPLICATION
// ============================================================================

func BuildChunkIndex(files []FileTree) map[string][]int {
	reportProgress(0, 100, "Building chunk index...", 0)

	indexMonoid := MapMonoid[string, []int](SliceMonoid[int]())

	result := FoldMap(files, indexMonoid, func(ft FileTree) map[string][]int {
		fileIdx := findFileIndex(files, ft)
		return FoldLeft(ft.Leaves, make(map[string][]int), func(acc map[string][]int, chunkHash string) map[string][]int {
			acc[chunkHash] = []int{fileIdx}
			return acc
		})
	})

	reportProgress(100, 100, "Chunk index built", 100)
	return result
}

func findFileIndex(files []FileTree, target FileTree) int {
	for i, f := range files {
		if f.Path == target.Path {
			return i
		}
	}
	return -1
}

func FindCandidates(sourceFile FileTree, chunkIndex map[string][]int, threshold float64) map[int]int {
	countMap := FoldLeft(sourceFile.Leaves, make(map[int]int),
		func(acc map[int]int, chunkHash string) map[int]int {
			if targets, exists := chunkIndex[chunkHash]; exists {
				for _, targetIdx := range targets {
					acc[targetIdx]++
				}
			}
			return acc
		})

	minSharedChunks := int(float64(len(sourceFile.Leaves)) * threshold)

	return FoldLeft(mapToSlice(countMap), make(map[int]int),
		func(acc map[int]int, pair struct {
			k int
			v int
		}) map[int]int {
			if pair.v >= minSharedChunks {
				acc[pair.k] = pair.v
			}
			return acc
		})
}

func mapToSlice[K comparable, V any](m map[K]V) []struct {
	k K
	v V
} {
	result := make([]struct {
		k K
		v V
	}, 0, len(m))
	for k, v := range m {
		result = append(result, struct {
			k K
			v V
		}{k, v})
	}
	return result
}

func CompareFiles(a, b FileTree) float64 {
	if hex.EncodeToString(a.Root) == hex.EncodeToString(b.Root) {
		return 1.0
	}

	if len(a.Leaves) == 0 || len(b.Leaves) == 0 {
		return 0.0
	}

	setB := FoldLeft(b.Leaves, make(map[string]bool, len(b.Leaves)),
		func(acc map[string]bool, leaf string) map[string]bool {
			acc[leaf] = true
			return acc
		})

	matches := FoldLeft(a.Leaves, 0, func(acc int, leaf string) int {
		if setB[leaf] {
			return acc + 1
		}
		return acc
	})

	return float64(matches) / float64(len(a.Leaves))
}

// ============================================================================
// SMART DUPLICATE GROUPS
// ============================================================================

func CreateSmartGroups(filesByRoot map[string][]FileTree, partialMatches map[string][]DuplicateMatch, visualMatches map[string][]DuplicateMatch, fileTrees []FileTree) []DuplicateGroup {
	groups := []DuplicateGroup{}

	// Exact duplicate groups
	for _, group := range filesByRoot {
		if len(group) <= 1 {
			continue
		}

		groupFiles := Map(group, func(ft FileTree) string { return ft.Path })

		// Calculate savings (keep largest, remove rest)
		totalSize := FoldLeft(group, int64(0), func(acc int64, ft FileTree) int64 {
			return acc + ft.Size
		})
		savings := totalSize - group[0].Size

		groups = append(groups, DuplicateGroup{
			Files:      groupFiles,
			Similarity: 1.0,
			Size:       totalSize,
			GroupType:  "exact",
			Savings:    savings,
		})
	}

	// Partial duplicate groups
	processed := make(map[string]bool)
	for srcPath, matches := range partialMatches {
		if processed[srcPath] {
			continue
		}

		// Create group with source and all its matches
		groupFiles := []string{srcPath}
		totalSize := int64(0)

		for _, match := range matches {
			if match.Similarity >= 0.8 && !processed[match.TargetPath] {
				groupFiles = append(groupFiles, match.TargetPath)
			}
		}

		if len(groupFiles) > 1 {
			for _, path := range groupFiles {
				processed[path] = true
			}

			groups = append(groups, DuplicateGroup{
				Files:      groupFiles,
				Similarity: 0.8,
				Size:       totalSize,
				GroupType:  "similar",
				Savings:    totalSize / 2, // Estimate
			})
		}
	}

	// Phase 2: Visual duplicate groups
	processedVisual := make(map[string]bool)
	for srcPath, matches := range visualMatches {
		if processedVisual[srcPath] {
			continue
		}

		groupFiles := []string{srcPath}
		var totalSimilarity float64
		var matchCount int
		var totalSize int64

		for _, match := range matches {
			if match.Similarity >= 0.85 && !processedVisual[match.TargetPath] {
				groupFiles = append(groupFiles, match.TargetPath)
				totalSimilarity += match.Similarity
				matchCount++
			}
		}

		if len(groupFiles) > 1 {
			for _, path := range groupFiles {
				processedVisual[path] = true
			}

			// Calculate average similarity
			avgSimilarity := 0.9
			if matchCount > 0 {
				avgSimilarity = totalSimilarity / float64(matchCount)
			}

			// Calculate total size (find files in fileTrees)
			for _, ft := range fileTrees {
				for _, gf := range groupFiles {
					if ft.Path == gf {
						totalSize += ft.Size
						break
					}
				}
			}

			// Estimate savings (keep one, remove rest)
			savings := totalSize
			if len(groupFiles) > 0 {
				// Assume we keep the first one
				for _, ft := range fileTrees {
					if ft.Path == groupFiles[0] {
						savings = totalSize - ft.Size
						break
					}
				}
			}

			groups = append(groups, DuplicateGroup{
				Files:      groupFiles,
				Similarity: avgSimilarity,
				Size:       totalSize,
				GroupType:  "visual",
				Savings:    savings,
			})
		}
	}

	return groups
}

// ============================================================================
// TREE BUILDING
// ============================================================================

func BuildFileTree(rootPath string, files []FileTree, matches map[string][]DuplicateMatch) FileNode {
	relFiles := FoldLeft(files, make(map[string]FileTree),
		func(acc map[string]FileTree, f FileTree) map[string]FileTree {
			rel, _ := filepath.Rel(rootPath, f.Path)
			acc[rel] = f
			return acc
		})

	root := FileNode{
		Path:         rootPath,
		Name:         filepath.Base(rootPath),
		IsDir:        true,
		Children:     []FileNode{},
		RelativePath: "",
	}

	for rel, ft := range relFiles {
		parts := strings.Split(rel, string(filepath.Separator))
		addToTree(&root, parts, ft, matches, rootPath)
	}

	return root
}

func addToTree(node *FileNode, parts []string, ft FileTree, matches map[string][]DuplicateMatch, rootPath string) {
	if len(parts) == 0 {
		return
	}

	if len(parts) == 1 {
		fileMatches := matches[ft.Path]
		bestMatch := FoldLeft(fileMatches, 0.0, func(acc float64, m DuplicateMatch) float64 {
			if m.Similarity > acc {
				return m.Similarity
			}
			return acc
		})

		node.Children = append(node.Children, FileNode{
			Path:         ft.Path,
			Name:         parts[0],
			IsDir:        false,
			Children:     []FileNode{},
			Matches:      fileMatches,
			BestMatch:    bestMatch,
			Size:         ft.Size,
			RelativePath: ft.Path,
		})
		return
	}

	dirName := parts[0]
	dirNode := findOrCreateDir(node, dirName, rootPath)
	addToTree(dirNode, parts[1:], ft, matches, rootPath)
}

func findOrCreateDir(node *FileNode, dirName string, rootPath string) *FileNode {
	for i := range node.Children {
		if node.Children[i].Name == dirName && node.Children[i].IsDir {
			return &node.Children[i]
		}
	}

	currentPath := filepath.Join(node.Path, dirName)
	relativePath, _ := filepath.Rel(rootPath, currentPath)

	node.Children = append(node.Children, FileNode{
		Name:         dirName,
		Path:         currentPath,
		RelativePath: relativePath,
		IsDir:        true,
		Children:     []FileNode{},
	})

	return &node.Children[len(node.Children)-1]
}

// ============================================================================
// MAIN DEDUPLICATION
// ============================================================================

func FindDuplicates(files []JSFile, threshold float64, chunkSize int) DedupResult {
	startTime := time.Now()

	reportProgress(0, 100, "Starting analysis...", 0)

	// Process all files with progress
	fileTrees := make([]FileTree, len(files))
	for i, f := range files {
		fileTrees[i] = ProcessFile(f, chunkSize, i, len(files))
	}

	reportProgress(30, 100, "Grouping files...", 30)

	// Group by merkle root
	filesByRoot := GroupBy(fileTrees, func(ft FileTree) string {
		return hex.EncodeToString(ft.Root)
	})

	reportProgress(50, 100, "Finding exact duplicates...", 50)

	// Process duplicates
	exactDups := processExactDuplicates(filesByRoot)

	reportProgress(70, 100, "Finding similar files...", 70)

	partialDups := processPartialDuplicates(fileTrees, exactDups.allMatches, threshold)

	reportProgress(80, 100, "Finding visually similar images...", 80)

	// Phase 2: Visual duplicates (optimized - skip exact matches)
	// Build set of files already in exact duplicate groups
	filesInExactGroups := make(map[string]bool)
	for _, group := range filesByRoot {
		if len(group) > 1 {
			for _, ft := range group {
				filesInExactGroups[ft.Path] = true
			}
		}
	}

	// Only check images NOT in exact duplicate groups
	imagesToCheck := Filter(fileTrees, func(ft FileTree) bool {
		return ft.IsImage && ft.PHash != 0 && !filesInExactGroups[ft.Path]
	})

	// Only check videos NOT in exact duplicate groups
	videosToCheck := Filter(fileTrees, func(ft FileTree) bool {
		return ft.IsVideo && len(ft.VideoHash) > 0 && !filesInExactGroups[ft.Path]
	})

	reportProgress(81, 100, fmt.Sprintf("Checking %d images and %d videos for visual similarity...", len(imagesToCheck), len(videosToCheck)), 81)

	// Combine images and videos for visual duplicate detection
	mediaToCheck := append(imagesToCheck, videosToCheck...)
	visualDups := findVisualDuplicates(mediaToCheck, 0.85)
	visualCount := len(visualDups)

	reportProgress(85, 100, "Creating smart groups...", 85)

	// Smart groups (now includes visual matches)
	smartGroups := CreateSmartGroups(filesByRoot, partialDups.allMatches, visualDups, fileTrees)

	reportProgress(90, 100, "Building file tree...", 90)

	// Combine results (Phase 1 + Phase 2)
	allMatches := MapMonoid[string, []DuplicateMatch](SliceMonoid[DuplicateMatch]()).Combine(
		exactDups.allMatches,
		MapMonoid[string, []DuplicateMatch](SliceMonoid[DuplicateMatch]()).Combine(
			partialDups.allMatches,
			visualDups,
		),
	)

	rootPath := "/"
	if len(files) > 0 {
		rootPath = filepath.Dir(files[0].Path)
	}

	tree := BuildFileTree(rootPath, fileTrees, allMatches)

	totalFiles := len(fileTrees)
	duplicateFileCount := exactDups.fullDupCount + partialDups.partialDupCount
	uniqueCount := totalFiles - duplicateFileCount

	processingTime := time.Since(startTime).Seconds()

	reportProgress(100, 100, "Analysis complete!", 100)

	return DedupResult{
		RootTree:        tree,
		AllMatches:      allMatches,
		DuplicateGroups: smartGroups,
		TotalFiles:      totalFiles,
		UniqueFiles:     uniqueCount,
		FullDupCount:    exactDups.fullDupCount,
		PartialDupCount: partialDups.partialDupCount,
		VisualDupCount:  visualCount,
		SpaceSaved:      exactDups.spaceSaved,
		ProcessingTime:  processingTime,
	}
}

type ExactDupsResult struct {
	allMatches   map[string][]DuplicateMatch
	groups       []DuplicateGroup
	fullDupCount int
	spaceSaved   int64
}

func processExactDuplicates(filesByRoot map[string][]FileTree) ExactDupsResult {
	duplicateGroups := Filter(mapToSlice(filesByRoot),
		func(pair struct {
			k string
			v []FileTree
		}) bool {
			return len(pair.v) > 1
		})

	type DupAcc struct {
		matches        map[string][]DuplicateMatch
		groups         []DuplicateGroup
		count          int
		saved          int64
		processedPaths map[string]bool
	}

	result := FoldLeft(duplicateGroups, DupAcc{
		matches:        make(map[string][]DuplicateMatch),
		groups:         []DuplicateGroup{},
		count:          0,
		saved:          0,
		processedPaths: make(map[string]bool),
	}, func(acc DupAcc, pair struct {
		k string
		v []FileTree
	}) DupAcc {
		group := pair.v
		groupFiles := Map(group, func(ft FileTree) string { return ft.Path })

		acc.groups = append(acc.groups, DuplicateGroup{
			Files:      groupFiles,
			Similarity: 1.0,
			Size:       group[0].Size,
		})

		for _, src := range group {
			matches := FoldLeft(group, []DuplicateMatch{},
				func(macc []DuplicateMatch, tgt FileTree) []DuplicateMatch {
					if src.Path != tgt.Path {
						return append(macc, DuplicateMatch{
							TargetPath: tgt.Path,
							Similarity: 1.0,
							SharedSize: src.Size,
							MatchType:  "exact",
						})
					}
					return macc
				})

			if len(matches) > 0 {
				acc.matches[src.Path] = matches
				acc.count++
				if !acc.processedPaths[src.Path] {
					acc.saved += src.Size
					acc.processedPaths[src.Path] = true
				}
			}
		}

		return acc
	})

	return ExactDupsResult{
		allMatches:   result.matches,
		groups:       result.groups,
		fullDupCount: result.count,
		spaceSaved:   result.saved,
	}
}

type PartialDupsResult struct {
	allMatches      map[string][]DuplicateMatch
	partialDupCount int
}

func processPartialDuplicates(fileTrees []FileTree, exactMatches map[string][]DuplicateMatch, threshold float64) PartialDupsResult {
	chunkIndex := BuildChunkIndex(fileTrees)

	type FileWithIndex struct {
		file  FileTree
		index int
	}

	filesWithIndices := Map(fileTrees, func(ft FileTree) FileWithIndex {
		idx := findFileIndex(fileTrees, ft)
		return FileWithIndex{file: ft, index: idx}
	})

	candidateFiles := Filter(filesWithIndices, func(fwi FileWithIndex) bool {
		_, hasExact := exactMatches[fwi.file.Path]
		return !hasExact
	})

	type PartialAcc struct {
		matches map[string][]DuplicateMatch
		count   int
	}

	result := FoldLeft(candidateFiles, PartialAcc{
		matches: make(map[string][]DuplicateMatch),
		count:   0,
	}, func(acc PartialAcc, fwi FileWithIndex) PartialAcc {
		src := fwi.file
		srcIdx := fwi.index

		candidates := FindCandidates(src, chunkIndex, threshold)

		matches := FoldLeft(mapToSlice(candidates), []DuplicateMatch{},
			func(macc []DuplicateMatch, pair struct {
				k int
				v int
			}) []DuplicateMatch {
				targetIdx := pair.k
				if targetIdx == srcIdx {
					return macc
				}

				tgt := fileTrees[targetIdx]
				srcRoot := hex.EncodeToString(src.Root)
				tgtRoot := hex.EncodeToString(tgt.Root)

				if srcRoot == tgtRoot {
					return macc
				}

				similarity := CompareFiles(src, tgt)

				if similarity >= threshold && similarity < 1.0 {
					return append(macc, DuplicateMatch{
						TargetPath: tgt.Path,
						Similarity: similarity,
						SharedSize: int64(float64(src.Size) * similarity),
						MatchType:  "partial",
					})
				}

				return macc
			})

		if len(matches) > 0 {
			acc.matches[src.Path] = matches
			acc.count++
		}

		return acc
	})

	return PartialDupsResult{
		allMatches:      result.matches,
		partialDupCount: result.count,
	}
}

// ============================================================================
// WASM EXPORTS
// ============================================================================

func analyzeFiles(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return map[string]interface{}{
			"error": "Expected 3 arguments: files, threshold, chunkSize",
		}
	}

	filesJS := args[0]
	threshold := args[1].Float()
	chunkSize := args[2].Int()

	// Set progress callback if provided
	if len(args) >= 4 && !args[3].IsUndefined() {
		progressCallback = args[3]
	}

	// Convert JS files to Go structs
	length := filesJS.Length()
	files := make([]JSFile, length)

	for i := 0; i < length; i++ {
		fileJS := filesJS.Index(i)

		dataJS := fileJS.Get("data")
		dataLen := dataJS.Get("length").Int()
		data := make([]byte, dataLen)
		js.CopyBytesToGo(data, dataJS)

		modTime := int64(0)
		if !fileJS.Get("modTime").IsUndefined() {
			modTime = int64(fileJS.Get("modTime").Int())
		}

		files[i] = JSFile{
			Name:    fileJS.Get("name").String(),
			Path:    fileJS.Get("path").String(),
			Size:    int64(fileJS.Get("size").Int()),
			Data:    data,
			ModTime: modTime,
		}
	}

	// Run deduplication
	result := FindDuplicates(files, threshold, chunkSize)

	// Convert result to JSON
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to marshal result: %v", err),
		}
	}

	return string(jsonBytes)
}

func main() {
	c := make(chan struct{})

	js.Global().Set("analyzeFiles", js.FuncOf(analyzeFiles))

	fmt.Println("🔍 pure-dupes WASM initialized")
	fmt.Println("✨ Phase 1 Features: Web Workers, Caching, Smart Groups, Progress")

	<-c
}
