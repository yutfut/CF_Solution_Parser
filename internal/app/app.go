package app

import (
	"context"
	"sync"
	"time"

	"scp/internal/cfclient"
	"scp/internal/reader"
	"scp/internal/worker"
	"scp/internal/writer"
)

func Run(mainCtx context.Context) {
	ctx, cancel := context.WithCancel(mainCtx)

	input := make(chan string, 20)
	output := make(chan []string, 20)

	r := reader.NewReader("./submissions.csv", input)

	go r.Read(ctx)

	w := writer.NewWriter("./data.csv", output)

	go w.Write(ctx)

	wg := &sync.WaitGroup{}

	workerGroup := worker.NewWorker(
		input,
		output,
		wg,
		cfclient.NewCFClient(),
	)

	wg.Add(1)

	go workerGroup.Worker()

	wg.Wait()

	cancel()

	time.Sleep(10 * time.Second)
}
