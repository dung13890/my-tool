package main

import (
	"github.com/dung13890/my-tool/cmd"
	"github.com/dung13890/my-tool/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	cli.AppHelpTemplate = config.AppHelpTemplate
	cli.CommandHelpTemplate = config.CommandHelpTemplate
	scraping := cmd.NewScraping()

	app := &cli.App{
		Name:                 "main",
		Usage:                "My Tools",
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			scraping,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
