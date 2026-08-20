package revoke

import (
	"sync"
	"time"
)

// Record 吊销记录。
type Record struct {
	Serial    string
	EntryID   string
	Reason    string
	RevokedAt time.Time
}

// List 吊销清单。
type List struct {
	mu   sync.Mutex
	recs []Record
}

func New() *List { return &List{} }

func (l *List) Add(r Record) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.recs = append(l.recs, r)
}

func (l *List) Snapshot() []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Record, len(l.recs))
	copy(out, l.recs)
	return out
}

func (l *List) Contains(serial string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, r := range l.recs {
		if r.Serial == serial {
			return true
		}
	}
	return false
}

func (l *List) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.recs)
}
