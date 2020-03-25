package runtime

//go:export malloc
func malloc(size uint) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
