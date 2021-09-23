package utils

import "strings"

func SiteTitle(args ...string) string {
	return strings.Join(args, " - ")
}
