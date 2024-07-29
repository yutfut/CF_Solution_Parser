package app

import (
	"context"
	"log"
	"scp/config"
	"scp/internal/proxy"
	"sync"
	"time"

	"scp/internal/reader"
	"scp/internal/worker"
	"scp/internal/writer"
)

const (
	path = "./config/config.json"
)

func Run(mainCtx context.Context) {
	ctx, cancel := context.WithCancel(mainCtx)

	conf, err := config.ReadConf(path)
	if err != nil {
		log.Fatal(err)
	}

	input := make(
		chan string,
		conf.Workers.InputChanel,
	)
	output := make(
		chan []string,
		conf.Workers.OutputChanel,
	)

	r := reader.NewReader(
		conf.Files.Input,
		input,
	)

	go r.Read(ctx)

	w := writer.NewWriter(
		conf.Files.Output,
		output,
	)

	go w.Write(ctx)

	wg := &sync.WaitGroup{}

	workerGroup := worker.NewWorker(
		input,
		output,
		wg,
		proxy.NewProxy(
			conf.Proxies,
		),
	)

	for i := 0; i < conf.Workers.WorkerCount; i++ {
		wg.Add(1)

		go workerGroup.Worker()
	}

	wg.Wait()

	cancel()

	time.Sleep(10 * time.Second)
}
