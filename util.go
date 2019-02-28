package zenfo

import (
	"strings"
)

func clean(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}
