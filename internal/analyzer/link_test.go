package analyzer_test

import (
	"net/url"
	"testing"
	"web-crawler/internal/analyzer"
)

func TestGetLinks_Success(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	htmlContent := `
		<html>
			<body>
				<a href="/page1">Link 1</a>
				<a href="https://example.com/page2">Link 2</a>
				<a href="https://example.com/page3#section">Link 3 with Fragment</a>
			</body>
		</html>
	`

	links := analyzer.GetLinks(baseURL, htmlContent)
	w1, _ := url.Parse("https://example.com/page1")
	w2, _ := url.Parse("https://example.com/page2")
	w3, _ := url.Parse("https://example.com/page3#section")

	expectedLinks := []*url.URL{w1, w2, w3}

	// Check if the extracted links match the expected links
	if !equalSlices(links, expectedLinks) {
		t.Errorf("Extracted links are not as expected. Got %v, want %v", links, expectedLinks)
	}
}

func equalSlices(a, b []*url.URL) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].String() != b[i].String() {
			return false
		}
	}
	return true
}

func TestGetLinks_Error(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	// Malformed HTML content that will cause an error during tokenization
	htmlContent := `<html><body><a href=":malformedlink">Link</a></body></html>`

	links := analyzer.GetLinks(baseURL, htmlContent)

	// Ensure that the result is an empty slice due to the error in tokenization
	if len(links) != 0 {
		t.Errorf("Expected no links due to error in tokenization, got %v", links)
	}
}

func TestGetLinks_MailTo(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	// Malformed HTML content that will cause an error during tokenization
	htmlContent := `<html><body><a href="mailto:mail">Link</a></body></html>`

	links := analyzer.GetLinks(baseURL, htmlContent)

	// Ensure that the result is an empty slice due to the error in tokenization
	if len(links) != 0 {
		t.Errorf("Expected no links due to error in tokenization, got %v", links)
	}
}
