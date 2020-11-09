package site_writer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Meromen/JsonlParser/html_helper"
	"github.com/Meromen/JsonlParser/jsonl_reader"
	"io/ioutil"
	"log"
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
		SiteWriters:    sync.Map{},
	}

	return r
}

func (r *SiteReceiver) Receive(ctx context.Context, siteChan <-chan jsonl_reader.Site) (successProcessed int) {
	for {
		select {
		case <-ctx.Done():
			return
		case site, more := <-siteChan:
			if !more {
				return
			}

			title, description, err := r.getPageTitleAndDescription(site.Url)
			if err != nil {
				log.Println(fmt.Sprintf("Error while working with %s:\n %v",site.Url , err))
				continue
			}

			for _, category := range site.Categories {
				r.writeStrToBufferByCategory(category, fmt.Sprintf("%s\t%s\t%s\n", site.Url, title, description))
			}
			successProcessed++
		}
	}
}

func (r *SiteReceiver) WriteRemainingTextFromBuffers() {
	r.SiteWriters.Range(func(key, val interface{}) bool {
		categoryName := key.(string)
		writer := val.(*SiteWriter)
		buf := writer.GetBuffer()
		err := r.writeToFileByCategory(categoryName, buf)
		if err != nil {
			log.Fatal(err)
		}
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
			err := r.writeToFileByCategory(category, buf)
			if err != nil {
				log.Fatal(err)
			}
			buf.Reset()
			r.BufferPool.Put(buf)
		}()
	}

	writer.Write(str)
}

func (r *SiteReceiver) getPageTitleAndDescription(url string) (title, description string, err error) {
	client := r.HttpClientPool.Get()
	res, err := client.Get(url)
	if err != nil {
		return title, description, errors.New(fmt.Sprintf("Failed to get Title and Description from page: %v", err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return title, description, errors.New(fmt.Sprintf("Failed to get Title and Description from page, HTTP Status Code: %d", res.StatusCode))
	}

	htmlString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return title, description, errors.New(fmt.Sprintf("Failed to get Title and Description from page: %v", err))
	}

	r.HttpClientPool.Put(client)

	title, description, err = html_helper.GetHtmlTitleAndDescription(htmlString)
	if err != nil {
		return title, description, errors.New(fmt.Sprintf("Failde to get Title and Description from page: %v", err))
	}

	return title, description, err
}

func (r *SiteReceiver) writeToFileByCategory(categoryName string, buf *bytes.Buffer) error {
	file, err := os.OpenFile(fmt.Sprintf("%s.tsv", categoryName), os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Failde to write data to file: %v", err))
	}
	defer file.Close()

	str := buf.String()
	_, err = file.WriteString(str)
	if err != nil {
		return errors.New(fmt.Sprintf("Failde to write data to file: %v", err))
	}

	return nil
}
