package site_writer

import (
	"Meromen/JsonlParser/html_helper"
	"Meromen/JsonlParser/jsonl_reader"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const maxBufferLength = 4096

type SiteReceiver struct {
	BufferPool     *BufferPool
	HttpClientPool *HttpClientPool
	SiteWriters    sync.Map
}

func NewSiteReceiver() *SiteReceiver {
	bufferPool := NewBufferPool()
	httpClientPool := NewHttpClientPool()
	r := &SiteReceiver{
		BufferPool:     &bufferPool,
		HttpClientPool: &httpClientPool,
		SiteWriters:   	sync.Map{},
	}

	return r
}

func (r *SiteReceiver) Receive(ctx context.Context, siteChan <-chan jsonl_reader.Site) {
	for {
		select {
		case <-ctx.Done():
			return
		case site, more := <-siteChan:
			if !more {
				return
			}
			fmt.Printf("%+v\n", site)

			title, description, err := r.getPageTitleAndDescription(site.Url)
			if err != nil {
				continue
			}

			for _, category := range site.Categories {
				r.writeStrToBufferByCategory(category, fmt.Sprintf("%s\t%s\t%s\n", site.Url, title, description))
			}
		}
	}
}

func (r *SiteReceiver) WriteRemainingTextFromBuffers()  {
	r.SiteWriters.Range(func(key, val interface{})bool {
		categoryName := key.(string)
		writer := val.(*SiteWriter)
		buf := writer.GetBuffer()
		r.writeToFileByCategory(categoryName, buf)
		return true
	})
}

func (r *SiteReceiver) writeStrToBufferByCategory(category, str string) {
	var writer *SiteWriter
	val, ok := r.SiteWriters.Load(category)
	if ok {
		writer = val.(*SiteWriter)
	} else {
		writer = &SiteWriter{
			m:      &sync.Mutex{},
			buffer: r.BufferPool.Get(),
		}
		r.SiteWriters.Store(category, writer)
	}


	if writer.Length() > maxBufferLength {
		buf := writer.GetBuffer()
		writer.SetBuffer(r.BufferPool.Get())
		go func() {
			r.writeToFileByCategory(category, buf)
			buf.Reset()
			r.BufferPool.Put(buf)
		}()
	}

	writer.Write(str)
	fmt.Println("received")
}

func (r *SiteReceiver) getPageTitleAndDescription(url string) (title string, description string, err error) {
	client := r.HttpClientPool.Get()
	res, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	htmlString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	r.HttpClientPool.Put(client)

	title = html_helper.GetHtmlTitle(string(htmlString))
	description = html_helper.GetHtmlDescription(string(htmlString))

	return title, description, err
}

func (r *SiteReceiver) writeToFileByCategory(categoryName string, buf *bytes.Buffer) {
	file, err := os.OpenFile(fmt.Sprintf("%s.tsv", categoryName), os.O_RDWR | os.O_CREATE | os.O_APPEND | os.O_SYNC, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	str := buf.String()
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}
}
