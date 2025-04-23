package shared

import (
	"regexp"
	"strings"
)

func Slugify(str string) string {
	str = strings.ToLower(str)

	re := regexp.MustCompile(`[^a-z0-9]+`)
	str = re.ReplaceAllString(str, "-")

	return strings.Trim(str, "-")
}
