package certvault

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"example.com/certvault/internal/persist"
	"example.com/certvault/internal/store"
)

// TestBug09_PersistTempFileClosed：Save 后临时文件句柄须关闭，否则 Windows 无法删残留/重入。
func TestBug09_PersistTempFileClosed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := persist.Snapshot{Entries: []store.Entry{{ID: "1", Name: "a", CertPEM: []byte("x")}}, SavedAt: time.Now()}
	if err := persist.Save(path, snap); err != nil {
		t.Fatal(err)
	}
	// 再次 Save 应成功；若上次 tmp 未 Close，部分平台 rename/remove 会失败或留下锁
	if err := persist.Save(path, snap); err != nil {
		t.Fatalf("second save: %v", err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "certvault-*.tmp"))
	if len(matches) != 0 {
		// 尝试删除，若仍被占用则失败
		for _, m := range matches {
			if err := os.Remove(m); err != nil {
				t.Fatalf("tmp still locked (%s): %v", runtime.GOOS, err)
			}
		}
		t.Fatalf("leftover tmp files: %v", matches)
	}
}
