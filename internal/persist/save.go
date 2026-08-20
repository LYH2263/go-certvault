package persist

import (
	"encoding/json"
	"os"
	"path/filepath"

	"example.com/certvault/internal/store"
)

// Save 原子写快照：临时文件 → Sync → Close → Rename。
func Save(path string, snap Snapshot) error {
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(sanitize(snap), "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "certvault-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	// BUG: missing tmp.Close() before rename
	return os.Rename(tmpName, path)
}

func sanitize(snap Snapshot) Snapshot {
	out := snap
	ents := make([]store.Entry, 0, len(snap.Entries))
	for _, e := range snap.Entries {
		e.Cert = nil
		e.Key = nil
		ents = append(ents, e)
	}
	out.Entries = ents
	return out
}
