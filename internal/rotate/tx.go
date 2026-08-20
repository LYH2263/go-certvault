package rotate

import (
	"example.com/certvault/internal/store"
)

// Tx 轮换事务：Apply 后可 Rollback 或 Commit。
type Tx struct {
	st       *store.Store
	oldID    string
	old      store.Entry
	newID    string
	applied  bool
	committed bool
}

func Begin(st *store.Store, oldID string) *Tx {
	old, _ := st.Get(oldID)
	return &Tx{st: st, oldID: oldID, old: old}
}

// Apply 写入新证并停用旧证。
func (t *Tx) Apply(neu store.Entry) (string, error) {
	neu.Active = true
	id, err := t.st.Put(neu)
	if err != nil {
		return "", err
	}
	old := t.old
	old.Active = false
	if err := t.st.Replace(t.oldID, old); err != nil {
		_ = t.st.Delete(id)
		return "", err
	}
	t.newID = id
	t.applied = true
	return id, nil
}

// Rollback 撤销 Apply。
func (t *Tx) Rollback() error {
	if !t.applied || t.committed {
		return nil
	}
	if t.newID != "" {
		_ = t.st.Delete(t.newID)
	}
	old := t.old
	old.Active = true
	_ = t.st.Replace(t.oldID, old)
	t.applied = false
	return nil
}

// Commit 确认事务。
func (t *Tx) Commit() {
	t.committed = true
}
