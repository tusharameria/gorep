package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	fileName := ""
	caseInsensetive := false
	flag.BoolVar(&caseInsensetive, "i", caseInsensetive, "True or False")
	flag.StringVar(&fileName, "f", fileName, "Need valid file address")
	flag.Parse()
	args := flag.Args()

	searchQuery := args[0]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 1
	for scanner.Scan() {
		text := scanner.Text()
		textTest := text
		if caseInsensetive {
			textTest = strings.ToLower(text)
			searchQuery = strings.ToLower(searchQuery)
		}
		if strings.Contains(textTest, searchQuery) {
			fmt.Printf("%d %s\n", i, text)
		}
		i++
	}
}
