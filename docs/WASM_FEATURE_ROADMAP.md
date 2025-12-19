# WASM Feature Roadmap & Suggestions

## 🎯 Current State

You have:
- ✅ Content-based deduplication (Merkle trees)
- ✅ Exact & partial matching
- ✅ Client-side processing
- ✅ Functional architecture

## 🚀 Immediate Wins (Low Hanging Fruit)

### 1. Web Workers - Parallel Processing
**Why:** Keep UI responsive, utilize multiple CPU cores
**Impact:** 🔥🔥🔥 High
**Effort:** 🛠️ Medium

```javascript
// wasm-worker.js
importScripts('wasm_exec.js');

const go = new Go();
WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
  .then(result => {
    go.run(result.instance);
    
    self.onmessage = (e) => {
      const {files, threshold, chunkSize} = e.data;
      const result = analyzeFiles(files, threshold, chunkSize);
      self.postMessage(result);
    };
  });
```

**Benefits:**
- Non-blocking UI
- Progress updates
- Cancellable operations
- True parallel processing

---

### 2. IndexedDB Caching
**Why:** Remember analyzed files, instant re-scans
**Impact:** 🔥🔥🔥 High
**Effort:** 🛠️ Low-Medium

```javascript
// Cache file hashes
const db = await openDB('pure-dupes', 1, {
  upgrade(db) {
    db.createObjectStore('file-hashes', {keyPath: 'path'});
    db.createObjectStore('analysis-results');
  }
});

// Store hash
await db.put('file-hashes', {
  path: file.path,
  hash: merkleRoot,
  size: file.size,
  timestamp: Date.now()
});

// Check cache before processing
const cached = await db.get('file-hashes', file.path);
if (cached && cached.size === file.size) {
  // Skip processing, use cached hash
}
```

**Benefits:**
- 10-100x faster re-analysis
- Persistent across sessions
- Incremental updates (only process new files)

---

### 3. Progress Reporting
**Why:** User feedback for long operations
**Impact:** 🔥🔥 Medium-High
**Effort:** 🛠️ Low

**Go side:**
```go
// Add to WASM exports
func reportProgress(current, total int, message string) {
    js.Global().Get("onProgress").Invoke(current, total, message)
}

// Use in processing
for i, file := range files {
    reportProgress(i, len(files), fmt.Sprintf("Processing %s", file.Name))
    // ... process file
}
```

**JS side:**
```javascript
window.onProgress = (current, total, message) => {
  setProgress({current, total, message});
};
```

---

## 🎨 High-Value Features

### 4. Image Similarity Detection ⭐
**Why:** Find duplicate photos with edits, crops, filters
**Impact:** 🔥🔥🔥🔥 Very High
**Effort:** 🛠️🛠️ Medium-High

**Approach: Perceptual Hashing (pHash)**

```go
// Add to Go WASM
func computeImagePHash(imageData []byte) string {
    img, _ := decodeImage(imageData)
    
    // Resize to 32x32
    resized := resize(img, 32, 32)
    
    // Convert to grayscale
    gray := toGrayscale(resized)
    
    // Apply DCT (Discrete Cosine Transform)
    dct := applyDCT(gray)
    
    // Keep low frequencies (8x8)
    lowFreq := dct[:8][:8]
    
    // Compute median
    median := computeMedian(lowFreq)
    
    // Generate hash (0/1 based on median)
    hash := ""
    for _, val := range lowFreq {
        if val > median {
            hash += "1"
        } else {
            hash += "0"
        }
    }
    
    return hash
}

// Hamming distance for similarity
func hammingDistance(hash1, hash2 string) int {
    distance := 0
    for i := range hash1 {
        if hash1[i] != hash2[i] {
            distance++
        }
    }
    return distance
}
```

**Use cases:**
- Find duplicate photos
- Detect edited versions (cropped, filtered, resized)
- Group similar images
- Threshold: < 10 hamming distance = very similar

**Libraries to consider:**
- Use Go's `image` package
- Or call browser's Canvas API for image processing

---

### 5. Audio Fingerprinting
**Why:** Find duplicate music files (different formats, bitrates)
**Impact:** 🔥🔥🔥 High (for music collections)
**Effort:** 🛠️🛠️🛠️ High

**Approach: Chromaprint/AcoustID**

```go
// Generate audio fingerprint
func computeAudioFingerprint(audioData []byte) string {
    // Decode audio (MP3, M4A, FLAC, etc.)
    samples := decodeAudio(audioData)
    
    // Apply FFT (Fast Fourier Transform)
    spectrum := applyFFT(samples)
    
    // Extract features
    features := extractChromaFeatures(spectrum)
    
    // Generate fingerprint
    fingerprint := generateFingerprint(features)
    
    return fingerprint
}
```

**Libraries:**
- Chromaprint algorithm (Rust WASM port available)
- Or use Web Audio API + WASM for processing

**Use cases:**
- Find duplicate songs in different formats
- Detect covers/remixes
- Match songs across quality levels

---

### 6. Text Similarity (Fuzzy Matching)
**Why:** Find duplicate documents with minor edits
**Impact:** 🔥🔥 Medium-High
**Effort:** 🛠️ Low-Medium

**Approach: MinHash + Simhash**

```go
// Text fingerprinting
func computeTextFingerprint(text string) uint64 {
    // Tokenize
    tokens := tokenize(text)
    
    // Generate shingles (n-grams)
    shingles := generateShingles(tokens, 3)
    
    // Compute simhash
    hash := simhash(shingles)
    
    return hash
}

// For documents
func analyzeDocument(docData []byte, fileType string) TextFingerprint {
    var text string
    
    switch fileType {
    case "txt":
        text = string(docData)
    case "pdf":
        text = extractTextFromPDF(docData)
    case "docx":
        text = extractTextFromDocx(docData)
    }
    
    return TextFingerprint{
        Hash: computeTextFingerprint(text),
        WordCount: countWords(text),
    }
}
```

**Use cases:**
- Duplicate documents
- Plagiarism detection
- Version tracking

---

## 🔮 Advanced Features

### 7. Machine Learning Similarity
**Why:** Semantic similarity, not just byte-level
**Impact:** 🔥🔥🔥🔥 Very High
**Effort:** 🛠️🛠️🛠️🛠️ Very High

**Approach: TensorFlow.js + WASM**

```javascript
// Load pre-trained model
const model = await tf.loadGraphModel('mobilenet/model.json');

// Generate embeddings
async function generateEmbedding(imageData) {
  const tensor = tf.browser.fromPixels(imageData);
  const resized = tf.image.resizeBilinear(tensor, [224, 224]);
  const normalized = resized.div(255.0);
  
  const embedding = model.predict(normalized.expandDims());
  return embedding.arraySync();
}

// Cosine similarity
function cosineSimilarity(a, b) {
  const dotProduct = a.reduce((sum, val, i) => sum + val * b[i], 0);
  const magA = Math.sqrt(a.reduce((sum, val) => sum + val * val, 0));
  const magB = Math.sqrt(b.reduce((sum, val) => sum + val * val, 0));
  return dotProduct / (magA * magB);
}
```

**Use cases:**
- Semantic image similarity
- Find conceptually similar images
- Content-based recommendations

---

### 8. Video Deduplication
**Why:** Find duplicate videos (re-encodes, edits)
**Impact:** 🔥🔥🔥 High
**Effort:** 🛠️🛠️🛠️🛠️ Very High

**Approach: Frame sampling + perceptual hashing**

```go
func analyzeVideo(videoData []byte) VideoFingerprint {
    // Sample frames (every N seconds)
    frames := extractFrames(videoData, samplingRate)
    
    // Compute pHash for each frame
    hashes := make([]string, len(frames))
    for i, frame := range frames {
        hashes[i] = computeImagePHash(frame)
    }
    
    // Generate video fingerprint
    return VideoFingerprint{
        FrameHashes: hashes,
        Duration: getDuration(videoData),
        Resolution: getResolution(videoData),
    }
}
```

**Challenges:**
- Large file sizes
- Processing time
- Memory constraints

---

## 🎁 UI/UX Enhancements

### 9. File Preview Generation
**Why:** Visual confirmation before deletion
**Impact:** 🔥🔥🔥 High
**Effort:** 🛠️ Medium

```javascript
// Generate thumbnails
async function generateThumbnail(file) {
  if (file.type.startsWith('image/')) {
    const img = new Image();
    img.src = URL.createObjectURL(file);
    await img.decode();
    
    const canvas = document.createElement('canvas');
    canvas.width = 200;
    canvas.height = 200;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0, 200, 200);
    
    return canvas.toDataURL();
  }
  
  // PDF thumbnails using PDF.js
  // Video thumbnails using video element
}
```

---

### 10. Batch Operations
**Why:** Delete/move multiple duplicates at once
**Impact:** 🔥🔥 Medium
**Effort:** 🛠️ Low

```javascript
// Select duplicates
const [selected, setSelected] = useState(new Set());

// Bulk actions
async function deleteSelected() {
  for (const path of selected) {
    // Use File System Access API
    await deleteFile(path);
  }
}
```

**Note:** File System Access API has limited delete support - may need to download list for manual deletion

---

## 📊 Suggested Priority Order

### Phase 1: Foundation (1-2 weeks)
1. ✅ Web Workers (parallel processing)
2. ✅ IndexedDB caching
3. ✅ Progress reporting

**Impact:** Makes the tool production-ready

### Phase 2: Content Intelligence (2-3 weeks)
4. ✅ Image similarity (pHash)
5. ✅ Text similarity (simhash)
6. ✅ File preview generation

**Impact:** 10x more useful, catches duplicates content-based matching misses

### Phase 3: Advanced (1-2 months)
7. ✅ Audio fingerprinting
8. ✅ ML-based similarity
9. ✅ Video deduplication

**Impact:** Professional-grade tool, unique in the market

### Phase 4: Polish (1 week)
10. ✅ Batch operations
11. ✅ Export results
12. ✅ Statistics dashboard

---

## 🛠️ Implementation Strategy

### Step 1: Add Type Detection

```go
type FileType int

const (
    FileTypeUnknown FileType = iota
    FileTypeImage
    FileTypeAudio
    FileTypeVideo
    FileTypeText
    FileTypeDocument
)

func detectFileType(data []byte, name string) FileType {
    // Check magic bytes
    if isJPEG(data) || isPNG(data) || isGIF(data) {
        return FileTypeImage
    }
    if isMP3(data) || isM4A(data) || isFLAC(data) {
        return FileTypeAudio
    }
    // ... etc
}
```

### Step 2: Extend Analysis Pipeline

```go
type FileAnalysis struct {
    Path            string
    Size            int64
    Type            FileType
    ContentHash     []byte          // Current Merkle root
    PerceptualHash  string          // For images
    AudioPrint      string          // For audio
    TextHash        uint64          // For documents
    Embedding       []float32       // For ML
}

func AnalyzeFile(file JSFile) FileAnalysis {
    analysis := FileAnalysis{
        Path: file.Path,
        Size: file.Size,
        Type: detectFileType(file.Data, file.Name),
    }
    
    // Content hash (always)
    analysis.ContentHash = computeMerkleRoot(file.Data)
    
    // Type-specific analysis
    switch analysis.Type {
    case FileTypeImage:
        analysis.PerceptualHash = computeImagePHash(file.Data)
    case FileTypeAudio:
        analysis.AudioPrint = computeAudioFingerprint(file.Data)
    case FileTypeText:
        analysis.TextHash = computeTextFingerprint(string(file.Data))
    }
    
    return analysis
}
```

### Step 3: Multi-Method Matching

```go
func FindDuplicates(analyses []FileAnalysis, threshold float64) DedupResult {
    matches := make(map[string][]DuplicateMatch)
    
    // Method 1: Exact content match (current)
    exactMatches := findExactMatches(analyses)
    
    // Method 2: Perceptual matching
    imageMatches := findSimilarImages(analyses, threshold)
    
    // Method 3: Audio matching
    audioMatches := findSimilarAudio(analyses, threshold)
    
    // Method 4: Text matching
    textMatches := findSimilarText(analyses, threshold)
    
    // Merge results
    return mergeMatches(exactMatches, imageMatches, audioMatches, textMatches)
}
```

---

## 📦 Library Recommendations

### Image Processing
- **Go:** `golang.org/x/image` (built-in)
- **WASM:** Use Canvas API for thumbnails
- **Algorithm:** pHash (perceptual hash)

### Audio Processing
- **Go:** `github.com/mjibson/go-dsp/fft`
- **Alternative:** Web Audio API
- **Algorithm:** Chromaprint

### Text Processing
- **Go:** Standard library sufficient
- **Algorithm:** Simhash or MinHash

### ML
- **TensorFlow.js:** Run in browser alongside WASM
- **Models:** MobileNet for images

---

## 🎯 Recommended First Steps

### Week 1: Web Workers + Progress
```bash
# Create files
touch wasm-worker.js
touch progress-handler.js

# Update Go to report progress
# Update UI to show progress bar
```

### Week 2: IndexedDB Caching
```bash
# Add idb library
npm install idb

# Implement caching layer
# Test with large file sets
```

### Week 3: Image Similarity
```bash
# Implement pHash in Go
# Add image type detection
# Update UI to show similar images
```

---

## 💡 Killer Feature Idea

**"Smart Duplicate Groups"**

Combine all matching methods:

```
Group 1: IMG_1234.jpg (original)
  - IMG_1234.png (exact, different format)
  - IMG_1234_edit.jpg (perceptual match, 95%)
  - IMG_1234_crop.jpg (perceptual match, 88%)
  
Group 2: song.mp3 (original)
  - song.m4a (audio fingerprint match, 98%)
  - song_320.mp3 (audio fingerprint match, 95%)
```

Users see one group with all related duplicates!

---

## 🚀 Quick Win: Implement Web Workers This Week

This gives you:
- ✅ Non-blocking UI
- ✅ Better UX
- ✅ Foundation for all future features
- ✅ Relatively easy to implement

Want me to implement Web Workers first? It's the biggest bang for buck! 🎯
