package core

import "fmt"

// Relation represents a relationship record between two objects.
type Relation struct {
	Name   string
	FromID string
	ToID   string
}

// ListRelations returns all relations where objectID is either the source or target.
func (v *Vault) ListRelations(objectID string) ([]Relation, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	rows, err := v.db.Query(
		"SELECT name, from_id, to_id FROM relations WHERE from_id = ? OR to_id = ?",
		objectID, objectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list relations: %w", err)
	}
	defer rows.Close()

	var rels []Relation
	for rows.Next() {
		var r Relation
		if err := rows.Scan(&r.Name, &r.FromID, &r.ToID); err != nil {
			return nil, fmt.Errorf("scan relation: %w", err)
		}
		rels = append(rels, r)
	}
	return rels, rows.Err()
}

// findRelationProperty finds a relation property by name in a type schema.
func findRelationProperty(schema *TypeSchema, name string) *Property {
	for i, p := range schema.Properties {
		if p.Name == name && p.Type == "relation" {
			return &schema.Properties[i]
		}
	}
	return nil
}

// LinkObjects creates a relation between two objects.
func (v *Vault) LinkObjects(fromID, relName, toID string) error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	// Get source object and schema
	fromObj, err := v.GetObject(fromID)
	if err != nil {
		return fmt.Errorf("get source object: %w", err)
	}

	fromSchema, err := v.LoadType(fromObj.Type)
	if err != nil {
		return fmt.Errorf("load source type: %w", err)
	}

	relProp := findRelationProperty(fromSchema, relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	// Validate target object
	toObj, err := v.GetObject(toID)
	if err != nil {
		return fmt.Errorf("get target object: %w", err)
	}
	if relProp.Target != "" && toObj.Type != relProp.Target {
		return fmt.Errorf("target type mismatch: expected %q, got %q", relProp.Target, toObj.Type)
	}

	// Update source frontmatter
	if relProp.Multiple {
		existing, _ := fromObj.Properties[relName]
		var arr []any
		if existing != nil {
			if existArr, ok := existing.([]any); ok {
				arr = existArr
			}
		}
		for _, item := range arr {
			if item == toID {
				return fmt.Errorf("relation already exists: %s -[%s]-> %s", fromID, relName, toID)
			}
		}
		arr = append(arr, toID)
		fromObj.Properties[relName] = arr
	} else {
		fromObj.Properties[relName] = toID
	}

	if err := v.writeObjectProperties(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}

	// Write to relations table
	_, err = v.db.Exec(
		"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		relName, fromID, toID,
	)
	if err != nil {
		return fmt.Errorf("insert relation: %w", err)
	}

	// Handle bidirectional
	if relProp.Bidirectional && relProp.Inverse != "" {
		toSchema, err := v.LoadType(toObj.Type)
		if err != nil {
			return fmt.Errorf("load target type: %w", err)
		}

		inverseProp := findRelationProperty(toSchema, relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		if inverseProp.Multiple {
			existing, _ := toObj.Properties[relProp.Inverse]
			var arr []any
			if existing != nil {
				if existArr, ok := existing.([]any); ok {
					arr = existArr
				}
			}
			arr = append(arr, fromID)
			toObj.Properties[relProp.Inverse] = arr
		} else {
			toObj.Properties[relProp.Inverse] = fromID
		}

		if err := v.writeObjectProperties(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}

		_, err = v.db.Exec(
			"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
			relProp.Inverse, toID, fromID,
		)
		if err != nil {
			return fmt.Errorf("insert inverse relation: %w", err)
		}
	}

	return nil
}

// UnlinkObjects removes a relation between two objects.
// If both is true and the relation is bidirectional, also removes the inverse.
func (v *Vault) UnlinkObjects(fromID, relName, toID string, both bool) error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	// Get source object and schema
	fromObj, err := v.GetObject(fromID)
	if err != nil {
		return fmt.Errorf("get source object: %w", err)
	}

	fromSchema, err := v.LoadType(fromObj.Type)
	if err != nil {
		return fmt.Errorf("load source type: %w", err)
	}

	relProp := findRelationProperty(fromSchema, relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	// Remove from source frontmatter
	if relProp.Multiple {
		existing, _ := fromObj.Properties[relName]
		if existArr, ok := existing.([]any); ok {
			var newArr []any
			for _, item := range existArr {
				if item != toID {
					newArr = append(newArr, item)
				}
			}
			if len(newArr) == 0 {
				fromObj.Properties[relName] = nil
			} else {
				fromObj.Properties[relName] = newArr
			}
		}
	} else {
		fromObj.Properties[relName] = nil
	}

	if err := v.writeObjectProperties(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}

	// Delete from relations table
	_, err = v.db.Exec(
		"DELETE FROM relations WHERE name = ? AND from_id = ? AND to_id = ?",
		relName, fromID, toID,
	)
	if err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}

	// Handle --both with bidirectional
	if both && relProp.Bidirectional && relProp.Inverse != "" {
		toObj, err := v.GetObject(toID)
		if err != nil {
			return fmt.Errorf("get target object: %w", err)
		}

		toSchema, err := v.LoadType(toObj.Type)
		if err != nil {
			return fmt.Errorf("load target type: %w", err)
		}

		inverseProp := findRelationProperty(toSchema, relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		if inverseProp.Multiple {
			existing, _ := toObj.Properties[relProp.Inverse]
			if existArr, ok := existing.([]any); ok {
				var newArr []any
				for _, item := range existArr {
					if item != fromID {
						newArr = append(newArr, item)
					}
				}
				if len(newArr) == 0 {
					toObj.Properties[relProp.Inverse] = nil
				} else {
					toObj.Properties[relProp.Inverse] = newArr
				}
			}
		} else {
			toObj.Properties[relProp.Inverse] = nil
		}

		if err := v.writeObjectProperties(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}

		_, err = v.db.Exec(
			"DELETE FROM relations WHERE name = ? AND from_id = ? AND to_id = ?",
			relProp.Inverse, toID, fromID,
		)
		if err != nil {
			return fmt.Errorf("delete inverse relation: %w", err)
		}
	}

	return nil
}
