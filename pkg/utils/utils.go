package utils

import "reflect"

func ContainsDifferentTypes(src []interface{}) bool {
	if len(src) < 2 {
		return false
	}

	baseType := reflect.TypeOf(src[0]).Kind()

	for _, val := range src[1:] {
		if reflect.TypeOf(val).Kind() != baseType {
			return true
		}
	}

	return false
}
