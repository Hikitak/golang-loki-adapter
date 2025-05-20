package loki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang-loki-adapter.local/pkg/models"
)

type LokiClient struct {
	cfg    *models.LokiConfig
	client *http.Client
}

func NewLokiClient(cfg *models.LokiConfig) *LokiClient {
	return &LokiClient{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func (c *LokiClient) Send(records []models.QueueRecord) error {
	payload := struct {
		Streams []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"streams"`
	}{
		Streams: []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		}{
			{
				Stream: c.cfg.Labels,
				Values: make([][]string, 0, len(records)),
			},
		},
	}

	for _, r := range records {
		payload.Streams[0].Values = append(
			payload.Streams[0].Values,
			[]string{
				fmt.Sprintf("%d", time.Now().UnixNano()),
				r.Data,
			},
		)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	for i := 0; i <= c.cfg.Retries; i++ {
		resp, err := c.client.Post(
			c.cfg.URL,
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			if i == c.cfg.Retries {
				return fmt.Errorf("http post failed: %w", err)
			}
			time.Sleep(time.Duration(1<<uint(i)) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			body, _ := ioutil.ReadAll(resp.Body)
			if i == c.cfg.Retries {
				return fmt.Errorf("loki responded with %d: %s",
					resp.StatusCode, string(body))
			}
			continue
		}
		return nil
	}
	return nil
}
