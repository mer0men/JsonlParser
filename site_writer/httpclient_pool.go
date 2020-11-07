package site_writer

import (
	"net/http"
	"sync"
	"time"
)

type HttpClientPool struct {
	pool *sync.Pool
}

func NewHttpClientPool() HttpClientPool {
	return HttpClientPool{
		pool: &sync.Pool{New: func() interface{} {
			return &http.Client{Timeout: time.Second * 30}
		}},
	}
}

func (m *HttpClientPool) Get() *http.Client{
	return m.pool.Get().(*http.Client)
}

func (m *HttpClientPool) Put(client *http.Client) {
	m.pool.Put(client)
}
