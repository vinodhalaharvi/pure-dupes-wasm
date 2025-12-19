// main.go
package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed index.html
var content embed.FS

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

func (m Monoid[A]) FoldMap(xs []A, f func(A) A) A {
	return FoldLeft(Map(xs, f), m.Empty(), m.Combine)
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

var IntSumMonoid = Monoid[int]{
	Empty:   func() int { return 0 },
	Combine: func(a, b int) int { return a + b },
}

var Int64SumMonoid = Monoid[int64]{
	Empty:   func() int64 { return 0 },
	Combine: func(a, b int64) int64 { return a + b },
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
// FUNCTOR INSTANCES (No interface - Go lacks higher-kinded types)
// Each type implements FMap directly with full type safety
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

func (m Maybe[A]) GetOrElse(defaultVal A) A {
	if m.present {
		return m.value
	}
	return defaultVal
}

func (m Maybe[A]) IsPresent() bool {
	return m.present
}

func (m Maybe[A]) Get() A {
	return m.value
}

// ============================================================================
// FREE APPLICATIVE
// ============================================================================

type FreeAp[F any, A any] interface {
	freeAp()
}

type Pure[F any, A any] struct {
	Value A
}

func (p Pure[F, A]) freeAp() {}

type Ap[F any, A any, B any] struct {
	Fn   FreeAp[F, func(A) B]
	Args FreeAp[F, A]
}

func (a Ap[F, A, B]) freeAp() {}

type Lift[F any, A any] struct {
	FA F
}

func (l Lift[F, A]) freeAp() {}

func LiftFreeAp[F any, A any](fa F) FreeAp[F, A] {
	return Lift[F, A]{FA: fa}
}

func PureFreeAp[F any, A any](a A) FreeAp[F, A] {
	return Pure[F, A]{Value: a}
}

// Natural transformation interpreter
type NatTrans[F any, G any] func(F) G

// ============================================================================
// THUNK AS FUNCTOR
// ============================================================================

type Thunk[A any] struct {
	Compute func() A
	Name    string
}

func (t Thunk[A]) FMap(f func(A) A) Thunk[A] {
	return Thunk[A]{
		Compute: func() A { return f(t.Compute()) },
		Name:    t.Name + ".fmap",
	}
}

func RunThunk[A any](t Thunk[A]) A {
	return t.Compute()
}

func ParallelRunThunks[A any](thunks []Thunk[A], numWorkers int) []A {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	results := make([]A, len(thunks))
	jobs := make(chan int, len(thunks))
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				results[idx] = RunThunk(thunks[idx])
			}
		}()
	}

	for i := range thunks {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return results
}

// ============================================================================
// UTILITY FUNCTIONS (FUNCTOR MAP)
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
	Depth      int
}

type DuplicateMatch struct {
	TargetPath string
	Similarity float64
	SharedSize int64
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
}

type DedupResult struct {
	RootTree        FileNode
	AllMatches      map[string][]DuplicateMatch
	DuplicateGroups []DuplicateGroup
	TotalFiles      int
	UniqueFiles     int
	FullDupCount    int
	PartialDupCount int
	SpaceSaved      int64
}

// ============================================================================
// MERKLE TREE (USING FOLDS AND MONOIDS)
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
// FILE PROCESSING (FREE APPLICATIVE STYLE)
// ============================================================================

type FileOp interface {
	fileOp()
}

type ReadFile struct {
	Path      string
	ChunkSize int
}

func (ReadFile) fileOp() {}

type ScanDir struct {
	Dir      string
	MaxDepth int
}

func (ScanDir) fileOp() {}

func calculateDepth(rootDir, filePath string) int {
	rel, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return 0
	}
	if rel == "." {
		return 0
	}
	return strings.Count(rel, string(filepath.Separator)) + 1
}

func ProcessFileThunk(path string, chunkSize int, rootDir string) Thunk[FileTree] {
	return Thunk[FileTree]{
		Compute: func() FileTree {
			file, err := os.Open(path)
			if err != nil {
				return FileTree{}
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				return FileTree{}
			}

			chunks := readFileChunks(file, chunkSize)
			hashes := Map(chunks, HashLeaf)
			tree := BuildMerkleTree(hashes, SHA256Monoid)
			root := tree.Hash

			leafBytes := collectLeaves(tree)
			leaves := Map(leafBytes, func(b []byte) string {
				return hex.EncodeToString(b)
			})

			depth := calculateDepth(rootDir, path)

			return FileTree{
				Path:       path,
				Root:       root,
				Tree:       tree,
				Size:       stat.Size(),
				ChunkCount: len(chunks),
				Leaves:     leaves,
				Depth:      depth,
			}
		},
		Name: fmt.Sprintf("ProcessFile(%s)", filepath.Base(path)),
	}
}

func readFileChunks(file *os.File, chunkSize int) [][]byte {
	type ChunkAcc struct {
		chunks [][]byte
		done   bool
	}

	buffer := make([]byte, chunkSize)
	var allChunks [][]byte

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			allChunks = append(allChunks, chunk)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	return allChunks
}

func ScanDirectoryThunk(dir string, maxDepth int) Thunk[[]string] {
	return Thunk[[]string]{
		Compute: func() []string {
			var files []string

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				depth := calculateDepth(dir, path)

				if info.IsDir() && depth > maxDepth {
					return filepath.SkipDir
				}

				if !info.IsDir() && depth <= maxDepth {
					files = append(files, path)
				}

				return nil
			})

			return files
		},
		Name: fmt.Sprintf("ScanDir(%s, maxDepth=%d)", filepath.Base(dir), maxDepth),
	}
}

// ============================================================================
// DEDUPLICATION (USING FOLDS AND MONOIDS)
// ============================================================================

func BuildChunkIndex(files []FileTree) map[string][]int {
	indexMonoid := MapMonoid[string, []int](SliceMonoid[int]())

	return FoldMap(files, indexMonoid, func(ft FileTree) map[string][]int {
		fileIdx := findFileIndex(files, ft)
		return FoldLeft(ft.Leaves, make(map[string][]int), func(acc map[string][]int, chunkHash string) map[string][]int {
			acc[chunkHash] = []int{fileIdx}
			return acc
		})
	})
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
// TREE BUILDING (FOLD-BASED)
// ============================================================================

func BuildFileTree(rootPath string, files []FileTree, matches map[string][]DuplicateMatch, maxDepth int) FileNode {
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
		addToTree(&root, parts, ft, matches, 1, maxDepth, rootPath)
	}

	return root
}

func addToTree(node *FileNode, parts []string, ft FileTree, matches map[string][]DuplicateMatch, depth int, maxDepth int, rootPath string) {
	if depth > maxDepth || len(parts) == 0 {
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
	addToTree(dirNode, parts[1:], ft, matches, depth+1, maxDepth, rootPath)
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
// MAIN DEDUPLICATION (COMPOSITIONAL)
// ============================================================================

func FindDuplicatesThunk(dir string, threshold float64, maxDepth int, numWorkers int, chunkSize int) Thunk[DedupResult] {
	return Thunk[DedupResult]{
		Compute: func() DedupResult {
			// Pipeline composition using thunks
			allFiles := RunThunk(ScanDirectoryThunk(dir, maxDepth))

			fileThunks := Map(allFiles, func(path string) Thunk[FileTree] {
				return ProcessFileThunk(path, chunkSize, dir)
			})

			fileTrees := Filter(ParallelRunThunks(fileThunks, numWorkers),
				func(ft FileTree) bool { return ft.Path != "" })

			// Group exact duplicates using fold
			filesByRoot := GroupBy(fileTrees, func(ft FileTree) string {
				return hex.EncodeToString(ft.Root)
			})

			// Process duplicates using monoids and folds
			exactDups := processExactDuplicates(filesByRoot)
			partialDups := processPartialDuplicates(fileTrees, exactDups.allMatches, threshold)

			// Combine results using monoids
			allMatches := MapMonoid[string, []DuplicateMatch](SliceMonoid[DuplicateMatch]()).Combine(
				exactDups.allMatches,
				partialDups.allMatches,
			)

			tree := BuildFileTree(dir, fileTrees, allMatches, maxDepth)

			totalFiles := len(fileTrees)
			duplicateFileCount := exactDups.fullDupCount + partialDups.partialDupCount
			uniqueCount := totalFiles - duplicateFileCount

			return DedupResult{
				RootTree:        tree,
				AllMatches:      allMatches,
				DuplicateGroups: exactDups.groups,
				TotalFiles:      totalFiles,
				UniqueFiles:     uniqueCount,
				FullDupCount:    exactDups.fullDupCount,
				PartialDupCount: partialDups.partialDupCount,
				SpaceSaved:      exactDups.spaceSaved,
			}
		},
		Name: fmt.Sprintf("FindDuplicates(dir:%s, thresh:%.0f%%, depth:%d, workers:%d, chunk:%d)",
			filepath.Base(dir), threshold*100, maxDepth, numWorkers, chunkSize),
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
// HTTP SERVER
// ============================================================================

func handleDedup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Dir        string  `json:"dir"`
		Threshold  float64 `json:"threshold"`
		MaxDepth   int     `json:"maxDepth"`
		NumWorkers int     `json:"numWorkers"`
		ChunkSize  int     `json:"chunkSize"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.NumWorkers <= 0 {
		req.NumWorkers = runtime.NumCPU()
	}

	if req.ChunkSize <= 0 {
		req.ChunkSize = 4096
	}

	dedupThunk := FindDuplicatesThunk(req.Dir, req.Threshold, req.MaxDepth, req.NumWorkers, req.ChunkSize)

	result := RunThunk(dedupThunk)
	json.NewEncoder(w).Encode(result)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Println("Please open your browser to:", url)
	}
}

func findAvailablePort(startPort int) (int, error) {
	for port := startPort; port < startPort+100; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found")
}

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	port, err := findAvailablePort(8080)
	if err != nil {
		fmt.Println("Error finding available port:", err)
		return
	}

	fsys, err := fs.Sub(content, ".")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc("/api/dedup", handleDedup)

	url := fmt.Sprintf("http://localhost:%d", port)

	fmt.Println("🔍 pure-dupes - Functional duplicate finder")
	fmt.Println("📊 Server:", url)
	fmt.Println("🌳 Merkle trees + Monoids + Folds + Type-safe Functors")
	fmt.Printf("💻 Default workers: %d (CPU cores)\n", runtime.NumCPU())
	fmt.Println()

	serverReady := make(chan bool)
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Println("Error starting server:", err)
			serverReady <- false
			return
		}
		serverReady <- true
		http.Serve(listener, nil)
	}()

	if <-serverReady {
		if waitForServer(url, 5*time.Second) {
			openBrowser(url)
		} else {
			fmt.Println("Server started but not responding. Please open your browser to:", url)
		}
	}

	select {}
}
