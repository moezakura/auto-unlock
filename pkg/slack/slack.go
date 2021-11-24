package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/moezakura/auto-unlock/pkg/config"
	"golang.org/x/xerrors"
)

type Slack struct {
	WebhookURL string
}

func NewSlack(cfg *config.Config) *Slack {
	return &Slack{
		WebhookURL: cfg.WebhookURL,
	}
}

func (s *Slack) SendMessage(text string) error {
	d := map[string]string{
		"text": text,
	}
	sd, err := json.Marshal(d)
	if err != nil {
		return xerrors.Errorf("failed to encode message: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		s.WebhookURL,
		bytes.NewReader(sd),
	)
	if err != nil {
		return xerrors.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return xerrors.Errorf("failed to sendMessage http response: %w")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return err
}
