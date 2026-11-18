package main

import (
	"encoding/json"
	"fmt"
)

func parseBrackets(unknown string) bool {

	closedBracketMap := map[rune]rune{
		'}': '{',
		']': '[',
		')': '(',
	}
	openBracketMap := map[rune]rune{
		'{': '}',
		'[': ']',
		'(': ')',
	}

	var stack []rune

	for _, char := range unknown {
		if _, ok := openBracketMap[char]; ok {
			stack = append(stack, char)
			continue
		} else if _, ok := closedBracketMap[char]; ok {
			if len(stack) == 0 {
				return false
			}
			n := len(stack) - 1
			if closedBracketMap[char] != stack[n] {
				return false
			}
			stack = stack[:n]
		}
	}

	return len(stack) == 0
}

// Optimized Switch implementation
func parseBracketsOptimized(unknown string) bool {
	var stack []rune

	for _, char := range unknown {
		switch char {
		case '{', '[', '(':
			stack = append(stack, char)
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func parseJson(unknown string) (map[string]any, error) {

	if valid := parseBracketsOptimized(unknown); !valid {
		fmt.Println("unable to parse")
		return nil, fmt.Errorf("Not a valid json")
	}
	index := 0
	prasedJson, err := recursion(unknown, &index)
	if err != nil {
		return nil, err
	}
	return prasedJson, nil
}

func recursion(unknown string, in *int) (map[string]any, error) {

	result := make(map[string]any)
	var key, value any
	var err error

	index := *in
	for ; index < len(unknown); index++ {

		removeWhiteSpaces(unknown, &index)
		char := unknown[index]
		switch char {

		case '{', ',':
			index++
			removeWhiteSpaces(unknown, &index)
			if index < len(unknown) && unknown[index] == '}' {
				index-- // Let the loop's index++ land on '}'
				continue
			}
			key, err = getString(unknown, &index)
			if err != nil {
				return nil, err
			}
			break
		case ':':
			index++
			removeWhiteSpaces(unknown, &index)
			char = unknown[index]
			if char == '{' {
				value, err = recursion(unknown, &index)
				if err != nil {
					return nil, err
				}
			} else if char == '"' {
				value, err = getString(unknown, &index)
				if err != nil {
					return nil, err
				}
			} else if char == '[' {
				value, err = parseArray(unknown, &index)
				if err != nil {
					return nil, err
				}
			}
			k, _ := key.(string)
			result[k] = value
			break
		case '}':
			*in = index
			return result, nil
		}

	}
	return result, nil
}

func parseArray(unknown string, in *int) ([]any, error) {
	result := []any{}
	index := *in + 1 // step past '['

	for index < len(unknown) {
		removeWhiteSpaces(unknown, &index)
		if index >= len(unknown) {
			return nil, fmt.Errorf("unexpected end of array")
		}
		if unknown[index] == ']' {
			*in = index
			return result, nil
		}

		var value any
		var err error
		switch unknown[index] {
		case '"':
			value, err = getString(unknown, &index)
		case '{':
			value, err = recursion(unknown, &index)
		case '[':
			value, err = parseArray(unknown, &index)
		default:
			return nil, fmt.Errorf("unsupported array value %c at %d", unknown[index], index)
		}
		if err != nil {
			return nil, err
		}
		result = append(result, value)

		index++
		removeWhiteSpaces(unknown, &index)
		if index >= len(unknown) {
			return nil, fmt.Errorf("unexpected end of array")
		}
		switch unknown[index] {
		case ',':
			index++
		case ']':
			*in = index
			return result, nil
		default:
			return nil, fmt.Errorf("expected , or ] but got %c", unknown[index])
		}
	}
	return nil, fmt.Errorf("unexpected end of array")
}

func getString(unknown string, i *int) (any, error) {

	var data string
	index := *i
	if unknown[index] != '"' {
		return "", fmt.Errorf("invalid json. missing \", currrent string is %c and index is %d", unknown[index], index)
	}
	index++
	for ; index < len(unknown); index++ {
		char := unknown[index]
		if char != '"' {
			data += string(char)
		} else {
			break
		}
	}
	*i = index
	return data, nil
}

func parseJsonOptimized(unknown string) {
	var data any
	if err := json.Unmarshal([]byte(unknown), &data); err != nil {
		//fmt.Println("there was an error while extracting")
		fmt.Println(err)
	}
	//fmt.Println(data)
}

func removeWhiteSpaces(unknown string, index *int) {
	for *index < len(unknown) && (unknown[*index] == ' ' || unknown[*index] == '\t' || unknown[*index] == '\n' || unknown[*index] == '\r') {
		*index++
	}
}

func main() {
	fmt.Println("-----json parser-----")

	unknownJSON := `{"name":"Alice",  "age":"30"  ,"skills":"Python","active":"true"}`
	//unknownbracketjson:=`{"message": "This is invalid structural JSON }"}`

	parseJsonOptimized(unknownJSON)

	res, err := parseJson(unknownJSON)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}
