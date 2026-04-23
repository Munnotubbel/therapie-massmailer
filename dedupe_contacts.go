package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dedupe_contacts.go <input_csv> [output_csv]")
		return
	}

	inputPath := os.Args[1]
	outputPath := inputPath
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	if len(records) < 2 {
		fmt.Println("CSV is empty or too short.")
		return
	}

	header := records[0]
	// Skip the header and keep unique emails
	uniqueRows := make(map[string][]string)
	var orderedEmails []string

	// Find email column index
	emailIdx := -1
	for i, col := range header {
		if strings.ToLower(col) == "email" {
			emailIdx = i
			break
		}
	}

	if emailIdx == -1 {
		fmt.Println("Error: 'Email' column not found.")
		return
	}

	for i := 1; i < len(records); i++ {
		email := strings.ToLower(strings.TrimSpace(records[i][emailIdx]))
		if email == "" {
			continue
		}
		if _, exists := uniqueRows[email]; !exists {
			uniqueRows[email] = records[i]
			orderedEmails = append(orderedEmails, email)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write(header)
	for _, email := range orderedEmails {
		writer.Write(uniqueRows[email])
	}

	fmt.Printf("Deduplication complete. Original: %d rows, Unique: %d rows.\n", len(records)-1, len(orderedEmails))
}
