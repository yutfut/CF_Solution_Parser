package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"scp/internal/cfclient"
	"sync"
)

func Worker(
	chanelIN chan string,
	chanelOUT chan []string,
) {
	client := cfclient.NewCFClient()

	for item := range chanelIN {

		solution, err := client.GetSolution(
			item,
		)
		if err != nil {
			return
		}

		chanelOUT <- []string{item, solution.Source}
	}
}

func ReadCSV(
	chanel chan string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	file, err := os.Open("./submissions.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		if record[1] == "id" {
			continue
		}

		chanel <- record[1]
	}

	close(chanel)
}

func WriteCSV(
	channel chan []string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	file, err := os.Create("./data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err = writer.Write(
		[]string{
			"id",
			"solution",
		},
	); err != nil {
		panic(err)
	}

	for item := range channel {
		if err = writer.Write(
			item,
		); err != nil {
			panic(err)
		}
	}

	close(channel)
}

func main() {
	chanelIN := make(chan string, 20)
	chanelOUT := make(chan []string, 20)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go ReadCSV(chanelIN, wg)

	go Worker(chanelIN, chanelOUT)

	go WriteCSV(chanelOUT, wg)

	wg.Wait()
}
