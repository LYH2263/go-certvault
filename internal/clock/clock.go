package clock

import "time"

// Clock 可注入时钟。
type Clock interface {
	Now() time.Time
}

// Real 系统时钟。
type Real struct{}

func (Real) Now() time.Time { return time.Now() }
