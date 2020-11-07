package jsonl_reader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type JsonlSiteReader struct {
	siteChan chan Site
}

func (jr *JsonlSiteReader) SiteChanel(c context.Context) <-chan Site {
	return jr.siteChan
}

func (jr *JsonlSiteReader) ReadFileToChannel(filename string) (<-chan Site, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	jr.siteChan = make(chan Site, 1000)

	go func() {
		for scanner.Scan() {
			jsonString :=  scanner.Text()
			if jsonString != "" {
				site := Site{}
				json.Unmarshal([]byte(jsonString), &site)
				if len(site.Categories) == 0 {
					site.Categories = append(site.Categories, "no_category")
				}
				jr.siteChan <- site
			}
		}
		close(jr.siteChan)
		file.Close()
		fmt.Println("complete reading ")
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	return jr.siteChan, nil
}
