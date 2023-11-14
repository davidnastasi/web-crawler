package analyzer

import (
	"strconv"
	"strings"
)

type RobotAnalyzer struct{}

type RobotRule struct {
	UserAgent string
	Allow     []string
	Disallow  []string
	CrawDelay int
}

func GetRobotsContent(value string) ([]RobotRule, error) {
	var rbs []RobotRule
	var rb *RobotRule
	lines := strings.Split(value, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "User-agent:") {
			if rb != nil {
				rbs = append(rbs, *rb)
			}
			rb = &RobotRule{}
			rb.UserAgent = strings.TrimSpace(line[len("User-agent:"):])
		} else if strings.HasPrefix(line, "Disallow:") {
			rb.Disallow = append(rb.Disallow, strings.TrimSpace(line[len("Disallow:"):]))
		} else if strings.HasPrefix(line, "Allow:") {
			rb.Allow = append(rb.Allow, strings.TrimSpace(line[len("Allow:"):]))
		} else if strings.HasPrefix(line, "Craw-Delay:") {
			strValue := strings.TrimSpace(line[len("Craw-Delay:"):])
			intValue, err := strconv.Atoi(strValue)
			if err != nil {
				return []RobotRule{}, err
			}
			rb.CrawDelay = intValue
		}
	}

	if rb != nil {
		rbs = append(rbs, *rb)
	}

	return rbs, nil
}
