package helper

import "strings"

// TrimWhiteSpace :- will check whitespace
func TrimWhiteSpace(field string) bool {
	if strings.TrimSpace(field) == "" {
		return false
	}
	return true
}
