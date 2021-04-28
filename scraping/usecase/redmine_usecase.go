package usecase

import (
	"github.com/dung13890/my-tool/entities"
	"github.com/dung13890/my-tool/scraping"
	"github.com/gocolly/colly"
)

type redmineUsecase struct {
	url string
}

func NewRedmineUsecase(url string) scraping.Usecase {
	return &redmineUsecase{
		url: url,
	}
}

func (r *redmineUsecase) SetUp() (*colly.Collector, error) {
	return nil, nil
}

func (r *redmineUsecase) Scraping() (t entities.Ticket, err error) {
	return
}
