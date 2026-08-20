package expire

import (
	"fmt"
	"time"
)

// Windows 常用扫描窗口。
func Windows() []time.Duration {
	return []time.Duration{
		24 * time.Hour,
		7 * 24 * time.Hour,
		30 * 24 * time.Hour,
		90 * 24 * time.Hour,
	}
}

// Label 窗口标签。
func Label(d time.Duration) string {
	days := int(d / (24 * time.Hour))
	if days < 1 {
		return "lt1d"
	}
	return fmt.Sprintf("%dd", days)
}

// Days 转 duration。
func Days(n int) time.Duration {
	if n < 0 {
		n = 0
	}
	return time.Duration(n) * 24 * time.Hour
}
