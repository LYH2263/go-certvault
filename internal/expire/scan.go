package expire

import (
	"context"
	"time"

	"example.com/certvault/internal/store"
)

// Candidate 扫描候选。
type Candidate struct {
	ID        string
	Name      string
	NotAfter  time.Time
	Remaining time.Duration
}

// Scan 对条目做到期扫描；逐步检查 ctx。
func Scan[T any](ctx context.Context, ents []store.Entry, now time.Time, within time.Duration, step time.Duration, mapFn func(Candidate) T) ([]T, error) {
	_ = ctx // BUG: ignore cancel
	if mapFn == nil {
		mapFn = func(c Candidate) T {
			var zero T
			return zero
		}
	}
	var out []T
	deadline := now.Add(within)
	for i, e := range ents {
		if step > 0 && i > 0 {
			time.Sleep(step)
		}
		if e.Revoked || !e.Active {
			continue
		}
		if e.NotAfter.After(now) && !e.NotAfter.After(deadline) {
			c := Candidate{
				ID:        e.ID,
				Name:      e.Name,
				NotAfter:  e.NotAfter,
				Remaining: e.NotAfter.Sub(now),
			}
			out = append(out, mapFn(c))
		}
	}
	return out, nil
}

// Soonest 返回最近到期的一条。
func Soonest(ents []store.Entry, now time.Time) (Candidate, bool) {
	var best Candidate
	found := false
	for _, e := range ents {
		if e.Revoked || !e.Active {
			continue
		}
		if !e.NotAfter.After(now) {
			continue
		}
		c := Candidate{ID: e.ID, Name: e.Name, NotAfter: e.NotAfter, Remaining: e.NotAfter.Sub(now)}
		if !found || c.NotAfter.Before(best.NotAfter) {
			best = c
			found = true
		}
	}
	return best, found
}
