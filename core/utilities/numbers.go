package utilities

import "fmt"

func AddNumberSuffix(n int) string {
	suffix := "th"
	switch n {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return fmt.Sprintf("%d%s", n, suffix)
}
