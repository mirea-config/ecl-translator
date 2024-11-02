package validator_test

import (
	"ecl-translator/pkg/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNameValid(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Name    string
		IsValid bool
	}{
		{"a", true},
		{"ab", true},
		{"111", false},
		{"2", false},
		{"1name", false},
		{"_name", true},
	}

	for i, test := range tests {
		valid := validator.IsNameValid(test.Name)
		if !a.Equal(test.IsValid, valid) {
			t.Errorf("test %d failed at name '%s'", i+1, test.Name)
		}
	}
}
