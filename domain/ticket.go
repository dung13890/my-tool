package domain

import (
	"github.com/gocolly/colly"
)

// Entity
type Ticket struct {
	Url              string
	Name             string
	Version          string
	Title            string
	EstimatedTime    int
	ActualTime       int
	BugCount         int
	PercentCompleted int
	Priority         string
	Category         string
	Bugs             []string
	Status           string
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
