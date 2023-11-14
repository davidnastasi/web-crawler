package fetcher_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-crawler/internal/fetcher"
)

func TestFetch(t *testing.T) {
	// Mock server to simulate different HTTP responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "Success response")
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/toomanyrequests":
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()

	urlSuccess := mockServer.URL + "/success"
	urlNotFound := mockServer.URL + "/notfound"
	urlTooManyRequests := mockServer.URL + "/toomanyrequests"
	urlInternalServerError := mockServer.URL + "/internalerror"

	successResult, successErr := fetcher.Fetch(urlSuccess, "")
	_, notFoundErr := fetcher.Fetch(urlNotFound, "")
	_, tooManyRequestsErr := fetcher.Fetch(urlTooManyRequests, "")
	_, internalServerErrorErr := fetcher.Fetch(urlInternalServerError, "")

	if successResult != "Success response" || successErr != nil {
		t.Errorf("Unexpected result for success URL. Got %v, %v, want %v, nil", successResult, successErr, "Success response")
	}

	if !errors.Is(notFoundErr, fetcher.ErrNotFound) {
		t.Errorf("Unexpected error for not found URL. Got %v, want %v", notFoundErr, fetcher.ErrNotFound)
	}

	if !errors.Is(tooManyRequestsErr, fetcher.ErrTooManyRequest) {
		t.Errorf("Unexpected error for too many requests URL. Got %v, want %v", tooManyRequestsErr, fetcher.ErrTooManyRequest)
	}

	if internalServerErrorErr == nil || internalServerErrorErr.Error() != "status code: 500 invalid response" {
		t.Errorf("Unexpected error for internal server error URL. Got %v, want %v", internalServerErrorErr, "status code: 500 invalid response")
	}
}

func TestFetch_GetFailure(t *testing.T) {
	// An intentionally invalid or unreachable URL
	invalidURL := "http://invalid.url"

	_, err := fetcher.Fetch(invalidURL, "")

	if err == nil {
		t.Errorf("Expected error for GET request failure")
	}
}

func TestFetch_FailToCreateRequest(t *testing.T) {
	// An intentionally invalid or unreachable URL
	invalidURL := "http://inv alid.url"

	_, err := fetcher.Fetch(invalidURL, "")

	if err == nil {
		t.Errorf("Expected error for GET request failure")
	}
}
