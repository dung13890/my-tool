package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/my-tool/domain"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
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
			g := &googleSheetsHandler{}

			g.setUp()
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

	// handler usecase
	sp.Stop()
	return nil
}

func (g *googleSheetsHandler) setUp() error {
	// b, err := ioutil.ReadFile("infrastructure/credentials.json")

	return nil
}
