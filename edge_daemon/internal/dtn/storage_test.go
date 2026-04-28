package dtn

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStorageAddAndRemove(t *testing.T) {
	storage, err := NewStorage(filepath.Join(t.TempDir(), "bundles.json"))
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	bundle := Bundle{ID: "bundle-1", CreatedAt: time.Now().UTC()}
	if err := storage.Add(bundle); err != nil {
		t.Fatalf("add bundle: %v", err)
	}
	if len(storage.Pending()) != 1 {
		t.Fatalf("expected 1 pending bundle")
	}
	if err := storage.Remove([]string{"bundle-1"}); err != nil {
		t.Fatalf("remove bundle: %v", err)
	}
	if len(storage.Pending()) != 0 {
		t.Fatalf("expected storage to be empty")
	}
}
