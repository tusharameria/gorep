package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	startTime := time.Now()
	caseInsensetive := false
	isRecursive := false
	totalLines := false
	totalFiles := false
	showTime := false
	coreWorkers := false
	numWorkers := 1
	flag.BoolVar(&caseInsensetive, "i", caseInsensetive, "True or False")
	flag.BoolVar(&isRecursive, "r", isRecursive, "True or False")
	flag.BoolVar(&showTime, "time", showTime, "True or False")
	flag.BoolVar(&totalLines, "totalLines", totalLines, "True or False")
	flag.BoolVar(&totalFiles, "totalFiles", totalFiles, "True or False")
	flag.BoolVar(&coreWorkers, "coreWorkers", coreWorkers, "True or False")
	flag.IntVar(&numWorkers, "workers", numWorkers, "Should be a positive number")
	flag.Parse()
	args := flag.Args()

	searchQuery := args[0]
	pathNames := args[1:]

	filePaths := make(chan string, 10)
	lineSlices := make(chan []string, 10)

	var totalNoLines int64
	var totalNoFiles int64

	var recursive func(string, bool)
	recursive = func(pathName string, firstLayer bool) {
		isFile, err := IsFile(pathName)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if isFile {
			filePaths <- pathName
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
		close(filePaths)
	}()

	numCores := runtime.NumCPU()
	numPool := numWorkers
	if coreWorkers {
		numPool = numCores
	}

	wg := sync.WaitGroup{}

	for i := 0; i < numPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pathName := range filePaths {
				if err := ProcessFile(pathName, searchQuery, caseInsensetive, totalLines, totalFiles, &totalNoLines, &totalNoFiles, lineSlices); err != nil {
					fmt.Printf("%s\n", err)
					os.Exit(1)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(lineSlices)
	}()

	for lines := range lineSlices {
		for _, line := range lines {
			fmt.Printf("%s", line)
		}
	}
	fmt.Printf("Number of workers in pool : %d\n", numPool)

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

func ProcessFile(pathName, searchQuery string, caseInsensetive, totalLines, totalFiles bool, totalNoLines, totalNoFiles *int64, lineSlices chan<- []string) error {
	file, err := os.Open(pathName)
	if err != nil {
		return err
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
				atomic.AddInt64(totalNoLines, 1)
			}
			found = true
		}
		n++
	}

	lineSlices <- fileLines

	if found && totalFiles {
		atomic.AddInt64(totalNoFiles, 1)
	}

	if err := scanner.Err(); err != nil {
		return err
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
