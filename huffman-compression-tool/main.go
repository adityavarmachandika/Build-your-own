package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("there is an error while fetching the file \n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if scanner.Err() != nil {
		fmt.Println("there is an error while using bufio newscanner", scanner)
	}
	var outText strings.Builder
	for scanner.Scan() {
		tempString := scanner.Text()
		outText.WriteString(tempString)
	}

	text := outText.String()
	wordMap := make(map[string]int)

	for _, char := range text {
		wordMap[string(char)]++
	}

	fmt.Println(wordMap)

}
