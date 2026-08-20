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
	_ = ctx // BUG: ignore cancel
	if d <= 0 {
		return nil
	}
	time.Sleep(d)
	return nil
}
