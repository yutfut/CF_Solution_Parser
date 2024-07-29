package cfclient

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"strings"
)

type CFClientInterface interface {
	GetSolution(item string) (Solution, error)
}

type CFClient struct {
	client  *resty.Client
	csrf    string
	cookies []*http.Cookie
}

func NewCFClient() CFClientInterface {
	cfClient := &CFClient{
		client: resty.New().
			SetDebug(false).
			SetProxy("217.29.63.91:11792:HJEUfj:paz7pN"),
	}

	if err := cfClient.getMeta(); err != nil {
		log.Fatal(err)
	}

	return cfClient
}

func (cf *CFClient) getMeta() error {
	tokenResp, err := cf.client.R().
		Get("https://codeforces.com")
	if err != nil {
		log.Printf("client.R() ::: %+v", err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(tokenResp.Body())))

	list := htmlquery.Find(doc, "//meta[@name='X-Csrf-Token']")

	if len(list) == 0 {
		return errors.New("no X-Csrf-Token found")
	}

	cf.csrf = list[0].Attr[len(list[0].Attr)-1].Val
	cf.cookies = tokenResp.Cookies()

	return nil
}

func (cf *CFClient) GetSolution(item string) (Solution, error) {
	for {
		solution := Solution{}

		resp, err := cf.client.R().
			SetHeaders(
				map[string]string{
					"X-Csrf-Token": cf.csrf,
					"Referer":      "https://codeforces.com/problemset/status",
				},
			).
			SetFormData(
				map[string]string{
					"submissionId": item,
					"csrf_token":   cf.csrf,
				},
			).SetCookies(cf.cookies).
			SetResult(&solution).
			Post("https://codeforces.com/data/submitSource")
		if err != nil {
			log.Printf("%s: %t", item, false)
			log.Printf("client.R() ::: %+v", err)
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf(
				"%s ::: %d",
				item,
				resp.StatusCode(),
			)
		}

		if solution.Source == "" {
			log.Printf("%s: %t ::: %s\n", item, false, resp.Body())
			return Solution{}, errors.New("solution not found")
		} else {
			log.Printf("%s: %s", item, solution.Source)
			return solution, nil
		}
	}
}
