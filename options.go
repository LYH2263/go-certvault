package certvault

import (
	"time"

	"example.com/certvault/internal/clock"
	"example.com/certvault/internal/signercfg"
)

// Option 配置 Vault。
type Option func(*Vault)

// WithClock 注入时钟（测试用 Fake）。
func WithClock(c clock.Clock) Option {
	return func(v *Vault) {
		if c != nil {
			v.clk = c
		}
	}
}

// WithPersistPath 启用 JSON 快照持久化。
func WithPersistPath(path string) Option {
	return func(v *Vault) { v.persistPath = path }
}

// WithMaxEntries 限制条目上限。
func WithMaxEntries(n int) Option {
	return func(v *Vault) {
		if n > 0 {
			v.maxEntries = n
		}
	}
}

// WithSigner 配置 CSR 签名器（自签或外部 CA 占位）。
func WithSigner(s signercfg.Signer) Option {
	return func(v *Vault) { v.signer = s }
}

// WithScanStep 扫描步进（用于可取消长扫描）。
func WithScanStep(d time.Duration) Option {
	return func(v *Vault) {
		if d > 0 {
			v.scanStep = d
		}
	}
}

// WithRotateIODelay 轮换模拟 I/O 等待（测试注入）。
func WithRotateIODelay(d time.Duration) Option {
	return func(v *Vault) { v.rotateIODelay = d }
}

// WithIOWaiter 注入可取消等待器。
func WithIOWaiter(w IOWaiter) Option {
	return func(v *Vault) {
		if w != nil {
			v.ioWait = w
		}
	}
}
