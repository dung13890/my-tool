package domain

import (
	"github.com/gocolly/colly"
)

// Entity
type Ticket struct {
	Url      string
	Title    string
	BugCount int
	Bugs     []string
	Status   string
}

// Usecase
type TicketUsecase interface {
	Scraping() (Ticket, error)
}

// Repository
type TicketRepository interface {
	SetUp() (*colly.Collector, error)
	Scraping() (Ticket, error)
}
