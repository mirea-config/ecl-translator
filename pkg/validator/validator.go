package validator

import (
	"fmt"
	"regexp"
)

func IsNameValid(name string) bool {
	if name == "" {
		return false
	}

	matched := false
	var err error
	if len(name) == 1 {
		matched, err = regexp.MatchString(`[a-zA-Z]`, name)
	} else {
		matched, err = regexp.MatchString(`^[_a-zA-Z][_a-zA-Z0-9]*$`, name)
	}
	if err != nil {
		fmt.Printf("tlang: %s\n", err.Error())
	}

	return matched
}
