package http

import (
	"fmt"
	"github.com/dung13890/my-tool/domain"
	"github.com/dung13890/my-tool/scraping/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type scrapingHandler struct {
	usecase domain.TicketUsecase
}

type params struct {
	WebhookSettingId string `json:"webhook_setting_id"`
	WebhookEventType string `json:"webhook_event_type"`
	WebhookEvent     struct {
		FromAccountId int    `json:"from_account_id"`
		ToAccountId   int    `json:"to_account_id"`
		RoomId        int    `json:"room_id"`
		MessageId     string `json:"message_id"`
		Body          string `json:"body"`
	} `json:"webhook_event"`
}

func NewServe() *cli.Command {
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"e"},
		Usage:   "start of serve",
		Action: func(ctx *cli.Context) error {
			e := echo.New()
			s := &scrapingHandler{}
			e.POST("/webhook", s.webhook)
			e.POST("/tickets/info", s.ticketInfo)
			e.Logger.Fatal(e.Start(":" + viper.GetString(`serve.port`)))
			return nil
		},
	}
}

func (s *scrapingHandler) webhook(c echo.Context) error {
	p := &params{}
	if err := c.Bind(p); err != nil {
		return err
	}
	ticket := regexp.MustCompile(`https://pherusa([-/\.\w\d])*|https://dev.sun-asterisk([-/\.\w\d])*`)
	url := ""
	switch {
	case ticket.MatchString(p.WebhookEvent.Body):
		url = ticket.FindString(p.WebhookEvent.Body)
		s.usecase = usecase.NewScrapingUsecase(url)
	default:
		return nil
	}
	if url == "" {
		return nil
	}

	t, err := s.usecase.Scraping()
	if err != nil {
		return err
	}

	s.reply(t, p)

	return nil
}

func (s *scrapingHandler) ticketInfo(c echo.Context) error {
	p := &struct {
		TicketUrl string `json:"ticket_url"`
	}{}
	if err := c.Bind(p); err != nil {
		return err
	}

	ticket := regexp.MustCompile(`https://pherusa([-/\.\w\d])*|https://dev.sun-asterisk([-/\.\w\d])*|https://framgiabrg.backlog([-/\.\w\d])*`)
	url := ""
	switch {
	case ticket.MatchString(p.TicketUrl):
		url = ticket.FindString(p.TicketUrl)
		s.usecase = usecase.NewScrapingUsecase(url)
	default:
		return nil
	}
	if url == "" {
		return nil
	}

	t, err := s.usecase.Scraping()
	if err != nil {
		return err
	}

	rs := &struct {
		Name          string `json:"name"`
		Title         string `json:"title"`
		TargetVersion string `json:"target_version"`
		Url           string `json:"url"`
		Status        string `json:"status"`
		EstimatedTime int    `json:"estimated_time"`
		ActualTime    int    `json:"actual_time"`
		Priority      string `json:"priority"`
	}{
		Name:          t.Title,
		Title:         t.Title,
		TargetVersion: t.Version,
		Url:           t.Url,
		Status:        t.Status,
		EstimatedTime: t.EstimatedTime,
		ActualTime:    t.ActualTime,
		Priority:      t.Priority,
	}
	return c.JSON(http.StatusOK, rs)
}

func (s *scrapingHandler) reply(t domain.Ticket, p *params) error {
	reply := fmt.Sprintf(
		"[rp aid=%d to=%d-%s]",
		p.WebhookEvent.FromAccountId,
		p.WebhookEvent.RoomId,
		p.WebhookEvent.MessageId,
	)

	body := fmt.Sprintf(
		"%s\n[info]Link:	%s\nTitle:	%s\nStatus:	%s\nBug:	%d\n[hr]%s[/info]",
		reply,
		t.Url,
		t.Title,
		t.Status,
		t.BugCount,
		strings.Join(t.Bugs, "\n"),
	)

	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%d/messages?body=%s", p.WebhookEvent.RoomId, url.QueryEscape(body))
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-ChatWorkToken", viper.GetString(`chatwork.token`))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	res, _ := client.Do(req)
	defer res.Body.Close()

	return nil
}
