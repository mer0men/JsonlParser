package jsonl_reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type JsonlSiteReader struct {
	siteChan chan Site
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
			jsonString := scanner.Text()
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
		if err := scanner.Err(); err != nil {
			log.Fatal(fmt.Sprintf("Failed to read file: %v", err))
		}
	}()

	return jr.siteChan, nil
}
