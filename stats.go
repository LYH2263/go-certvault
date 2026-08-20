package certvault

// Stats 返回运行统计。
func (v *Vault) Stats() Stats {
	v.mu.Lock()
	defer v.mu.Unlock()
	st := Stats{
		Imports:   v.imports,
		Revokes:   v.revokes,
		Rotates:   v.rotates,
		Scans:     v.scans,
		Closed:    v.closed,
		PersistOn: v.persistPath != "",
	}
	if v.st != nil {
		st.Entries = v.st.Len()
		st.Revoked = v.st.RevokedCount()
	}
	return st
}
