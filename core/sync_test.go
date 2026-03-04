package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVault_SyncIndex_NewFile(t *testing.T) {
	v := setupTestVault(t)

	// Manually create an object file (bypassing NewObject API)
	typeDir := filepath.Join(v.ObjectsDir(), "book")
	os.MkdirAll(typeDir, 0755)
	os.WriteFile(filepath.Join(typeDir, "test-book.md"), []byte("---\ntitle: Test Book\n---\n\nHello world.\n"), 0644)

	// Also need type schema for the type directory to be valid
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	if err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	// Verify object is now in DB
	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("len(objs) = %d, want 1", len(objs))
	}
	if objs[0].ID != "book/test-book" {
		t.Errorf("ID = %q, want %q", objs[0].ID, "book/test-book")
	}
}

func TestVault_SyncIndex_UpdatedFile(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	// Create via API (body is empty in DB)
	v.NewObject("book", "test-book")

	// Manually edit the file to add body
	objPath := v.ObjectPath("book", "test-book")
	os.WriteFile(objPath, []byte("---\ntitle: Updated\n---\n\nNew body content.\n"), 0644)

	if err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("len(objs) = %d, want 1", len(objs))
	}
	if strings.TrimSpace(objs[0].Body) != "New body content." {
		t.Errorf("Body = %q, want %q", objs[0].Body, "New body content.")
	}
}

func TestVault_SyncIndex_DeletedFile(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	// Create via API
	v.NewObject("book", "test-book")

	// Delete the file
	os.Remove(v.ObjectPath("book", "test-book"))

	if err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 0 {
		t.Errorf("len(objs) = %d, want 0 (deleted file should be removed from DB)", len(objs))
	}
}

func TestVault_SyncIndex_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	err := v.SyncIndex()
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}
