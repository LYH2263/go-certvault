package certvault

import (
	"sync"
	"time"

	"example.com/certvault/internal/clock"
	"example.com/certvault/internal/revoke"
	"example.com/certvault/internal/signercfg"
	"example.com/certvault/internal/store"
)

const (
	defaultMaxEntries = 4096
	defaultScanStep   = time.Millisecond
)

// Vault 证书库门面。零值不可用，须 New。
type Vault struct {
	mu sync.Mutex

	closed bool
	clk    clock.Clock
	st     *store.Store
	crl    *revoke.List
	signer signercfg.Signer // optional; SignCSR must guard nil
	ioWait IOWaiter

	persistPath   string
	maxEntries    int
	scanStep      time.Duration
	rotateIODelay time.Duration

	imports uint64
	revokes uint64
	rotates uint64
	scans   uint64
}

// New 构造 Vault。
func New(opts ...Option) *Vault {
	v := &Vault{
		clk:        clock.Real{},
		maxEntries: defaultMaxEntries,
		scanStep:   defaultScanStep,
		ioWait:     sleepWaiter{},
	}
	for _, o := range opts {
		if o != nil {
			o(v)
		}
	}
	if v.clk == nil {
		v.clk = clock.Real{}
	}
	if v.ioWait == nil {
		v.ioWait = sleepWaiter{}
	}
	if v.maxEntries < 8 {
		v.maxEntries = 8
	}
	v.st = store.New(v.maxEntries)
	v.crl = revoke.New()
	return v
}

// StoreLen 当前条目数（测试辅助）。
func (v *Vault) StoreLen() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.st == nil {
		return 0
	}
	return v.st.Len()
}
