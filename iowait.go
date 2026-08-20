package certvault

import (
	"context"
	"time"
)

// IOWaiter 抽象轮换等路径上的可取消等待。
type IOWaiter interface {
	Wait(ctx context.Context, d time.Duration) error
}

type sleepWaiter struct{}

func (sleepWaiter) Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
