package pemutil

// Clone 复制字节切片；nil 保持 nil。
func Clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

// CloneNonNil 空也返回新切片。
func CloneNonNil(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
