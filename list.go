package certvault

// List 返回全部条目视图（含私钥拷贝）。
func (v *Vault) List() []EntryView {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return nil
	}
	ents := v.st.List()
	out := make([]EntryView, 0, len(ents))
	for _, e := range ents {
		out = append(out, toView(e, true))
	}
	return out
}

// ListPublic 列表但不导出私钥内容。
func (v *Vault) ListPublic() []EntryView {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return nil
	}
	ents := v.st.List()
	out := make([]EntryView, 0, len(ents))
	for _, e := range ents {
		out = append(out, toView(e, false))
	}
	return out
}
