package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/my-tool/domain"
	"github.com/dung13890/my-tool/scraping/usecase"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

type scrapingHandler struct {
	usecase domain.TicketUsecase
}

func NewScraping() *cli.Command {
	return &cli.Command{
		Name:    "scraping",
		Aliases: []string{"s"},
		Usage:   "scraping site",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() > 0 {
				url := ctx.Args().Get(0)
				s := &scrapingHandler{
					usecase: usecase.NewScrapingUsecase(url),
				}
				s.exec()
			}
			return nil
		},
	}
}

func (s *scrapingHandler) exec() error {
	// Instantiate spinner processes
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
	sp.Suffix = fmt.Sprintf(" [%s]: Processing...", "scraper")
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]: Completed!\n", green("✔"), "scraper")
	sp.Start()
	t, err := s.usecase.Scraping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	sp.Stop()

	s.printInfo(t)

	return nil
}

func (s *scrapingHandler) printInfo(t domain.Ticket) error {
	fmt.Printf(
		"============Ticket===========\nLink:	%s\nTitle:	%s\nStatus:	%s\nBug:	%d\n=============================\n%s\n",
		t.Url,
		t.Title,
		t.Status,
		t.BugCount,
		strings.Join(t.Bugs, "\n"),
	)
	return nil
}
