package cfclient

import (
	"log"
	"net/http"
	"strings"

	"scp/internal/proxy"
	"scp/pkg/errors"

	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
)

type CFClientInterface interface {
	GetSolution(item string) (Solution, error)
}

type cfClient struct {
	client       *resty.Client
	csrf         string
	cookies      []*http.Cookie
	proxyManager proxy.ProxyInterface
}

func NewCFClient(
	proxyManager proxy.ProxyInterface,
) CFClientInterface {
	client := &cfClient{
		client: resty.New().
			SetDebug(false).
			SetProxy(
				proxyManager.Get(),
			),
		proxyManager: proxyManager,
	}

	if err := client.getMeta(); err != nil {
		log.Fatal(err)
	}

	return client
}

func (cf *cfClient) getMeta() error {
	tokenResp, err := cf.client.R().
		Get(getMeta)
	if err != nil {
		log.Printf("client.R() ::: %+v", err)
		return err
	}

	doc, err := htmlquery.Parse(
		strings.NewReader(
			string(
				tokenResp.Body(),
			),
		),
	)

	list := htmlquery.Find(
		doc,
		xQuery,
	)

	if len(list) == 0 {
		return errors.CSRFTokenNotFoundError
	}

	cf.csrf = list[0].Attr[len(list[0].Attr)-1].Val
	cf.cookies = tokenResp.Cookies()

	return nil
}

func (cf *cfClient) recreateClient() {
	cf.client = resty.New().
		SetDebug(false).
		SetProxy(
			cf.proxyManager.Get(),
		)

	if err := cf.getMeta(); err != nil {
		log.Fatal(err)
	}
}

func (cf *cfClient) GetSolution(item string) (Solution, error) {
	for {
		solution := Solution{}

		resp, err := cf.client.R().
			SetHeaders(
				map[string]string{
					csrfHeader:    cf.csrf,
					refererHeader: referer,
				},
			).
			SetFormData(
				map[string]string{
					submissionId: item,
					csrfToken:    cf.csrf,
				},
			).SetCookies(cf.cookies).
			SetResult(&solution).
			Post(getSubmit)
		if err != nil {
			log.Printf(
				"%s: %t\nclient.R() ::: %+v",
				item,
				false,
				err,
			)
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
		} else {
			log.Printf("%s: %s", item, solution.Source)
			return solution, nil
		}

		cf.recreateClient()
	}
}
