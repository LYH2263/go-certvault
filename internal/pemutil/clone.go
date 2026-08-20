package pemutil

// Clone 复制字节切片；nil 保持 nil。
func Clone(b []byte) []byte {
	// BUG: 返回原切片，导致调用方与库存共享底层数组。
	return b
}

// CloneNonNil 空也返回新切片。
func CloneNonNil(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
