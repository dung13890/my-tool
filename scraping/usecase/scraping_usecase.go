package usecase

import (
	"github.com/dung13890/my-tool/domain"
	"github.com/dung13890/my-tool/scraping/repository"
	"regexp"
)

type scrapingUsecase struct {
	repo domain.TicketRepository
}

func NewScrapingUsecase(url string) domain.TicketUsecase {
	u := &scrapingUsecase{}
	pherusa := regexp.MustCompile(`https://pherusa([-/\.\w\d])*`)
	redmine := regexp.MustCompile(`https://dev.sun-asterisk([-/\.\w\d])*`)
	switch {
	case pherusa.MatchString(url):
		u.repo = repository.NewPherusaRepository(url)
	case redmine.MatchString(url):
		u.repo = repository.NewRedmineRepository(url)
	default:
		return nil
	}

	return u
}

func (u scrapingUsecase) Scraping() (domain.Ticket, error) {
	return u.repo.Scraping()
}
