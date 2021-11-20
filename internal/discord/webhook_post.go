package discord

import (
	"bytes"
	"encoding/json"
	"github.com/roessland/csgostate/internal/logger"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Log        logger.Logger
	WebhookURL string
}

func NewClient(webhookURL string, log logger.Logger) *Client {
	c := &Client{}
	c.WebhookURL = webhookURL
	c.Log = log
	return c
}

type webhookSimpleMsg struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// TODO return error
func (c *Client) Post(msg string) {
	if c.WebhookURL == "" {
		c.Log.Infow("missing discord webhook url; logging instead", "msg", msg)
		return
	}
	body, _ := json.Marshal(webhookSimpleMsg{
		Username: "CS:GO State",
		Content:  msg,
	})
	req, err := http.NewRequest(http.MethodPost, c.WebhookURL, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	c.Log.Infow("posted to discord webhook",
		"msg", msg,
		"response_code", resp.StatusCode,
		"response_body", string(respBody))
}
