package worker

import (
	"scp/internal/cfclient"
	"sync"
)

type WorkerInterface interface {
	Worker()
}

type Worker struct {
	input  chan string
	output chan []string
	wg     *sync.WaitGroup
	client cfclient.CFClientInterface
}

func NewWorker(
	input chan string,
	output chan []string,
	wg *sync.WaitGroup,
	client cfclient.CFClientInterface,
) WorkerInterface {
	return &Worker{
		input:  input,
		output: output,
		wg:     wg,
		client: client,
	}
}

func (w *Worker) Worker() {
	defer w.wg.Done()

	for item := range w.input {

		solution, err := w.client.GetSolution(
			item,
		)
		if err != nil {
			return
		}

		w.output <- []string{item, solution.Source}
	}
}
