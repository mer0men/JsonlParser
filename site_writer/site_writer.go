package site_writer

import (
	"bytes"
	"sync"
)

type SiteWriter struct {
	m      *sync.Mutex
	buffer *bytes.Buffer
}

func (w *SiteWriter) SetBuffer(buffer *bytes.Buffer) {
	w.m.Lock()
	w.buffer = buffer
	w.m.Unlock()
}

func (w *SiteWriter) GetBuffer() *bytes.Buffer {
	return w.buffer
}

func (w *SiteWriter) Write(s string) {
	w.m.Lock()
	w.buffer.WriteString(s)
	w.m.Unlock()
}

func (w *SiteWriter) Length() int {
	w.m.Lock()
	l := w.buffer.Len()
	w.m.Unlock()
	return l
}
