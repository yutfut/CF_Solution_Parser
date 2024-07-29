package worker

import (
	"scp/internal/cfclient"
	"scp/internal/proxy"
	"sync"
)

type WorkerInterface interface {
	Worker()
}

type worker struct {
	input        chan string
	output       chan []string
	wg           *sync.WaitGroup
	proxyManager proxy.ProxyInterface
}

func NewWorker(
	input chan string,
	output chan []string,
	wg *sync.WaitGroup,
	proxyManager proxy.ProxyInterface,
) WorkerInterface {
	return &worker{
		input:        input,
		output:       output,
		wg:           wg,
		proxyManager: proxyManager,
	}
}

func (w *worker) Worker() {
	defer w.wg.Done()

	client := cfclient.NewCFClient(
		w.proxyManager,
	)

	for item := range w.input {

		solution, err := client.GetSolution(
			item,
		)
		if err != nil {
			return
		}

		w.output <- []string{item, solution.Source}
	}
}
