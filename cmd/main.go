package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	caseInsensetive := false
	flag.BoolVar(&caseInsensetive, "i", caseInsensetive, "True or False")
	flag.Parse()
	args := flag.Args()

	searchQuery := args[0]
	fileNames := args[1:]

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		n := 1
		for scanner.Scan() {
			text := scanner.Text()
			textTest := text
			if caseInsensetive {
				textTest = strings.ToLower(text)
				searchQuery = strings.ToLower(searchQuery)
			}
			if strings.Contains(textTest, searchQuery) {
				fmt.Printf("%s:%d %s\n", fileName, n, text)
			}
			n++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if err := file.Close(); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}
}
