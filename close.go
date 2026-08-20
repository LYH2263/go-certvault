package certvault

// Close 刷盘后关闭；之后写入返回 ErrClosed。
func (v *Vault) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return nil
	}
	err := v.persistLocked()
	v.closed = true
	// 保留 st 以便只读诊断；写入路径检查 closed。
	return err
}

// Closed 是否已关闭。
func (v *Vault) Closed() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.closed
}
