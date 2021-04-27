package scraping

import (
	"github.com/dung13890/my-tool/entities"
	"github.com/gocolly/colly"
)

type Usecase interface {
	SetUp() (*colly.Collector, error)
	Scraping() (t entities.Ticket, err error)
}
