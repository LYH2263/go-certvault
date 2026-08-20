package certvault

import "example.com/certvault/internal/revoke"

// Revoke 吊销条目并记入 CRL。
func (v *Vault) Revoke(id, reason string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return ErrClosed
	}
	ent, ok := v.st.Get(id)
	if !ok {
		return ErrNotFound
	}
	if ent.Revoked {
		return ErrRevoked
	}
	if err := v.st.MarkRevoked(id); err != nil {
		return err
	}
	v.crl.Add(revoke.Record{
		Serial:    ent.Serial,
		EntryID:   id,
		Reason:    reason,
		RevokedAt: v.clk.Now(),
	})
	v.revokes++
	return v.persistLocked()
}

// RevocationList 导出吊销记录。
func (v *Vault) RevocationList() []revoke.Record {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.crl == nil {
		return nil
	}
	return v.crl.Snapshot()
}
