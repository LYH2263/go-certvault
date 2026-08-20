package certvault

import (
	"context"
	"time"

	"example.com/certvault/internal/expire"
)

// ScanExpiring 扫描 within 内到期的活动证书；尊重 ctx 取消。
func (v *Vault) ScanExpiring(within time.Duration) ([]ExpiringHit, error) {
	return v.ScanExpiringContext(context.Background(), within)
}

// ScanExpiringContext 带取消的到期扫描。
func (v *Vault) ScanExpiringContext(ctx context.Context, within time.Duration) ([]ExpiringHit, error) {
	if within < 0 {
		return nil, ErrBadRequest
	}
	v.mu.Lock()
	if v.closed || v.st == nil {
		v.mu.Unlock()
		return nil, ErrClosed
	}
	ents := v.st.List()
	now := v.clk.Now()
	step := v.scanStep
	v.mu.Unlock()

	hits, err := expire.Scan(ctx, ents, now, within, step, func(e expire.Candidate) ExpiringHit {
		return ExpiringHit{
			ID:        e.ID,
			Name:      e.Name,
			NotAfter:  e.NotAfter,
			Remaining: e.Remaining,
		}
	})
	if err != nil {
		return nil, mapCtxErr(err)
	}
	v.mu.Lock()
	v.scans++
	v.mu.Unlock()
	return hits, nil
}

func mapCtxErr(err error) error {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return ErrCanceled
	}
	return err
}
