package certvault

import (
	"fmt"

	"example.com/certvault/internal/persist"
)

func (v *Vault) persistLocked() error {
	if v.persistPath == "" {
		return nil
	}
	if v.st == nil {
		return ErrClosed
	}
	snap := persist.Snapshot{
		Entries:   v.st.List(),
		Revoked:   v.crl.Snapshot(),
		SavedAt:   v.clk.Now(),
		Imports:   v.imports,
		Revokes:   v.revokes,
		Rotates:   v.rotates,
		Scans:     v.scans,
	}
	if err := persist.Save(v.persistPath, snap); err != nil { // Save must close temp
		return fmt.Errorf("%w: %v", ErrPersist, err)
	}
	return nil
}

// LoadPersist 从快照恢复（启动用）。
func (v *Vault) LoadPersist() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.persistPath == "" {
		return nil
	}
	snap, err := persist.Load(v.persistPath)
	if err != nil {
		return err
	}
	return v.st.ReplaceAll(snap.Entries)
}
