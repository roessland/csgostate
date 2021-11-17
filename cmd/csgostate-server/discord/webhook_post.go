package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	WebhookURL string
}

func NewClient(webhookURL string) *Client {
	c := &Client{}
	c.WebhookURL = webhookURL
	return c
}

type webhookSimpleMsg struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *Client) Post(msg string) {
	if c.WebhookURL == "" {
		fmt.Println("DISCORD POST:", msg)
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
	httpClient := &http.Client{Timeout: 10*time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}
