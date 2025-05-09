package utils

import (
	"os"
	"path/filepath"
	"time"
)

func CleanOldWasmCache(homeDir string, maxAge time.Duration) {
	cacheDir := filepath.Join(homeDir, "wasm", "wasm", "cache")
	_ = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if time.Since(info.ModTime()) > maxAge {
			os.Remove(path)
		}
		return nil
	})
}
