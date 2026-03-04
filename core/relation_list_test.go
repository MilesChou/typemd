package core

import "testing"

func TestVault_ListRelations_Empty(t *testing.T) {
	v := setupRelationTestVault(t)
	v.NewObject("book", "test")

	rels, err := v.ListRelations("book/test")
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 0 {
		t.Errorf("len(rels) = %d, want 0", len(rels))
	}
}

func TestVault_ListRelations_Forward(t *testing.T) {
	v := setupRelationTestVault(t)
	v.NewObject("book", "test")
	v.NewObject("person", "alan")
	v.LinkObjects("book/test", "author", "person/alan")

	// Bidirectional link creates 2 DB rows: forward (author) + inverse (books)
	// ListRelations returns both since book/test is from_id on forward and to_id on inverse
	rels, err := v.ListRelations("book/test")
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 2 {
		t.Fatalf("len(rels) = %d, want 2", len(rels))
	}

	// Verify forward relation exists
	found := false
	for _, r := range rels {
		if r.Name == "author" && r.FromID == "book/test" && r.ToID == "person/alan" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected forward relation {author, book/test, person/alan} in %+v", rels)
	}
}

func TestVault_ListRelations_BothDirections(t *testing.T) {
	v := setupRelationTestVault(t)
	v.NewObject("book", "test")
	v.NewObject("person", "alan")
	v.LinkObjects("book/test", "author", "person/alan")

	// person/alan should have both: inverse (books, from_id=person/alan) + forward (author, to_id=person/alan)
	rels, err := v.ListRelations("person/alan")
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 2 {
		t.Fatalf("len(rels) = %d, want 2", len(rels))
	}

	// Verify inverse relation exists
	found := false
	for _, r := range rels {
		if r.Name == "books" && r.FromID == "person/alan" && r.ToID == "book/test" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected inverse relation {books, person/alan, book/test} in %+v", rels)
	}
}

func TestVault_ListRelations_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.ListRelations("book/test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}
