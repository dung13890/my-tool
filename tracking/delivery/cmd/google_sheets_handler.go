package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/my-tool/domain"
	scraperUsecase "github.com/dung13890/my-tool/scraping/usecase"
	trackingUsecase "github.com/dung13890/my-tool/tracking/usecase"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"sync"
	"time"
)

type googleSheetsHandler struct {
	googleSheetsUsecase domain.GoogleSheetsUsecase
	ticketUsecase       domain.TicketUsecase
}

func NewGoogleSheet() *cli.Command {
	return &cli.Command{
		Name:    "googleSheet",
		Aliases: []string{"g"},
		Usage:   "google sheets tracking",
		Action: func(ctx *cli.Context) error {
			u, err := trackingUsecase.NewGoogleSheetsUsecase()
			if err != nil {
				log.Fatal(err)
			}
			g := &googleSheetsHandler{
				googleSheetsUsecase: u,
			}
			g.exec()
			return nil
		},
	}
}

func (g *googleSheetsHandler) exec() error {
	// Instantiate spinner processes
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
	sp.Suffix = fmt.Sprintf(" [%s]: Processing...", "scraper")
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]: Completed!\n", green("✔"), "scraper")
	sp.Start()
	// Fetch item from spreadsheets
	items, err := g.googleSheetsUsecase.Tracking()
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Update bug into spreadsheets
	if len(items) > 0 {
		var wg sync.WaitGroup
		notifyChan := make(chan domain.GoogleSheets, len(items))
		errChan := make(chan error, len(items))

		for _, item := range items {
			wg.Add(1)
			go g.update(item, &wg, notifyChan, errChan)
		}
		wg.Wait()
		close(notifyChan)
		close(errChan)

		if len(notifyChan) > 0 {
			var notifyItem []domain.GoogleSheets
			for nc := range notifyChan {
				notifyItem = append(notifyItem, nc)
			}
			// Notify chatwork when bug > 5
			g.googleSheetsUsecase.Notify(notifyItem)
		}
		if len(errChan) > 0 {
			for e := range errChan {
				log.Fatal(e)
			}
		}

	}
	// handler usecase
	sp.Stop()

	return nil
}

func (g *googleSheetsHandler) update(item *domain.GoogleSheets, wg *sync.WaitGroup, ch chan domain.GoogleSheets, errC chan error) error {
	defer wg.Done()
	g.ticketUsecase = scraperUsecase.NewScrapingUsecase(item.Ticket)
	ticket, err := g.ticketUsecase.Scraping()
	if err != nil {
		return err
	}
	if ticket.BugCount >= viper.GetInt(`spreadsheet.bug_limit`) && ticket.BugCount != item.CurrentBug {
		item.CurrentBug = ticket.BugCount
		ch <- *item
	}
	if err := g.googleSheetsUsecase.Update(item, ticket.BugCount); err != nil {
		errC <- err
	}

	return nil
}
