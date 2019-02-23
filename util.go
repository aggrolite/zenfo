package zenfo

import (
	"strings"
)

func cleanWhiteSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
