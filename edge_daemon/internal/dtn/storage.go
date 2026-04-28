package dtn

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	mu      sync.RWMutex
	path    string
	bundles []Bundle
}

func NewStorage(path string) (*Storage, error) {
	storage := &Storage{path: path, bundles: make([]Bundle, 0)}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &storage.bundles)
	}
	return storage, nil
}

func (storage *Storage) Add(bundle Bundle) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.bundles = append(storage.bundles, bundle)
	return storage.persistLocked()
}

func (storage *Storage) Pending() []Bundle {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	copyBundles := make([]Bundle, len(storage.bundles))
	copy(copyBundles, storage.bundles)
	return copyBundles
}

func (storage *Storage) Remove(ids []string) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	if len(ids) == 0 {
		return nil
	}
	removeSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		removeSet[id] = struct{}{}
	}
	filtered := storage.bundles[:0]
	for _, bundle := range storage.bundles {
		if _, found := removeSet[bundle.ID]; found {
			continue
		}
		filtered = append(filtered, bundle)
	}
	storage.bundles = append([]Bundle(nil), filtered...)
	return storage.persistLocked()
}

func (storage *Storage) persistLocked() error {
	body, err := json.MarshalIndent(storage.bundles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storage.path, body, 0o644)
}
