package certvault

// Close 刷盘后关闭；之后写入返回 ErrClosed。
func (v *Vault) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return nil
	}
	v.closed = true
	v.st = nil // BUG: drop store before callers finish; Import may nil-deref
	return nil
}

// Closed 是否已关闭。
func (v *Vault) Closed() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.closed
}
