package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/urfave/cli/v2"
	"time"
)

type scraping struct {
	url string
}

func NewScraping() *cli.Command {
	return &cli.Command{
		Name:    "scraping",
		Aliases: []string{"s"},
		Usage:   "scraping site",
		Action: func(ctx *cli.Context) error {
			s := &scraping{}
			if ctx.NArg() > 0 {
				s.url = ctx.Args().Get(0)
			} else {
				s.url = "https://google.com"
			}
			s.exec()
			return nil
		},
	}
}

func (s *scraping) exec() error {
	// Instantiate default collector
	c := colly.NewCollector()

	fmt.Printf("Scraping data from site: [%s]\n", s.url)
	// Instantiate spinner processes
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
	sp.Suffix = fmt.Sprintf(" [%s]: Processing...", "scraper")
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]: Completed!\n", green("✔"), "scraper")
	sp.Start()

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on url
	c.Visit(s.url)

	sp.Stop()

	return nil
}
