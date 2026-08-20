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
	// 任一失败路径都要关闭并清理临时文件，避免 Windows 上的句柄泄漏与残留锁。
	// rename 成功后 tmpName 已不存在，Close/Remove 均为无副作用操作。
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
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
