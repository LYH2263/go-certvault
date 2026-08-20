package validate

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrEmptyName = errors.New("empty name")
	ErrBadName   = errors.New("invalid name")
)

// MetaName 校验业务名。
func MetaName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrEmptyName
	}
	if len(name) > 128 {
		return ErrBadName
	}
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return ErrBadName
	}
	return nil
}

// DNSName 粗检 DNS。
func DNSName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || len(s) > 253 {
		return false
	}
	for _, p := range strings.Split(s, ".") {
		if p == "" || len(p) > 63 {
			return false
		}
	}
	return true
}
