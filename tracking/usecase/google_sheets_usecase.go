package usecase

import (
	"fmt"
	"github.com/dung13890/my-tool/domain"
	"github.com/dung13890/my-tool/tracking/repository"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

type googleSheetsUsecase struct {
	repo domain.GoogleSheetsRepository
}

func NewGoogleSheetsUsecase() (domain.GoogleSheetsUsecase, error) {
	r, err := repository.NewGoogleSheetsRepository()
	if err != nil {
		return nil, err
	}
	return &googleSheetsUsecase{
		repo: r,
	}, nil
}

func (g *googleSheetsUsecase) Tracking() ([]*domain.GoogleSheets, error) {
	return g.repo.Fetch()
}

func (g *googleSheetsUsecase) Update(item *domain.GoogleSheets, bug int) error {
	return g.repo.Update(item, bug)
}

func (g *googleSheetsUsecase) Notify(items []domain.GoogleSheets) error {
	var content string
	for _, item := range items {
		content += fmt.Sprintf("Bug: (%d) %s\n", item.CurrentBug, item.Ticket)
	}

	body := fmt.Sprintf(
		"%s\n[info][title]Cảnh báo nhiều bug (devil)[/title]%s[/info]",
		viper.GetString(`chatwork.leader`),
		content,
	)

	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%s/messages?body=%s", viper.GetString(`chatwork.room`), url.QueryEscape(body))
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-ChatWorkToken", viper.GetString(`chatwork.token`))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	res, err := client.Do(req)
	defer res.Body.Close()

	return nil
}
