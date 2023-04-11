package statistics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

const URL = "https://api.sapling.ai/api/v1/statistics"

type statistics struct {
	URL    string
	Key    string
	client *http.Client
}

func (s statistics) Get(text string) (Response, error) {
	var (
		responseParams Response
		params         = RequestParams{
			Key:       s.Key,
			Text:      text,
			SessionID: uuid.NewV4().String(),
		}
	)

	paramsMarshalled, err := json.Marshal(params)
	if err != nil {
		return responseParams, err
	}

	req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader(paramsMarshalled))
	if err != nil {
		return responseParams, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return responseParams, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseParams); err != nil {
		return responseParams, err
	}

	return responseParams, nil
}

func NewClientV1(op Options) Statistics {
	return statistics{
		URL: URL,
		Key: op.Key,
		client: &http.Client{
			Timeout: op.Timeout,
		},
	}
}

type Options struct {
	Key     string
	Timeout time.Duration
}

type Statistics interface {
	Get(string) (Response, error)
}

type Response struct {
	Chars          int     `json:"chars`
	Readability    float64 `json:"readability"`
	ReadingTimeMin float64 `json:"reading_time_min"`
	ReadingTimeSec float64 `json:"reading_time_sec"`
	Sentences      int     `json:"sentences"`
	Words          int     `json:"words"`
}

type RequestParams struct {
	Key       string `json:"key"`
	Text      string `json:"text"`
	SessionID string `json:"session_id"`
}
