package core

import (
	"fmt"
	"strings"
)

// ErrAmbiguousObject is returned when a prefix matches more than one object.
type ErrAmbiguousObject struct {
	Prefix   string
	Matches  []*Object
}

func (e *ErrAmbiguousObject) Error() string {
	names := make([]string, len(e.Matches))
	for i, o := range e.Matches {
		names[i] = o.ID
	}
	return fmt.Sprintf("ambiguous prefix %q matches %d objects: %s", e.Prefix, len(e.Matches), strings.Join(names, ", "))
}

// ResolveObject resolves an object ID that may be a full ID or a name prefix.
//
// Resolution order:
//  1. Exact match on full ID (e.g. "book/clean-code-01jqr3k5mp...")
//  2. Prefix match: if the input looks like "type/prefix", find objects of that
//     type whose filename starts with the given prefix.
//
// Returns ErrAmbiguousObject if multiple objects match the prefix.
// Returns an error wrapping os.ErrNotExist if no match is found.
func (v *Vault) ResolveObject(id string) (*Object, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	// 1. Exact match — try reading the file directly.
	obj, err := v.GetObject(id)
	if err == nil {
		return obj, nil
	}

	// 2. Prefix match via SQLite.
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid object ID: %s", id)
	}
	typeName, prefix := parts[0], parts[1]

	rows, err := v.db.Query(
		"SELECT id, type, filename, properties, body FROM objects WHERE type = ? AND filename LIKE ?",
		typeName, prefix+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("prefix query: %w", err)
	}
	defer rows.Close()

	matches, err := scanObjects(rows)
	if err != nil {
		return nil, fmt.Errorf("scan prefix results: %w", err)
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("object not found: %s", id)
	case 1:
		// Re-read from file to get the full content (body etc.)
		return v.GetObject(matches[0].ID)
	default:
		return nil, &ErrAmbiguousObject{Prefix: id, Matches: matches}
	}
}
