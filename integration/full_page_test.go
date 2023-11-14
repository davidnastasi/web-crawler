package integration_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"web-crawler/internal/crawler"
)

func TestFullPage(t *testing.T) {

	c := crawler.NewCrawler()
	err := c.Run("http://localhost:8989/", 5, 10)
	require.NoError(t, err)
	assert.Len(t, c.GetVisitedURLs(), 5)

}
