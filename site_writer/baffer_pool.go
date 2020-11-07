package site_writer

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool *sync.Pool
}

func NewBufferPool() BufferPool {
	return BufferPool{
		pool: &sync.Pool{New: func() interface{} {
			return &bytes.Buffer{}
		}},
	}
}

func (m *BufferPool) Get() *bytes.Buffer{
	return m.pool.Get().(*bytes.Buffer)
}

func (m *BufferPool) Put(buffer *bytes.Buffer) {
	m.pool.Put(buffer)
}
