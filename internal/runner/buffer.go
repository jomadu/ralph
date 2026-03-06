package runner

import (
	"io"
)

// BoundedBuffer captures output with a maximum size, discarding oldest content when full.
type BoundedBuffer struct {
	data []byte
	max  int
}

// NewBoundedBuffer creates a buffer with the specified maximum size in bytes.
func NewBoundedBuffer(maxBytes int) *BoundedBuffer {
	return &BoundedBuffer{
		data: make([]byte, 0, maxBytes),
		max:  maxBytes,
	}
}

// Write implements io.Writer. When data exceeds max, oldest bytes are discarded.
func (b *BoundedBuffer) Write(p []byte) (n int, err error) {
	n = len(p)
	b.data = append(b.data, p...)
	if len(b.data) > b.max {
		b.data = b.data[len(b.data)-b.max:]
	}
	return n, nil
}

// Bytes returns the current buffer contents.
func (b *BoundedBuffer) Bytes() []byte {
	return b.data
}

// MultiWriter creates a writer that duplicates writes to all provided writers.
func MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}
