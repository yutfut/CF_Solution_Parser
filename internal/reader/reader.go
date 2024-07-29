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

type Reader struct {
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

	return &Reader{
		file:   file,
		output: output,
	}
}

func (r *Reader) Read(ctx context.Context) {
	reader := csv.NewReader(r.file)

	for {
		select {
		case <-ctx.Done():
			log.Println("Reader cancelled")
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

func (r *Reader) Close() {
	if err := r.file.Close(); err != nil {
		log.Fatal(err)
	}
}
