package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// filterCondition represents a single key=value query condition.
type filterCondition struct {
	Key   string
	Value string
}

// parseFilter parses a filter string like "type=book status=reading"
// into a slice of filterConditions. Empty string returns nil, nil.
// Returns an error if any token is not in key=value form.
func parseFilter(filter string) ([]filterCondition, error) {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil, nil
	}

	parts := strings.Fields(filter)
	conditions := make([]filterCondition, 0, len(parts))

	for _, part := range parts {
		idx := strings.Index(part, "=")
		if idx < 1 {
			return nil, fmt.Errorf("invalid filter condition: %q (expected key=value)", part)
		}
		conditions = append(conditions, filterCondition{
			Key:   part[:idx],
			Value: part[idx+1:],
		})
	}

	return conditions, nil
}

// scanObjects scans a sql.Rows result into a slice of Objects.
func scanObjects(rows *sql.Rows) ([]*Object, error) {
	var results []*Object
	for rows.Next() {
		var obj Object
		var propsJSON string
		if err := rows.Scan(&obj.ID, &obj.Type, &obj.Filename, &propsJSON, &obj.Body); err != nil {
			return nil, fmt.Errorf("scan object: %w", err)
		}
		if err := json.Unmarshal([]byte(propsJSON), &obj.Properties); err != nil {
			return nil, fmt.Errorf("unmarshal properties: %w", err)
		}
		results = append(results, &obj)
	}
	return results, rows.Err()
}

// QueryObjects queries objects using key=value filter syntax.
// Multiple conditions are combined with AND.
// "type" is a special key that filters on the objects.type column.
// Other keys filter on JSON properties using json_extract.
// An empty filter returns all objects.
func (v *Vault) QueryObjects(filter string) ([]*Object, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	conditions, err := parseFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("parse filter: %w", err)
	}

	query := "SELECT id, type, filename, properties, body FROM objects"
	var whereClauses []string
	var args []any

	for _, c := range conditions {
		if c.Key == "type" {
			whereClauses = append(whereClauses, "type = ?")
			args = append(args, c.Value)
		} else {
			whereClauses = append(whereClauses, "json_extract(properties, ?) = ?")
			args = append(args, "$."+c.Key, c.Value)
		}
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := v.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query objects: %w", err)
	}
	defer rows.Close()

	return scanObjects(rows)
}

// SearchObjects performs full-text search using FTS5.
// Searches across filename, properties, and body.
// Returns nil, nil for empty keyword.
func (v *Vault) SearchObjects(keyword string) ([]*Object, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, nil
	}

	query := `SELECT o.id, o.type, o.filename, o.properties, o.body
		FROM objects o
		JOIN objects_fts fts ON o.rowid = fts.rowid
		WHERE objects_fts MATCH ?`

	rows, err := v.db.Query(query, keyword)
	if err != nil {
		return nil, fmt.Errorf("search objects: %w", err)
	}
	defer rows.Close()

	return scanObjects(rows)
}

// RebuildIndex rebuilds the FTS5 index from the objects table.
// Useful when the index may be out of sync with the objects table.
func (v *Vault) RebuildIndex() error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	_, err := v.db.Exec("INSERT INTO objects_fts(objects_fts) VALUES('rebuild')")
	if err != nil {
		return fmt.Errorf("rebuild index: %w", err)
	}

	return nil
}
