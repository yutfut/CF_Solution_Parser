package writer

import (
	"context"
	"encoding/csv"
	"log"
	"os"
)

type WriterInterface interface {
	Write(ctx context.Context)
	Close()
}

type writer struct {
	file  *os.File
	input chan []string
}

func NewWriter(
	filePath string,
	input chan []string,
) WriterInterface {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	return &writer{
		file:  file,
		input: input,
	}
}

func (w *writer) Write(ctx context.Context) {
	writer := csv.NewWriter(w.file)

	if err := writer.Write(
		[]string{
			"id",
			"solution",
		},
	); err != nil {
		panic(err)
	}

	for item := range w.input {
		select {
		case <-ctx.Done():
			log.Println("Closing writer")
			w.Close()
			return
		default:
			if err := writer.Write(
				item,
			); err != nil {
				panic(err)
			}
		}
	}
}

func (w *writer) Close() {
	if err := w.file.Close(); err != nil {
		log.Fatal(err)
	}
}
