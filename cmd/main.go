package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	caseInsensetive := false
	isRecursive := false
	totalLines := false
	totalFiles := false
	showTime := false
	flag.BoolVar(&caseInsensetive, "i", caseInsensetive, "True or False")
	flag.BoolVar(&isRecursive, "r", isRecursive, "True or False")
	flag.BoolVar(&showTime, "t", showTime, "True or False")
	flag.BoolVar(&totalLines, "tl", totalLines, "True or False")
	flag.BoolVar(&totalFiles, "tf", totalFiles, "True or False")
	flag.Parse()
	args := flag.Args()

	searchQuery := args[0]
	pathNames := args[1:]

	lineSlices := make(chan []string)

	totalNoLines := 0
	totalNoFiles := 0

	wg := sync.WaitGroup{}

	var recursive func(string, bool)
	recursive = func(pathName string, firstLayer bool) {
		isFile, err := IsFile(pathName)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if isFile {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := ProcessFile(pathName, searchQuery, caseInsensetive, totalLines, totalFiles, &totalNoLines, &totalNoFiles, lineSlices); err != nil {
					fmt.Printf("%s\n", err)
					os.Exit(1)
				}
			}()
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

	go func() {
		for _, pathName := range pathNames {
			recursive(pathName, true)
		}
		wg.Wait()
		close(lineSlices)
	}()

	for lines := range lineSlices {
		for _, line := range lines {
			fmt.Printf("%s", line)
		}
	}

	if totalLines {
		fmt.Printf("Total number of lines matched : %d\n", totalNoLines)
	}

	if totalFiles {
		fmt.Printf("Total number of files matched : %d\n", totalNoFiles)
	}

	if showTime {
		timeinMs := time.Since(startTime).Milliseconds()
		if timeinMs < 1000 {
			fmt.Printf("Time taken : %dms\n", timeinMs)
		} else {
			fmt.Printf("Time taken : %fs\n", float64(timeinMs)/1000)
		}
	}
}

func ProcessFile(pathName, searchQuery string, caseInsensetive, totalLines, totalFiles bool, totalNoLines, totalNoFiles *int, lineSlices chan<- []string) error {
	file, err := os.Open(pathName)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	n := 1
	found := false
	if caseInsensetive {
		searchQuery = strings.ToLower(searchQuery)
	}

	fileLines := []string{}
	for scanner.Scan() {
		text := scanner.Text()
		textTest := text
		if caseInsensetive {
			textTest = strings.ToLower(text)
		}
		if strings.Contains(textTest, searchQuery) {
			fileLines = append(fileLines, fmt.Sprintf("%s:%d:%s\n", pathName, n, text))
			if totalLines {
				*totalNoLines++
			}
			found = true
		}
		n++
	}

	lineSlices <- fileLines

	if found && totalFiles {
		*totalNoFiles++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	return nil
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
