package analyzer

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

	severalContent := `
		User-agent: *
		Disallow: /private
		
		User-agent: googlebot
		Disallow: /restricted
	`

	validResult, validErr := GetRobotsContent(validContent)
	invalidResult, invalidErr := GetRobotsContent(invalidContent)
	severalResult, severalErr := GetRobotsContent(severalContent)

	expectedValidResult := []RobotRule{
		{
			UserAgent: "*",
			Disallow:  []string{"/private"},
			Allow:     []string{"/public"},
			CrawDelay: 500,
		},
	}

	var expectedInvalidResult []RobotRule

	expectedSeveralResult := []RobotRule{
		{
			UserAgent: "*",
			Disallow:  []string{"/private"},
		},
		{
			UserAgent: "googlebot",
			Disallow:  []string{"/restricted"},
		},
	}

	if !reflect.DeepEqual(validResult[0], expectedValidResult[0]) || validErr != nil {
		t.Errorf("Unexpected result for valid content. Got %v, %v, want %v, nil", validResult, validErr, expectedValidResult)
	}

	if len(invalidResult) != 0 || invalidErr == nil {
		t.Errorf("Unexpected result for invalid content. Got %v, %v, want %v, error", invalidResult, invalidErr, expectedInvalidResult)
	}

	if len(severalResult) != 2 || severalErr != nil {
		t.Errorf("Unexpected result for several content. Got %v, %v, want %v, error", severalResult, severalErr, expectedSeveralResult)
	}
}
