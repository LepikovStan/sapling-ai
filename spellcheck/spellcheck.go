package spellcheck

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

const URL = "https://api.sapling.ai/api/v1/spellcheck"

type spellcheck struct {
	URL    string
	Key    string
	client *http.Client
}

func (s spellcheck) Get(text string) (Response, error) {
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

func NewClientV1(op Options) SpellCheck {
	return spellcheck{
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

type SpellCheck interface {
	Get(string) (Response, error)
}

type Edit struct {
	SentenceStart int    `json:"sentence_start"`
	Start         int    `json:"start"`
	End           int    `json:"end"`
	ID            string `json:"id"`
	Replacement   string `json:"replacement"`
	Sentence      string `json:"sentence"`
}

type Response struct {
	Edits []Edit `json:"edits"`
}

type RequestParams struct {
	Key       string `json:"key"`
	Text      string `json:"text"`
	SessionID string `json:"session_id"`

	MinLength     int    `json:"min_length"`
	MultipleEdits bool   `json:"multiple_edits"`
	Lang          string `json:"lang"`
}
