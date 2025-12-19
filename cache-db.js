// cache-db.js - IndexedDB wrapper for file hash caching
const DB_NAME = 'pure-dupes-cache';
const DB_VERSION = 1;
const HASH_STORE = 'file-hashes';
const RESULTS_STORE = 'analysis-results';

class CacheDB {
    constructor() {
        this.db = null;
    }
    
    async init() {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(DB_NAME, DB_VERSION);
            
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                this.db = request.result;
                resolve();
            };
            
            request.onupgradeneeded = (event) => {
                const db = event.target.result;
                
                // Store for individual file hashes
                if (!db.objectStoreNames.contains(HASH_STORE)) {
                    const hashStore = db.createObjectStore(HASH_STORE, {keyPath: 'path'});
                    hashStore.createIndex('hash', 'hash', {unique: false});
                    hashStore.createIndex('size', 'size', {unique: false});
                    hashStore.createIndex('modTime', 'modTime', {unique: false});
                }
                
                // Store for complete analysis results
                if (!db.objectStoreNames.contains(RESULTS_STORE)) {
                    db.createObjectStore(RESULTS_STORE, {keyPath: 'id'});
                }
            };
        });
    }
    
    // Store file hash
    async putFileHash(path, hash, size, modTime, chunks) {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([HASH_STORE], 'readwrite');
            const store = transaction.objectStore(HASH_STORE);
            
            const data = {
                path,
                hash,
                size,
                modTime,
                chunks,
                timestamp: Date.now()
            };
            
            const request = store.put(data);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
        });
    }
    
    // Get file hash
    async getFileHash(path) {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([HASH_STORE], 'readonly');
            const store = transaction.objectStore(HASH_STORE);
            const request = store.get(path);
            
            request.onsuccess = () => resolve(request.result);
            request.onerror = () => reject(request.error);
        });
    }
    
    // Check if file is cached and unchanged
    async isCached(path, size, modTime) {
        const cached = await this.getFileHash(path);
        if (!cached) return false;
        
        // Check if file hasn't changed
        return cached.size === size && cached.modTime === modTime;
    }
    
    // Store complete analysis result
    async putAnalysisResult(id, result) {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([RESULTS_STORE], 'readwrite');
            const store = transaction.objectStore(RESULTS_STORE);
            
            const data = {
                id,
                result,
                timestamp: Date.now()
            };
            
            const request = store.put(data);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
        });
    }
    
    // Get analysis result
    async getAnalysisResult(id) {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([RESULTS_STORE], 'readonly');
            const store = transaction.objectStore(RESULTS_STORE);
            const request = store.get(id);
            
            request.onsuccess = () => resolve(request.result?.result);
            request.onerror = () => reject(request.error);
        });
    }
    
    // Get all cached file hashes
    async getAllHashes() {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([HASH_STORE], 'readonly');
            const store = transaction.objectStore(HASH_STORE);
            const request = store.getAll();
            
            request.onsuccess = () => resolve(request.result);
            request.onerror = () => reject(request.error);
        });
    }
    
    // Clear all cached data
    async clear() {
        const hashTransaction = this.db.transaction([HASH_STORE], 'readwrite');
        await hashTransaction.objectStore(HASH_STORE).clear();
        
        const resultsTransaction = this.db.transaction([RESULTS_STORE], 'readwrite');
        await resultsTransaction.objectStore(RESULTS_STORE).clear();
    }
    
    // Get cache statistics
    async getStats() {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction([HASH_STORE], 'readonly');
            const store = transaction.objectStore(HASH_STORE);
            const countRequest = store.count();
            
            countRequest.onsuccess = () => {
                resolve({
                    cachedFiles: countRequest.result,
                    dbSize: 'N/A' // Browser doesn't provide easy way to get DB size
                });
            };
            countRequest.onerror = () => reject(countRequest.error);
        });
    }
}

// Export singleton instance
const cacheDB = new CacheDB();
