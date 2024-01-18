package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func searchInFile(filePath, searchWord string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lineNumber := 0
	scanner := bufio.NewScanner(file)
	// Increase the buffer size to handle longer lines
	const maxTokenSize = 64 * 1024 // 64 KB
	buf := make([]byte, maxTokenSize)
	scanner.Buffer(buf, maxTokenSize)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.Contains(line, searchWord) {
			// Use regular expression to find and highlight the search word
			re := regexp.MustCompile(searchWord)
			highlightedLine := re.ReplaceAllStringFunc(line, func(match string) string {
				return fmt.Sprintf("\033[1;31m%s\033[0m", match) // ANSI escape code for red text
			})

			fmt.Printf("File: %s, Line: %d, Content: %s\n", filePath, lineNumber, highlightedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func searchInDirectory(directory, searchWord string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := searchInFile(path, searchWord)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rsearch <directory> <search_word>")
		os.Exit(1)
	}

	searchDirectory := os.Args[1]
	searchWord := os.Args[2]

	err := searchInDirectory(searchDirectory, searchWord)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
