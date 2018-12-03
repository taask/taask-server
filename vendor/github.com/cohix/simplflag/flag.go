package simplflag

import (
	"os"
	"strings"
)

// GetFlag finds a named flag in the args list and returns its value
func GetFlag(name string) string {
	args := os.Args

	for _, arg := range args {
		if strings.HasPrefix(arg, "-"+name) || strings.HasPrefix(arg, "--"+name) {
			i := strings.Index(arg, "=")
			if i > -1 {
				return arg[i+1:]
			}
		}
	}
	return ""
}

// CheckFlag determines if a flag exists, even if it doesn't have a value
func CheckFlag(name string) (string, bool) {
	args := os.Args

	for _, arg := range args {
		if strings.HasPrefix(arg, "-"+name) || strings.HasPrefix(arg, "--"+name) {
			i := strings.Index(arg, "=")
			if i > -1 {
				return arg[i+1:], true
			}

			return "", true
		}
	}
	return "", false
}

// ArgsNoFlags returns the args list provided, minus any - or -- flags
func ArgsNoFlags(args []string) []string {
	ret := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") == false && strings.HasPrefix(arg, "--") == false {
			ret = append(ret, arg)
		}
	}

	return ret
}
