package aidetect

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const URL = "https://api.sapling.ai/api/v1/aidetect"

type aiDetect struct {
	URL    string
	Key    string
	client *http.Client
}

func (s aiDetect) Get(text string, scores bool) (Response, error) {
	var (
		responseParams Response
		params         = RequestParams{
			Key:        s.Key,
			Text:       text,
			SentScores: scores,
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

func NewClientV1(op Options) AIDetect {
	return aiDetect{
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

type AIDetect interface {
	Get(string, bool) (Response, error)
}

type SentenceScore struct {
	Sentence string  `json:"sentence"`
	Score    float64 `json:"score"`
}

type Response struct {
	Score          float64         `json:"score"`
	SentenceScores []SentenceScore `json:"sentence_scores"`
	Truncated      bool            `json:"truncated"`
	UsedTokens     int             `json:"used_tokens"`
	Hash           string          `json:"hash"`
}

type RequestParams struct {
	Key        string `json:"key"`
	Text       string `json:"text"`
	SentScores bool   `json:"sent_scores"`
}
