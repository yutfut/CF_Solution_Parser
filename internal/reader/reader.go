package reader

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type ReaderInterface interface {
	Read(ctx context.Context)
	Close()
}

type reader struct {
	file   *os.File
	output chan string
}

func NewReader(
	filePath string,
	output chan string,
) ReaderInterface {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	return &reader{
		file:   file,
		output: output,
	}
}

func (r *reader) Read(ctx context.Context) {
	reader := csv.NewReader(r.file)

	for {
		select {
		case <-ctx.Done():
			log.Println("reader cancelled")
			r.Close()
			return
		default:
			record, err := reader.Read()
			if err != nil {
				fmt.Println(err)
				break
			}
			if record[1] == "id" {
				continue
			}

			r.output <- record[1]
		}
	}
}

func (r *reader) Close() {
	if err := r.file.Close(); err != nil {
		log.Fatal(err)
	}
}
