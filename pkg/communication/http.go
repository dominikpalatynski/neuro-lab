package communication

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Body       json.RawMessage `json:"body"`
	StatusCode int             `json:"statusCode"`
}

func SendRequest(method string, url string, body []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	return &Response{
		Body:       json.RawMessage(respBody),
		StatusCode: resp.StatusCode,
	}, err
}
