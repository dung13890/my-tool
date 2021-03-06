package repository

import (
	"errors"
	"fmt"
	"github.com/dung13890/my-tool/domain"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"time"
)

type redmineRepository struct {
	url string
}

func NewRedmineRepository(url string) domain.TicketRepository {
	return &redmineRepository{
		url: url,
	}
}

func (re *redmineRepository) SetUp() (*colly.Collector, error) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("dev.sun-asterisk.com"),
	)
	c.SetRequestTimeout(400 * time.Second)

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  viper.GetString(`redmine.name`),
		Value: viper.GetString(`redmine.value`),
	}
	cookies = append(cookies, cookie)
	// Authenticate
	if err := c.SetCookies(re.url, cookies); err != nil {
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

func (re *redmineRepository) Scraping() (t domain.Ticket, err error) {
	t.Url = re.url
	c, err := re.SetUp()
	if err != nil {
		return
	}

	c.OnHTML(".subject h3", func(e *colly.HTMLElement) {
		t.Title = e.Text
	})

	c.OnHTML(".status.attribute .value", func(e *colly.HTMLElement) {
		t.Status = e.Text
	})

	c.OnHTML("#relations", func(e *colly.HTMLElement) {
		count := 0
		bugs := []string{}
		e.ForEach("table tr.issue", func(_ int, el *colly.HTMLElement) {
			matched, _ := regexp.MatchString(`^.* - Bug #.*$`, el.ChildText("td.subject"))
			if matched {
				if el.ChildText("td.status") != "Reject" {
					count += 1
					bugs = append(bugs, fmt.Sprintf("[%s] - %s", el.ChildText("td.status"), el.ChildText("td.subject")))
				}
			}
		})

		t.BugCount = count
		t.Bugs = bugs
	})

	// Start scraping on url
	c.Visit(re.url)

	return
}
