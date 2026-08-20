package rotate

import "time"

// Diff 轮换差异摘要。
type Diff struct {
	OldSerial string
	NewSerial string
	OldAfter  time.Time
	NewAfter  time.Time
	Extended  time.Duration
}

// ComputeDiff 计算差异。
func ComputeDiff(oldSerial, newSerial string, oldAfter, newAfter time.Time) Diff {
	return Diff{
		OldSerial: oldSerial,
		NewSerial: newSerial,
		OldAfter:  oldAfter,
		NewAfter:  newAfter,
		Extended:  newAfter.Sub(oldAfter),
	}
}
