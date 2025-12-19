// wasm-worker.js - Web Worker for background WASM processing
importScripts('wasm_exec.js');

let wasmReady = false;
const go = new Go();

// Load WASM module
WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
    .then((result) => {
        go.run(result.instance);
        wasmReady = true;
        self.postMessage({type: 'ready'});
        console.log('✅ WASM Worker ready');
    })
    .catch((err) => {
        console.error('❌ Worker failed to load WASM:', err);
        self.postMessage({type: 'error', error: err.message});
    });

// Handle messages from main thread
self.onmessage = async function(e) {
    const {type, data} = e.data;
    
    if (type === 'analyze') {
        if (!wasmReady) {
            self.postMessage({
                type: 'error',
                error: 'WASM module not ready yet'
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
