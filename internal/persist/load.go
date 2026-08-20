package persist

import (
	"encoding/json"
	"os"
)

// Load 读取快照。
func Load(path string) (Snapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Exists 文件是否存在。
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
