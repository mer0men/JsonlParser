package main

import (
	jsonl_reader "Meromen/JsonlParser/jsonl_reader"
	site_writer "Meromen/JsonlParser/site_writer"
	"context"
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
		panic(err)
	}

	r := site_writer.NewSiteReceiver()
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			r.Receive(ctx, siteChan)
			wg.Done()
		}()
	}

	wg.Wait()

	r.WriteRemainingTextFromBuffers()
	//r.Receive(context.Background(), siteChan)
}
