package persist

import (
	"time"

	"example.com/certvault/internal/revoke"
	"example.com/certvault/internal/store"
)

// Snapshot 持久化快照（不含解析后的 Cert/Key 指针）。
type Snapshot struct {
	Entries []store.Entry  `json:"entries"`
	Revoked []revoke.Record `json:"revoked"`
	SavedAt time.Time      `json:"saved_at"`
	Imports uint64         `json:"imports"`
	Revokes uint64         `json:"revokes"`
	Rotates uint64         `json:"rotates"`
	Scans   uint64         `json:"scans"`
}
