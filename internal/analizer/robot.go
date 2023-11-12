package analizer

import (
	"strconv"
	"strings"
)

type RobotAnalyzer struct{}

type RobotContent struct {
	UserAgent string
	Allow     []string
	Disallow  []string
	CrawDelay *int
}

func GetRobotContent(value string) (RobotContent, error) {
	rb := RobotContent{}
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "User-agent:") {
			rb.UserAgent = strings.TrimSpace(line[len("User-agent:"):])
		} else if strings.HasPrefix(line, "Disallow:") {
			rb.Disallow = append(rb.Disallow, strings.TrimSpace(line[len("Disallow:"):]))
		} else if strings.HasPrefix(line, "Allow:") {
			rb.Disallow = append(rb.Disallow, strings.TrimSpace(line[len("Allow:"):]))
		} else if strings.HasPrefix(line, "Craw-Delay:") {
			strValue := strings.TrimSpace(line[len("Craw-Delay:"):])
			intValue, err := strconv.Atoi(strValue)
			if err != nil {
				return RobotContent{}, err
			}
			rb.CrawDelay = &intValue
		}

	}
	return rb, nil
}
