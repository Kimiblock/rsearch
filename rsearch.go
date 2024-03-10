package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const maxTokenSize = 128 * 1024

func isBinary(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512) // Read the first 512 bytes to detect binary files
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	return false
}

func isUTF8(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Use the golang.org/x/text/encoding/unicode package to detect UTF-8 encoding
	decoder := unicode.UTF8.NewDecoder()
	_, err = transform.NewReader(file, decoder).Read(make([]byte, 512))

	return err == nil
}

func searchInFile(filePath, searchWord string) error {
	if isBinary(filePath) || !isUTF8(filePath) {
		return nil // Skip binary or non-UTF-8 files
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lineNumber := 0
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lineNumber++

		if strings.Contains(line, searchWord) {
			// Use regular expression to find and highlight the search word
			re := regexp.MustCompile(searchWord)
			highlightedLine := re.ReplaceAllStringFunc(line, func(match string) string {
				return fmt.Sprintf("\033[1;31m%s\033[0m", match) // ANSI escape code for red text
			})

			// Print both original and changed lines
			fmt.Printf("File: %s, Line: %d\n", filePath, lineNumber)
			fmt.Printf("Hit:  %s\n", highlightedLine)
			fmt.Println(strings.Repeat("-", 40))
		}
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
