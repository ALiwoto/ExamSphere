package appValues

import "strings"

func NormalizeUserId(userId string) string {
	return strings.ToLower(strings.TrimSpace(userId))
}
