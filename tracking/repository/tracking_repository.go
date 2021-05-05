package repository

import (
	"errors"
	"fmt"
	"github.com/dung13890/my-tool/domain"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"strconv"
)

type googleSheetsRepository struct {
	srv *sheets.Service
}

func NewGoogleSheetsRepository() (domain.GoogleSheetsRepository, error) {
	b, err := ioutil.ReadFile("infrastructure/credentials.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: Unable to read client secret file: %v", err))
	}
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: Unable to parse client secret file to config: %v", err))
	}
	client := config.Client(oauth2.NoContext)
	srv, err := sheets.New(client)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: Unable to retrieve Sheets client: %v", err))
	}
	g := &googleSheetsRepository{
		srv: srv,
	}

	return g, nil
}

func (g *googleSheetsRepository) Fetch() ([]*domain.GoogleSheets, error) {
	var items []*domain.GoogleSheets
	spreadsheetId := viper.GetString(`spreadsheet.id`)
	readRange := viper.GetString(`spreadsheet.range`)

	resp, err := g.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors: Unable to retrieve data from sheet: %v", err))
	}

	if len(resp.Values) != 0 {
		for i, row := range resp.Values {
			if len(row) != 0 && row[3].(string) == "Testing" {
				item := &domain.GoogleSheets{
					Row:        i,
					Ticket:     row[0].(string),
					Status:     row[3].(string),
					CurrentBug: 0,
				}
				if len(row) > 4 {
					item.CurrentBug, _ = strconv.Atoi(row[4].(string))
				}
				items = append(items, item)
			}
		}
	}

	return items, nil
}

func (g *googleSheetsRepository) Update(item *domain.GoogleSheets, bug int) error {
	index := item.Row + 3
	spreadsheetId := viper.GetString(`spreadsheet.id`)
	writeRange := fmt.Sprintf("%s%d", viper.GetString(`spreadsheet.bug_range`), index)
	var rb sheets.ValueRange
	rb.Values = append(rb.Values, []interface{}{bug})

	_, err := g.srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &rb).ValueInputOption("RAW").Do()

	return err
}
