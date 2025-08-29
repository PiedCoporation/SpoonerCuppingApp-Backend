package stringutils

import (
	"strings"
	"unicode"
)

// pascal case / camel case -> snake case
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var b strings.Builder
	runes := []rune(s)
	n := len(runes)

	for i := range n {
		b.WriteRune(unicode.ToLower(runes[i]))

		nextIsUpper := i+1 < n && unicode.IsUpper(runes[i+1])
		overIsLowerOrNil := (i+2 >= n && unicode.IsLower(runes[i])) ||
			(i+2 < n && unicode.IsLower(runes[i+2]))
		if nextIsUpper && overIsLowerOrNil {
			b.WriteRune('_')
		}

	}

	return b.String()
}
