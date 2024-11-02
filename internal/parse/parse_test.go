package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSingleLineComment(t *testing.T) {
	comment := "Sample comment"
	result := parseSingleLineComment(comment)
	assert.Equal(t, "! Sample comment", result)
}

func TestParseMultilineComment(t *testing.T) {
	input := []interface{}{"Line 1", "Line 2"}
	result, err := parseMultilineComment(input)
	require.NoError(t, err)
	assert.Equal(t, "|#\nLine 1\nLine 2\n|#", result)
}

func TestParseArray(t *testing.T) {
	values := []interface{}{"value1", "value2", "3"}
	result, err := parseArray(values)
	require.NoError(t, err)
	assert.Equal(t, "[ @\"value1\", @\"value2\", @\"3\" ]", result)
}

func TestParseVar(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		value   interface{}
		output  string
		isErr   bool
	}{
		{"ValidInt", "a", 5, "a = 5", false},
		{"ValidString", "b", "hello", "b = @\"hello\"", false},
		{"ValidFloat", "c", 3.14, "c = 3.14", false},
		{"InvalidName", "1invalid", 5, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseVar(tt.varName, tt.value)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.output, result)
			}
		})
	}
}

func TestParseVal(t *testing.T) {
	assert.Equal(t, "@\"hello\"", parseVal("hello"))
	assert.Equal(t, "42", parseVal(42))
}

func TestParseConst(t *testing.T) {
	constants = make(map[string]string)
	result, err := parseConst("pi", 3.14)
	require.NoError(t, err)
	assert.Equal(t, "def pi = 3.14;", result)
	assert.Equal(t, "pi = 3.14", constants["pi"])
}

func TestEvalConst(t *testing.T) {
	constants = map[string]string{
		"pi":  "3.14",
		"arr": "[ 1, 2 ]",
	}

	result, err := evalConst("?(pi)")
	require.NoError(t, err)
	assert.Equal(t, "3.14", result)

	result, err = evalConst("?(arr)")
	require.NoError(t, err)
	assert.Equal(t, "[ 1, 2 ]", result)

	_, err = evalConst("?(unknown)")
	assert.Error(t, err)
}
