package cmd

import (
	"errors"
	"fmt"

	"github.com/typemd/typemd/core"
)

// resolveVault creates a Vault with the given path, defaulting to "." if empty.
func resolveVault(path string) *core.Vault {
	if path == "" {
		return core.NewVault(".")
	}
	return core.NewVault(path)
}

// resolveObject resolves an object by full ID or prefix, and formats an
// actionable error message when the result is ambiguous.
func resolveObject(vault *core.Vault, id string) (*core.Object, error) {
	obj, err := vault.ResolveObject(id)
	if err == nil {
		return obj, nil
	}

	var ambig *core.ErrAmbiguousObject
	if errors.As(err, &ambig) {
		fmt.Println("Multiple objects match the prefix. Please use the full ID:")
		for _, m := range ambig.Matches {
			fmt.Printf("  %s\n", m.ID)
		}
		return nil, fmt.Errorf("ambiguous prefix: %s", id)
	}

	return nil, err
}
