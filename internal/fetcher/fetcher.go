package fetcher

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	ErrNotFound        = errors.New("page not found")
	ErrTooManyRequest  = errors.New("too many request")
	ErrInvalidResponse = errors.New("invalid response")
)

func Fetch(url string, userAgent string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request %s: %v\n", url, err)
		return "", err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching request %s: %v\n", url, err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return string(body), nil
	case http.StatusNotFound:
		return "", ErrNotFound
	case http.StatusTooManyRequests:
		return "", ErrTooManyRequest
	default:
		return "", fmt.Errorf("status code: %d %w", resp.StatusCode, ErrInvalidResponse)
	}
}
