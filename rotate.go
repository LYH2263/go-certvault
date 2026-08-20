package certvault

import (
	"context"
	"fmt"

	"example.com/certvault/internal/pemutil"
	"example.com/certvault/internal/rotate"
	"example.com/certvault/internal/store"
)

// Rotate 用新证替换旧证；持久化失败则回滚。
func (v *Vault) Rotate(oldID string, certPEM, keyPEM []byte) (string, error) {
	return v.RotateContext(context.Background(), oldID, certPEM, keyPEM)
}

// RotateContext 可取消的轮换。
func (v *Vault) RotateContext(ctx context.Context, oldID string, certPEM, keyPEM []byte) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", mapCtxErr(err)
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return "", ErrClosed
	}
	old, ok := v.st.Get(oldID)
	if !ok {
		return "", ErrNotFound
	}
	if old.Revoked {
		return "", ErrRevoked
	}
	delay := v.rotateIODelay
	waiter := v.ioWait
	v.mu.Unlock()
	if delay > 0 && waiter != nil {
		if err := waiter.Wait(ctx, delay); err != nil {
			v.mu.Lock()
			return "", mapCtxErr(err)
		}
	}
	v.mu.Lock()
	if v.closed || v.st == nil {
		return "", ErrClosed
	}
	parsed, err := pemutil.ParsePair(certPEM, keyPEM)
	if err != nil {
		return "", err
	}
	neu := store.Entry{
		Name:        old.Name,
		Tags:        append([]string(nil), old.Tags...),
		Description: old.Description,
		Active:      true,
		CertPEM:     pemutil.Clone(certPEM),
		KeyPEM:      pemutil.Clone(keyPEM),
		Cert:        parsed.Cert,
		Key:         parsed.Key,
		Subject:     parsed.Subject,
		Issuer:      parsed.Issuer,
		Serial:      parsed.Serial,
		NotBefore:   parsed.NotBefore,
		NotAfter:    parsed.NotAfter,
		Fingerprint: parsed.Fingerprint,
		ImportedAt:  v.clk.Now(),
	}
	tx := rotate.Begin(v.st, oldID)
	newID, err := tx.Apply(neu)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRotateFailed, err)
	}
	if err := v.persistLocked(); err != nil {
		// 持久化失败：回滚事务，活动指针保持原证。
		_ = tx.Rollback()
		return "", fmt.Errorf("%w: %v", ErrPersist, err)
	}
	tx.Commit()
	v.rotates++
	return newID, nil
}

// PreviewRotate 不落库的轮换预览。
func (v *Vault) PreviewRotate(oldID string, certPEM, keyPEM []byte) (RotatePreview, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return RotatePreview{}, ErrClosed
	}
	old, ok := v.st.Get(oldID)
	if !ok {
		return RotatePreview{}, ErrNotFound
	}
	parsed, err := pemutil.ParsePair(certPEM, keyPEM)
	if err != nil {
		return RotatePreview{}, err
	}
	return RotatePreview{
		OldID:     oldID,
		NewSerial: parsed.Serial,
		OldAfter:  old.NotAfter,
		NewAfter:  parsed.NotAfter,
	}, nil
}
