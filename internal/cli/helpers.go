package cli

import (
	"encoding/json"
	"os"

	"profilepress-pp-cli/internal/store"
)

func openStore(path string) (*store.Store, error) { return store.Open(path) }

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeFile(path string, b []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(b)
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
