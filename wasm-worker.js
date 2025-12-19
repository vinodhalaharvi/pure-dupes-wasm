// wasm-worker.js - Web Worker for background WASM processing

// Try to load wasm_exec.js with better error handling
try {
    importScripts('wasm_exec.js');
} catch (err) {
    self.postMessage({
        type: 'error',
        error: 'Failed to load wasm_exec.js. Make sure you ran ./build.sh first! Error: ' + err.message
    });
    throw new Error('wasm_exec.js not found - did you run ./build.sh?');
}

let wasmReady = false;
let go = null;

// Initialize Go runtime
try {
    go = new Go();
} catch (err) {
    self.postMessage({
        type: 'error',
        error: 'Failed to initialize Go runtime: ' + err.message
    });
}

// Load WASM module
if (go) {
    fetch('main.wasm')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch main.wasm - did you run ./build.sh?');
            }
            return response.arrayBuffer();
        })
        .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
        .then(result => {
            go.run(result.instance);
            wasmReady = true;
            self.postMessage({type: 'ready'});
            console.log('✅ WASM Worker ready');
        })
        .catch((err) => {
            console.error('❌ Worker failed to load WASM:', err);
            self.postMessage({
                type: 'error',
                error: 'WASM loading failed: ' + err.message + '. Make sure you ran ./build.sh!'
            });
        });
}

// Handle messages from main thread
self.onmessage = async function(e) {
    const {type, data} = e.data;
    
    if (type === 'analyze') {
        if (!wasmReady) {
            self.postMessage({
                type: 'error',
                error: 'WASM module not ready yet. Please wait a moment and try again.'
            });
            return;
        }
        
        try {
            const {files, threshold, chunkSize} = data;
            
            // Progress callback
            const progressCallback = (progress) => {
                self.postMessage({
                    type: 'progress',
                    data: progress
                });
            };
            
            // Call WASM function with progress callback
            const resultJSON = analyzeFiles(
                files,
                threshold,
                chunkSize,
                progressCallback
            );
            
            const result = JSON.parse(resultJSON);
            
            if (result.error) {
                self.postMessage({
                    type: 'error',
                    error: result.error
                });
            } else {
                self.postMessage({
                    type: 'complete',
                    data: result
                });
            }
        } catch (err) {
            self.postMessage({
                type: 'error',
                error: err.message
            });
        }
    }
};
