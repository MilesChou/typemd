package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SyncIndex scans the objects directory, upserts all found objects into the DB,
// removes DB entries for deleted files, and rebuilds the FTS index.
func (v *Vault) SyncIndex() error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	// Collect all object IDs found on disk
	diskIDs := make(map[string]bool)

	objsDir := v.ObjectsDir()
	if _, err := os.Stat(objsDir); os.IsNotExist(err) {
		// No objects directory — just clean DB and return
		_, err := v.db.Exec("DELETE FROM objects")
		if err != nil {
			return fmt.Errorf("clean objects: %w", err)
		}
		return v.RebuildIndex()
	}

	// Walk objects/<type>/<name>.md
	err := filepath.Walk(objsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		rel, err := filepath.Rel(objsDir, path)
		if err != nil {
			return nil
		}

		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		if len(parts) != 2 {
			return nil
		}
		typeName := parts[0]
		filename := strings.TrimSuffix(parts[1], ".md")
		id := typeName + "/" + filename

		// Read and parse file
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		props, body, err := parseFrontmatter(data)
		if err != nil {
			return nil // skip unparseable files
		}

		propsJSON, err := json.Marshal(props)
		if err != nil {
			return nil
		}

		// Upsert into DB
		_, err = v.db.Exec(`
			INSERT INTO objects (id, type, filename, properties, body)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				type = excluded.type,
				filename = excluded.filename,
				properties = excluded.properties,
				body = excluded.body
		`, id, typeName, filename, string(propsJSON), body)
		if err != nil {
			return fmt.Errorf("upsert object %s: %w", id, err)
		}

		diskIDs[id] = true
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk objects: %w", err)
	}

	// Remove DB entries that no longer exist on disk
	rows, err := v.db.Query("SELECT id FROM objects")
	if err != nil {
		return fmt.Errorf("list db objects: %w", err)
	}
	defer rows.Close()

	var toDelete []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan id: %w", err)
		}
		if !diskIDs[id] {
			toDelete = append(toDelete, id)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate db objects: %w", err)
	}

	for _, id := range toDelete {
		if _, err := v.db.Exec("DELETE FROM objects WHERE id = ?", id); err != nil {
			return fmt.Errorf("delete stale object %s: %w", id, err)
		}
	}

	return v.RebuildIndex()
}
