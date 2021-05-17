package repository

import (
	"errors"
	"fmt"
	"github.com/dung13890/my-tool/domain"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"net/http"
	_ "regexp"
	"strconv"
	"time"
)

type backlogRepository struct {
	url string
}

func NewBacklogRepository(url string) domain.TicketRepository {
	return &backlogRepository{
		url: url,
	}
}

func (b *backlogRepository) SetUp() (*colly.Collector, error) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("framgiabrg.backlog.com"),
	)
	c.SetRequestTimeout(400 * time.Second)

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  viper.GetString(`backlog.name`),
		Value: viper.GetString(`backlog.value`),
	}
	cookies = append(cookies, cookie)
	// Authenticate
	if err := c.SetCookies(b.url, cookies); err != nil {
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

func (b *backlogRepository) Scraping() (t domain.Ticket, err error) {
	t.Url = b.url
	c, err := b.SetUp()
	if err != nil {
		return
	}

	c.OnHTML(".content-main .ticket__title-group h2 span", func(e *colly.HTMLElement) {
		t.Title = e.Text
	})
	c.OnHTML(".content-main .ticket__header", func(e *colly.HTMLElement) {
		t.Status = e.ChildText("span.status-label")
	})
	c.OnHTML(".content-main .ticket__properties", func(e *colly.HTMLElement) {
		e.ForEach(".ticket__properties-item", func(_ int, el *colly.HTMLElement) {
			// Fetch priority
			if el.Attr("class") == "ticket__properties-item -priority" {
				t.Priority = el.ChildText(".ticket__properties-value")
			}

			// Fetch category
			if el.Attr("class") == "ticket__properties-item -category" {
				t.Category = el.ChildText(".ticket__properties-value")
			}

			// Fetch version
			if el.Attr("class") == "ticket__properties-item -milestones" {
				t.Version = el.ChildText(".ticket__properties-value")
			}

			// Fetch EST
			if el.Attr("class") == "ticket__properties-item -estimated-hours" {
				t.EstimatedTime, _ = strconv.Atoi(el.ChildText(".ticket__properties-value"))
			}

			// Fetch AcT
			if el.Attr("class") == "ticket__properties-item -actual-hours" {
				t.ActualTime, _ = strconv.Atoi(el.ChildText(".ticket__properties-value"))
			}
		})
	})
	// Start scraping on url
	c.Visit(b.url)

	return
}
