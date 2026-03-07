package core

import (
	"errors"
	"strings"
	"testing"
)

func TestVault_ResolveObject_ExactMatch(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "golang-in-action")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	obj, err := v.ResolveObject(created.ID)
	if err != nil {
		t.Fatalf("ResolveObject() error = %v", err)
	}
	if obj.ID != created.ID {
		t.Errorf("ID = %q, want %q", obj.ID, created.ID)
	}
}

func TestVault_ResolveObject_PrefixMatch(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "golang-in-action")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// prefix without ULID
	obj, err := v.ResolveObject("book/golang-in-action")
	if err != nil {
		t.Fatalf("ResolveObject() prefix error = %v", err)
	}
	if obj.ID != created.ID {
		t.Errorf("ID = %q, want %q", obj.ID, created.ID)
	}
}

func TestVault_ResolveObject_PartialPrefix(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "golang-in-action")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	obj, err := v.ResolveObject("book/golang")
	if err != nil {
		t.Fatalf("ResolveObject() partial prefix error = %v", err)
	}
	if obj.ID != created.ID {
		t.Errorf("ID = %q, want %q", obj.ID, created.ID)
	}
}

func TestVault_ResolveObject_Ambiguous(t *testing.T) {
	v := setupTestVault(t)

	if _, err := v.NewObject("book", "go-in-action"); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	if _, err := v.NewObject("book", "go-programming"); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	_, err := v.ResolveObject("book/go")
	if err == nil {
		t.Fatal("expected ErrAmbiguousObject, got nil")
	}

	var ambig *ErrAmbiguousObject
	if !errors.As(err, &ambig) {
		t.Fatalf("expected ErrAmbiguousObject, got %T: %v", err, err)
	}
	if len(ambig.Matches) != 2 {
		t.Errorf("Matches len = %d, want 2", len(ambig.Matches))
	}
	if !strings.Contains(ambig.Error(), "book/go") {
		t.Errorf("error message should contain prefix, got: %s", ambig.Error())
	}
}

func TestVault_ResolveObject_NotFound(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.ResolveObject("book/nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent object, got nil")
	}
}

func TestVault_ResolveObject_ExactTakesPriority(t *testing.T) {
	v := setupTestVault(t)

	// Create two objects: one named "go", one named "go-advanced"
	obj1, err := v.NewObject("book", "go")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	_, err = v.NewObject("book", "go-advanced")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// Exact match on full ID should always win
	obj, err := v.ResolveObject(obj1.ID)
	if err != nil {
		t.Fatalf("ResolveObject() exact ID error = %v", err)
	}
	if obj.ID != obj1.ID {
		t.Errorf("ID = %q, want %q", obj.ID, obj1.ID)
	}
}

func TestVault_ResolveObject_DBNotOpen(t *testing.T) {
	v := setupTestVault(t)
	v.Close()

	_, err := v.ResolveObject("book/test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_ResolveObject_InvalidID(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.ResolveObject("no-slash-here")
	if err == nil {
		t.Fatal("expected error for invalid ID format, got nil")
	}
}
