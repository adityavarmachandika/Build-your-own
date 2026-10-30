package main

import (
	"fmt"

)


func parseBrackets(unknown string)(bool){

	closedBracketMap:=map[rune]rune{
		'}':'{',
		']':'[',
		')':'(',
	}
	openBracketMap:=map[rune]rune{
		'{':'}',
		'[':']',
		'(':')',
	}


	var stack []rune

	for _,char :=range unknown{
		if _,ok :=openBracketMap[char]; ok{
			stack = append(stack, char)
			continue
		}else if _,ok:=closedBracketMap[char]; ok{
			if len(stack)==0 {
				return false
			}
			n:=len(stack)-1
			if closedBracketMap[char]!=stack[n]{
				return false
			}
			stack= stack[:n]
		}
	}

	

	return len(stack)==0
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

func main(){
	fmt.Println("-----json parser-----")

	//unknownJSON := `{"name": "Alice", "age": 30, "skills": ["Go", "Python"], "active": true}`
	unknownbracketjson:=`{"message": "This is invalid structural JSON }["}`
	
	//if valid:=parseBrackets(unknownbracketjson); valid{
	if valid:=parseBracketsOptimized(unknownbracketjson); valid{
		fmt.Println("parsed succesfull")
		return
	}

	fmt.Println("unable to parse")

	
}