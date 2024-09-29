package cmdio

import "strings"

func fmtout(s string) string {
	if s == "" {
		return " <empty>\n"
	}
	s = strings.TrimRight(s, "\n")
	s = strings.ReplaceAll(s, "\n", "\n\t")
	return "\n\t" + s + "\n"
}
