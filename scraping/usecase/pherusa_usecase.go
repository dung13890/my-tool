package usecase

import (
	"errors"
	"fmt"
	"github.com/dung13890/my-tool/entities"
	"github.com/dung13890/my-tool/scraping"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"time"
)

type pherusaUsecase struct {
	url string
}

func NewPherusaUsecase(url string) scraping.Usecase {
	return &pherusaUsecase{
		url: url,
	}
}

func (p *pherusaUsecase) SetUp() (*colly.Collector, error) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("pherusa-redmine.sun-asterisk.vn"),
	)
	c.SetRequestTimeout(400 * time.Second)

	// Read config from json config
	viper.SetConfigFile(`infrastructure/config.json`)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: not exists or is wrong json format", err))
	}

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  viper.GetString(`pherusa.name`),
		Value: viper.GetString(`pherusa.value`),
	}
	cookies = append(cookies, cookie)
	// Authenticate
	if err := c.SetCookies(p.url, cookies); err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: have errors from cookies", err))
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(errors.New(fmt.Sprintf("Errors: have errors from response", err)))
		return
	})

	return c, nil
}

func (p *pherusaUsecase) Scraping() (t entities.Ticket, err error) {
	t.Url = p.url
	c, err := p.SetUp()
	if err != nil {
		return
	}

	c.OnHTML(".subject h3", func(e *colly.HTMLElement) {
		t.Title = e.Text
	})
	c.OnHTML(".status.attribute .value", func(e *colly.HTMLElement) {
		t.Status = e.Text
	})
	c.OnHTML("#issue_tree", func(e *colly.HTMLElement) {
		count := 0
		bugs := []string{}
		e.ForEach("table tr.issue", func(_ int, el *colly.HTMLElement) {
			matched, _ := regexp.MatchString(`.*QA Bug.*`, el.ChildText("td.subject"))
			if matched {
				count += 1
				bugs = append(bugs, fmt.Sprintf("[%s] - %s", el.ChildText("td.status"), el.ChildText("td.subject")))
			}
		})

		t.BugCount = count
		t.Bugs = bugs
	})

	// Start scraping on url
	c.Visit(p.url)

	return
}
