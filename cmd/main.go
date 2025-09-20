package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	caseInsensetive := false
	isRecursive := false
	totalLines := false
	totalFiles := false
	flag.BoolVar(&caseInsensetive, "i", caseInsensetive, "True or False")
	flag.BoolVar(&isRecursive, "r", isRecursive, "True or False")
	flag.BoolVar(&totalLines, "tl", totalLines, "True or False")
	flag.BoolVar(&totalFiles, "tf", totalFiles, "True or False")
	flag.Parse()
	args := flag.Args()

	searchQuery := args[0]
	pathNames := args[1:]

	totalNoLines := 0
	totalNoFiles := 0

	var recursive func(string, bool)
	recursive = func(pathName string, firstLayer bool) {
		isFile, err := IsFile(pathName)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if isFile {
			file, err := os.Open(pathName)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			n := 1
			found := false
			for scanner.Scan() {
				text := scanner.Text()
				textTest := text
				if caseInsensetive {
					textTest = strings.ToLower(text)
					searchQuery = strings.ToLower(searchQuery)
				}
				if strings.Contains(textTest, searchQuery) {
					fmt.Printf("%s:%d %s\n", pathName, n, text)
					if totalLines {
						totalNoLines++
					}
					found = true
				}
				n++
			}
			if found && totalFiles {
				totalNoFiles++
			}

			if err := scanner.Err(); err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			if err := file.Close(); err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
		} else {
			if firstLayer || isRecursive {
				isDir, err := IsDir(pathName)
				if err != nil {
					fmt.Printf("%s\n", err)
					os.Exit(1)
				}

				if isDir {
					newPaths, err := os.ReadDir(pathName)
					if err != nil {
						fmt.Printf("%s\n", err)
						os.Exit(1)
					}
					for _, newPath := range newPaths {
						recursive(filepath.Join(pathName, newPath.Name()), false)
					}
				}
			}
		}
	}

	for _, pathName := range pathNames {
		recursive(pathName, true)
	}

	if totalLines {
		fmt.Printf("Total number of lines matched : %d\n", totalNoLines)
	}
	if totalFiles {
		fmt.Printf("Total number of files matched : %d\n", totalNoFiles)
	}
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}
