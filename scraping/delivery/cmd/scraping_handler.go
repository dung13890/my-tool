package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/my-tool/entities"
	"github.com/dung13890/my-tool/scraping"
	"github.com/dung13890/my-tool/scraping/usecase"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"regexp"
	"strings"
	"time"
)

type scrapingHandler struct {
	usecase scraping.Usecase
	url     string
}

func NewScraping() *cli.Command {
	return &cli.Command{
		Name:    "scraping",
		Aliases: []string{"s"},
		Usage:   "scraping site",
		Action: func(ctx *cli.Context) error {
			s := &scrapingHandler{}
			if ctx.NArg() > 0 {
				s.url = ctx.Args().Get(0)
				s.setUp()
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
		return err
	}
	sp.Stop()

	s.printInfo(t)

	return nil
}

func (s *scrapingHandler) setUp() error {
	pherusa := regexp.MustCompile(`https://pherusa([-/\.\w\d])*`)
	redmine := regexp.MustCompile(`https://dev.sun-asterisk([-/\.\w\d])*`)
	switch {
	case pherusa.MatchString(s.url):
		s.usecase = usecase.NewPherusaUsecase(s.url)
	case redmine.MatchString(s.url):
		s.usecase = usecase.NewRedmineUsecase(s.url)
	default:
		return nil
	}

	return nil
}

func (s *scrapingHandler) printInfo(t entities.Ticket) error {
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
