package main

import (
	"Meromen/JsonlParser/jsonl_reader"
	"Meromen/JsonlParser/site_writer"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

const workerCount = 10

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	sr := jsonl_reader.JsonlSiteReader{}
	siteChan, err := sr.ReadFileToChannel("500.jsonl")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to read file: %v", err))
	}

	totalSuccessProcessed := 0
	r := site_writer.NewSiteReceiver()
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			successProcessed := r.Receive(ctx, siteChan)
			totalSuccessProcessed +=successProcessed
			wg.Done()
		}()
	}

	wg.Wait()
	r.WriteRemainingTextFromBuffers()

	log.Println(fmt.Sprintf("Receive complete: %d successfuly processed sites", totalSuccessProcessed))
}
