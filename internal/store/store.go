package store

import (
	"errors"
	"sync"

	"example.com/certvault/internal/idgen"
)

var (
	ErrFull      = errors.New("store: full")
	ErrNotFound  = errors.New("store: not found")
	ErrDuplicate = errors.New("store: duplicate")
)

// Store 内存证书库存。
type Store struct {
	mu      sync.RWMutex
	max     int
	byID    map[string]Entry
	order   []string
}

func New(max int) *Store {
	if max < 1 {
		max = 1
	}
	return &Store{
		max:  max,
		byID: make(map[string]Entry),
	}
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.byID)
}

func (s *Store) RevokedCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, e := range s.byID {
		if e.Revoked {
			n++
		}
	}
	return n
}

func (s *Store) Put(e Entry) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.byID) >= s.max {
		return "", ErrFull
	}
	if e.ID == "" {
		e.ID = idgen.New("crt")
	}
	if _, ok := s.byID[e.ID]; ok {
		return "", ErrDuplicate
	}
	if e.Active && e.Name != "" {
		for id, old := range s.byID {
			if old.Name == e.Name && old.Active && !old.Revoked {
				old.Active = false
				s.byID[id] = old
			}
		}
	}
	s.byID[e.ID] = e.Clone()
	s.order = append(s.order, e.ID)
	return e.ID, nil
}

func (s *Store) Get(id string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.byID[id]
	if !ok {
		return Entry{}, false
	}
	return e.Clone(), true
}

func (s *Store) GetActive(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, id := range s.order {
		e := s.byID[id]
		if e.Name == name && e.Active && !e.Revoked {
			return e.Clone(), true
		}
	}
	return Entry{}, false
}

func (s *Store) List() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.order))
	for _, id := range s.order {
		if e, ok := s.byID[id]; ok {
			out = append(out, e.Clone())
		}
	}
	return out
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return ErrNotFound
	}
	delete(s.byID, id)
	n := s.order[:0]
	for _, x := range s.order {
		if x != id {
			n = append(n, x)
		}
	}
	s.order = n
	return nil
}

func (s *Store) MarkRevoked(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	if !ok {
		return ErrNotFound
	}
	e.Revoked = true
	e.Active = false
	s.byID[id] = e
	return nil
}

func (s *Store) SetActive(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	if !ok {
		return ErrNotFound
	}
	for oid, old := range s.byID {
		if old.Name == e.Name && old.Active {
			old.Active = false
			s.byID[oid] = old
		}
	}
	e.Active = true
	s.byID[id] = e
	return nil
}

func (s *Store) Replace(id string, e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return ErrNotFound
	}
	e.ID = id
	s.byID[id] = e.Clone()
	return nil
}

func (s *Store) ReplaceAll(ents []Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID = make(map[string]Entry, len(ents))
	s.order = s.order[:0]
	for _, e := range ents {
		s.byID[e.ID] = e.Clone()
		s.order = append(s.order, e.ID)
	}
	return nil
}

// PeekKeyPEM 返回库存内私钥切片（故意不拷贝——仅内部 rotate 使用时应 Clone）。
// 对外 API 必须走 Clone。
func (s *Store) PeekKeyPEM(id string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.byID[id]
	if !ok {
		return nil, false
	}
	return e.KeyPEM, true
}
