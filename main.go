package main

import (
	"github.com/dung13890/my-tool/config"
	"github.com/dung13890/my-tool/scraping/delivery/cmd"
	"github.com/dung13890/my-tool/scraping/delivery/http"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	viper.SetConfigFile(`infrastructure/config.json`)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Errors: not exists or is wrong json format", err)
	}
	cli.AppHelpTemplate = config.AppHelpTemplate
	cli.CommandHelpTemplate = config.CommandHelpTemplate
	scraping := cmd.NewScraping()
	serve := http.NewServe()

	app := &cli.App{
		Name:                 "main",
		Usage:                "My Tools",
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			scraping,
			serve,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
