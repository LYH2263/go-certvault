package validate

const (
	MaxPEMBytes   = 1 << 20
	MaxTags       = 32
	MaxTagLen     = 64
	MaxDescLen    = 512
	MaxDNSNames   = 64
	DefaultCSRDay = 365
	MaxCSRDay     = 3650
)

// ClampDays 限制签发天数。
func ClampDays(d int) int {
	if d <= 0 {
		return DefaultCSRDay
	}
	if d > MaxCSRDay {
		return MaxCSRDay
	}
	return d
}
