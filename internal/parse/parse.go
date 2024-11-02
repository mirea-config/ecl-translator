package parse

import (
	"ecl-translator/internal/models"
	"ecl-translator/pkg/utils"
	"ecl-translator/pkg/validator"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var constants = make(map[string]string)

// В случае успешного парсинга возвращает слайс готовых строк исходного кода учебного конфигурационного языка.
// В противном случае возращает ошибку
func ParseJsonInput(input []byte) ([]string, error) {
	var jsonInput models.JsonInput

	if err := json.Unmarshal(input, &jsonInput); err != nil {
		return []string{}, err
	}

	lines := make([]string, 0, len(jsonInput.Tokens))
	for key, val := range jsonInput.Tokens {
		keyType := strings.Split(key, "_")[0]
		switch keyType {
		case "scomm":
			valString, ok := val.(string)
			if !ok {
				return []string{}, fmt.Errorf("single line comment must be a string")
			}
			lines = append(lines, parseSingleLineComment(valString))

		case "mcomm":
			commLines, ok := val.([]interface{})

			if !ok {
				return []string{}, fmt.Errorf("multi line comment must be an array of strings")
			}
			multiline, err := parseMultilineComment(commLines)
			if err != nil {
				return []string{}, err
			}
			lines = append(lines, multiline)

		case "var":
			varInfoMap, ok := val.(map[string]interface{})
			if !ok {
				return []string{}, fmt.Errorf("\"var\" key must map to json object with \"name\" and \"value\" fields, but none provided")
			}

			varName, ok := varInfoMap["name"]
			if !ok {
				return []string{}, fmt.Errorf("object has no \"name\" field")
			}
			varValue, ok := varInfoMap["value"]
			if !ok {
				return []string{}, fmt.Errorf("object has no \"value\" field")
			}

			varNameStr, ok := varName.(string)
			if !ok {
				return []string{}, fmt.Errorf("\"name\" must be a string")
			}

			variable, err := parseVar(varNameStr, varValue)
			if err != nil {
				return []string{}, err
			}
			lines = append(lines, variable)

		case "const":
			constInfoMap, ok := val.(map[string]interface{})
			if !ok {
				return []string{}, fmt.Errorf("\"const\" key must map to json object with \"name\" and \"value\" fields, but none provided")
			}

			constName, ok := constInfoMap["name"]
			if !ok {
				return []string{}, fmt.Errorf("object has no \"name\" field")
			}
			constValue, ok := constInfoMap["value"]
			if !ok {
				return []string{}, fmt.Errorf("object has no \"value\" field")
			}

			constNameStr, ok := constName.(string)
			if !ok {
				return []string{}, fmt.Errorf("\"name\" must be a string")
			}

			constant, err := parseConst(constNameStr, constValue)
			if err != nil {
				return []string{}, err
			}
			lines = append(lines, constant)

		default:
			return []string{}, fmt.Errorf("\"%s\" key is unknown", key)
		}
	}

	return lines, nil
}

func parseSingleLineComment(comment string) string {
	return fmt.Sprintf("! %s", comment)
}

func parseMultilineComment(commLines []interface{}) (string, error) {
	if utils.ContainsDifferentTypes(commLines) {
		return "", fmt.Errorf("failed to parse multiline comment: array contains different types")
	}

	if _, ok := commLines[0].(string); !ok {
		return "", fmt.Errorf("failed to parse multiline comment: array must contain strings")
	}

	strCommLines := make([]string, len(commLines))
	for i, line := range commLines {
		strCommLines[i] = line.(string)
	}

	return fmt.Sprintf("|#\n%s\n|#", strings.Join(strCommLines, "\n")), nil
}

func parseArray(values []interface{}) (string, error) {
	if utils.ContainsDifferentTypes(values) {
		return "", fmt.Errorf("failed to parse array: array contains different value types")
	}

	strValues := make([]string, len(values))
	for i, val := range values {
		strValues[i] = parseVal(val)
	}

	return fmt.Sprintf("[ %s ]", strings.Join(strValues, ", ")), nil
}

func parseVar(name string, value interface{}) (string, error) {
	if !validator.IsNameValid(name) {
		return "", fmt.Errorf("'%s' is invalid variable name", name)
	}

	valType := reflect.TypeOf(value).Kind()

	switch valType {
	case reflect.Int, reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%s = %v", name, value), nil

	case reflect.String:
		strValue := value.(string)
		matched, err := regexp.MatchString(`\?\(.*\)`, strValue)
		if err != nil {
			return "", err
		}
		if matched {
			return evalConst(strValue)
		}
		return fmt.Sprintf("%s = @\"%s\"", name, strValue), nil

	case reflect.Slice:
		valSlice := value.([]interface{})

		arr, err := parseArray(valSlice)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", name, arr), nil

	default:
		return "", fmt.Errorf("unexpected value type: %s", valType.String())
	}
}

func parseVal(value interface{}) string {
	valType := reflect.TypeOf(value).Kind()
	switch valType {
	case reflect.String:
		return fmt.Sprintf("@\"%v\"", value)

	default:
		return fmt.Sprintf("%v", value)
	}
}

func parseConst(name string, value interface{}) (string, error) {
	valExpr, err := parseVar(name, value)
	if err != nil {
		return "", err
	}

	constants[name] = valExpr

	return fmt.Sprintf("def %s;", valExpr), nil
}

func evalConst(token string) (string, error) {
	name, found := strings.CutPrefix(token, "?(")
	if !found {
		return "", fmt.Errorf("bad const evaluation syntax")
	}
	name, found = strings.CutSuffix(name, ")")
	if !found {
		return "", fmt.Errorf("bad const evaluation syntax")
	}

	val, ok := constants[name]
	if !ok {
		return "", fmt.Errorf("constant name '%s' is unknown", name)
	}

	return val, nil
}
