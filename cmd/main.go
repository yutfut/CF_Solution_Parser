package main

import (
	"encoding/csv"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Solution struct {
	Source string `json:"source"`
}

func GetMeta() (string, []*http.Cookie) {
	client := resty.New().SetDebug(false)

	tokenResp, err := client.R().Get("https://codeforces.com")
	if err != nil {
		log.Printf("client.R() ::: %+v", err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(tokenResp.Body())))

	list := htmlquery.Find(doc, "//meta[@name='X-Csrf-Token']")

	CSRF := list[0].Attr[len(list[0].Attr)-1].Val

	return CSRF, tokenResp.Cookies()
}

func Worker(
	chanelIN chan string,
	chanelOUT chan []string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	client := resty.New().
		SetDebug(false).
		SetProxy("http://h3v3n5:motBAz@217.29.53.100:11937")

	CSRF, cookies := GetMeta()

	solution := Solution{}

	for item := range chanelIN {
		resp, err := client.R().
			SetHeaders(
				map[string]string{
					"X-Csrf-Token": CSRF,
					"Referer":      "https://codeforces.com/problemset/status",
				},
			).
			SetFormData(
				map[string]string{
					"submissionId": item,
					"csrf_token":   CSRF,
				},
			).SetCookies(cookies).
			SetResult(&solution).
			Post("https://codeforces.com/data/submitSource")
		if err != nil {
			log.Printf("client.R() ::: %+v", err)
		}

		if solution.Source == "" {
			log.Printf("client.R() ::: %+v\n", resp.Body())
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
		fmt.Println(record[1])

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
		fmt.Println("writer", item)
		if err = writer.Write(
			item,
		); err != nil {
			panic(err)
		}
	}

	close(channel)
}

func main() {
	chanelIN := make(chan string)
	chanelOUT := make(chan []string)

	wg := &sync.WaitGroup{}

	wg.Add(3)

	go ReadCSV(chanelIN, wg)

	go Worker(chanelIN, chanelOUT, wg)

	go WriteCSV(chanelOUT, wg)

	wg.Wait()
}
