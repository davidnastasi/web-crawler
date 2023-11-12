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

func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", url, err)
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
