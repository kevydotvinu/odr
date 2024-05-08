package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Define flag for search query
	db := flag.String("db", "", "Path to the OVN database file")
	all := flag.Bool("all", false, "Display all output")
	search := flag.String("search", "", "Search for a specific table")
	uuid := flag.String("uuid", "", "UUID to search for")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Count the number of flags provided
	var flagCount int
	flag.Visit(func(f *flag.Flag) {
		flagCount++
	})

	// Check the number of flags provided
	if flagCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: At least one flag is required.\n")
		flag.Usage()
		os.Exit(1)
	} else if flagCount > 2 {
		fmt.Fprintf(os.Stderr, "Error: Only two flag can be used at a time.\n")
		flag.Usage()
		os.Exit(1)
	}

	// Find the leader_nbdb file recursively from the current directory
	var filePath string
	if *db != "" {
		filePath = *db
	} else {
		var err error
		filePath, err = findFile("leader_nbdb")
		if err != nil {
			fmt.Println("Error finding file:", err)
			return
		}
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Calculate buffer size considering the file size
	var bufferSize = 1000
	if fileSize > int64(bufferSize) {
		bufferSize = int(fileSize)
	}

	// Create a scanner with a larger buffer size
	scanner := bufio.NewScanner(file)
	buf := make([]byte, bufferSize)
	scanner.Buffer(buf, bufferSize)

	counter := 0

	// Iterate through each line
	for scanner.Scan() {
		counter++

		// Check if the line number is odd
		if counter%2 != 0 {
			if *search != "" {
				// Search for the query in the current line
				if strings.Contains(scanner.Text(), *search) {
					if !strings.Contains(scanner.Text(), "OVSDB CLUSTER") {
						output := replaceUnicodeEscape(scanner.Text())
						fmt.Println(output)
					}
				}
			} else if *all == true {
				if !strings.Contains(scanner.Text(), "OVSDB CLUSTER") {
					output := replaceUnicodeEscape(scanner.Text())
					fmt.Println(output)
				}
			} else if *uuid != "" {
				if !strings.Contains(scanner.Text(), "OVSDB CLUSTER") {
					output := replaceUnicodeEscape(scanner.Text())
					fmt.Println(output)
				}
			}
		} else {
			// Parse the JSON using the "encoding/json" package
			var data map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &data)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}

			if *search != "" {
				// Print only the part of JSON containing the search term
				printJSON(data, *search)
			} else if *uuid != "" {
				printJSON(data, *uuid)
			} else {
				// Print the formatted JSON if no search term is provided
				prettyJSON, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					fmt.Println("Error formatting JSON:", err)
					return
				}
				// Search and replace the unicode in the string
				output := replaceUnicodeEscape(string(prettyJSON))
				fmt.Println(output)
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}

// Function to print only the part of JSON containing the search term
func printJSON(data interface{}, search string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == search {
				prettyJSON, err := json.MarshalIndent(value, "", "  ")
				if err != nil {
					fmt.Println("Error formatting JSON:", err)
					return
				}
				output := replaceUnicodeEscape(string(prettyJSON))
				fmt.Println(output)
			} else {
				printJSON(value, search)
			}
		}
	case []interface{}:
		for _, element := range v {
			printJSON(element, search)
		}
	}
}

// Function to search and replace unicode
func replaceUnicodeEscape(input string) string {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		unicodeCode, _ := strconv.ParseInt(match[2:], 16, 32)
		return string(unicodeCode)
	})
}

// Function to recursively find a file with a given name
func findFile(fileName string) (string, error) {
	var filePath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == fileName {
			filePath = path
			return nil
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("file %s not found", fileName)
	}
	return filePath, nil
}
