package store

import "strings"

// FindByFingerprint 按指纹查找。
func (s *Store) FindByFingerprint(fp string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fp = strings.ToLower(fp)
	for _, e := range s.byID {
		if strings.ToLower(e.Fingerprint) == fp {
			return e.Clone(), true
		}
	}
	return Entry{}, false
}

// FindBySerial 按序列号查找。
func (s *Store) FindBySerial(serial string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	serial = strings.ToUpper(serial)
	for _, e := range s.byID {
		if strings.ToUpper(e.Serial) == serial {
			return e.Clone(), true
		}
	}
	return Entry{}, false
}

// Names 去重名称列表。
func (s *Store) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := map[string]struct{}{}
	var out []string
	for _, id := range s.order {
		e := s.byID[id]
		if _, ok := seen[e.Name]; ok {
			continue
		}
		seen[e.Name] = struct{}{}
		out = append(out, e.Name)
	}
	return out
}
