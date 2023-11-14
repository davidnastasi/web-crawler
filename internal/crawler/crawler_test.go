package crawler_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-crawler/internal/crawler"
)

func TestCrawlSuccess(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/successcase":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/successcase/child1\"/></body></html>")
		case "/successcase/child1":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/successcase/child2\"/></body></html>")
		case "/successcase/child2":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/successcase/child1\"/></body></html>")
		}
	}))
	defer mockServer.Close()

	c := crawler.NewCrawler()
	err := c.Run(mockServer.URL+"/successcase", 4, 0)
	require.NoError(t, err)
	require.Len(t, c.GetVisitedURLs(), 3)

	c = crawler.NewCrawler()
	err = c.Run(mockServer.URL+"/successcase", 2, 0)
	require.NoError(t, err)
	require.Len(t, c.GetVisitedURLs(), 2)
}

func TestCrawlChildNotFound(t *testing.T) {
	c := crawler.NewCrawler()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/notfoundcase":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/notfoundcase/child1\"/></body></html>")
		case "/notfoundcase/child1":
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()
	err := c.Run(mockServer.URL+"/notfoundcase", 3, 0)
	require.NoError(t, err)
	require.Len(t, c.GetVisitedURLs(), 1)
}

func TestCrawlFragment(t *testing.T) {
	c := crawler.NewCrawler()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fragmentcase":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/fragmentcase/child1\"/></body></html>")
		case "/fragmentcase/child1":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/fragmentcase/child1#some\"/></body></html>")
		}
	}))
	defer mockServer.Close()
	err := c.Run(mockServer.URL+"/fragmentcase", 3, 0)
	require.NoError(t, err)
	require.Len(t, c.GetVisitedURLs(), 2)
}

func TestCrawlRobot(t *testing.T) {
	c := crawler.NewCrawler()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/robot.txt":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "User-agent: *\nDisallow: /robotcase/child2/\n")
		case "/robotcase":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/robotcase/child1\"/></body></html>")
		case "/robotcase/child1":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "<html><body><a href=\"/robotcase/child2\"/></body></html>")
		case "/robotcase/child2":
			w.WriteHeader(http.StatusOK)

		}
	}))
	defer mockServer.Close()
	err := c.Run(mockServer.URL+"/robotcase", 3, 0)
	require.NoError(t, err)
	require.Len(t, c.GetVisitedURLs(), 2)
}

func TestCrawlRobotTxtFails(t *testing.T) {
	c := crawler.NewCrawler()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robot.txt" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()
	err := c.Run(mockServer.URL+"/robotcase", 3, 0)
	require.Error(t, err)

}

func TestTooManyRequestError(t *testing.T) {
	c := crawler.NewCrawler()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/toomanyrequestcase":
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, "<html><body><a href=\"/toomanyrequestcase/child1\"/></body></html>")
		}
	}))
	defer mockServer.Close()
	err := c.Run(mockServer.URL+"/toomanyrequestcase", 3, 0)
	require.NoError(t, err)
	require.Empty(t, c.GetVisitedURLs())
}

func TestCrawlInvalidURL(t *testing.T) {
	c := crawler.NewCrawler()
	err := c.Run(":invalid_url:", 3, 0)
	require.Error(t, err)
}
