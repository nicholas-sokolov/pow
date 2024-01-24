package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Quote struct {
	Text   string `json:"content"`
	Author string
}

func GetQuote() (*Quote, error) {
	url := "https://api.quotable.io/quotes/random"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error to create request, %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, %v", err)
	}

	var qs []Quote
	err = json.Unmarshal(body, &qs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data from response, %v", err)
	}
	return &qs[0], nil
}
