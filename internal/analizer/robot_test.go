package analizer

import (
	"reflect"
	"testing"
)

func TestGetRobotContent(t *testing.T) {
	// Valid robots.txt content
	validContent := `
		User-agent: *
		Disallow: /private
		Allow: /public
		Craw-Delay: 500
	`

	// Invalid robots.txt content with a Craw-Delay value that cannot be converted to an integer
	invalidContent := `
		User-agent: *
		Disallow: /private
		Craw-Delay: invalid
	`

	validResult, validErr := GetRobotContent(validContent)
	invalidResult, invalidErr := GetRobotContent(invalidContent)

	crawlDelay := 500
	expectedValidResult := RobotContent{
		UserAgent: "*",
		Disallow:  []string{"/private", "/public"},
		CrawDelay: &crawlDelay,
	}

	expectedInvalidResult := RobotContent{}

	if !reflect.DeepEqual(validResult, expectedValidResult) || validErr != nil {
		t.Errorf("Unexpected result for valid content. Got %v, %v, want %v, nil", validResult, validErr, expectedValidResult)
	}

	if !reflect.DeepEqual(invalidResult, expectedInvalidResult) || invalidErr == nil {
		t.Errorf("Unexpected result for invalid content. Got %v, %v, want %v, error", invalidResult, invalidErr, expectedInvalidResult)
	}
}
