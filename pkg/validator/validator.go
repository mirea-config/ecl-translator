package validator

import (
	"fmt"
	"regexp"
)

func IsNameValid(name string) bool {
	if name == "" {
		return false
	}

	matched, err := regexp.MatchString(`[_a-zA-Z0-9].*`, name)
	if err != nil {
		fmt.Printf("tlang: %s\n", err.Error())
	}

	return matched
}
