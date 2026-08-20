package certvault

// Close 先刷盘再关闭；之后写入返回 ErrClosed。
func (v *Vault) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return nil
	}
	// 先刷盘：此时 store 仍在，快照才能含已导入条目。
	err := v.persistLocked()
	v.closed = true
	v.st = nil // 刷盘完成后再丢弃 store。
	return err
}

// Closed 是否已关闭。
func (v *Vault) Closed() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.closed
}
