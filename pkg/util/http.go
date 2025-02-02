package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func PostJSON(url string, val any, timeout time.Duration) (*http.Response, error) {
	payload, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: timeout,
	}
	return client.Post(url, "application/json", bytes.NewBuffer(payload))
}
