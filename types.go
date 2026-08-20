package certvault

import "time"

// Meta 导入时的业务元数据。
type Meta struct {
	Name        string
	Tags        []string
	Description string
	Active      bool
}

// EntryView 对外可见的证书条目视图（私钥仅长度）。
type EntryView struct {
	ID          string
	Name        string
	Subject     string
	Issuer      string
	Serial      string
	NotBefore   time.Time
	NotAfter    time.Time
	Tags        []string
	Description string
	Active      bool
	Revoked     bool
	CertPEM     []byte
	KeyPEM      []byte
	KeyLen      int
	Fingerprint string
}

// ExpiringHit 到期扫描命中。
type ExpiringHit struct {
	ID        string
	Name      string
	NotAfter  time.Time
	Remaining time.Duration
}

// RotatePreview 轮换预览。
type RotatePreview struct {
	OldID     string
	NewSerial string
	OldAfter  time.Time
	NewAfter  time.Time
}

// Stats 运行计数。
type Stats struct {
	Imports   uint64
	Revokes   uint64
	Rotates   uint64
	Scans     uint64
	Entries   int
	Revoked   int
	Closed    bool
	PersistOn bool
}
