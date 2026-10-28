package main

import (
	"fmt"

)



func parseBrackets(unknown string)(bool){

	bracketMap:=map[rune]rune{
		'}':'{',
		']':'[',
		')':'(',
	}


	var stack []rune

	for _,char :=range unknown{
		if _,ok :=bracketMap[char]; !ok{
			stack = append(stack, char)
			continue
		}
		if len(stack)==0{
			return false
		}
		n:=len(stack)-1
		if bracketMap[char]!=stack[n]{
			return false
		}
		stack= stack[:n]
		for _,char :=range stack{
			fmt.Printf("%c",char)
		}
	}

	

	return len(stack)==0
}



func main(){
	fmt.Println("-----json parser-----")

	//unknownJSON := `{"name": "Alice", "age": 30, "skills": ["Go", "Python"], "active": true}`
	unknownbracketjson:=`([])`

	if valid:=parseBrackets(unknownbracketjson); valid{
		fmt.Println("parsed succesfull")
		return
	}

	fmt.Println("unable to parse")

	
}