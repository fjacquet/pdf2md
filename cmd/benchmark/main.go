package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	// Paths
	groundTruthDir := "READoc/data/arxiv/markdown"
	predictionDir := "READoc/output/arxiv/pdf2md"

	// Get list of ground truth files
	files, err := os.ReadDir(groundTruthDir)
	if err != nil {
		fmt.Printf("Error reading ground truth directory: %v\n", err)
		os.Exit(1)
	}

	totalScore := 0.0
	count := 0

	fmt.Println("Benchmarking pdf2md against READoc (ArXiv subset)...")
	fmt.Println("---------------------------------------------------")
	fmt.Printf("%-20s | %-10s | %s\n", "Filename", "Similarity", "Status")
	fmt.Println("---------------------------------------------------")

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		gtPath := filepath.Join(groundTruthDir, file.Name())
		predPath := filepath.Join(predictionDir, file.Name())

		// Check if prediction exists
		if _, err := os.Stat(predPath); os.IsNotExist(err) {
			fmt.Printf("%-20s | %-10s | %s\n", file.Name(), "N/A", "MISSING")
			continue
		}

		// Read files
		gtContent, err := os.ReadFile(gtPath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", gtPath, err)
			continue
		}
		predContent, err := os.ReadFile(predPath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", predPath, err)
			continue
		}

		// Calculate similarity
		score := calculateSimilarity(string(gtContent), string(predContent))
		totalScore += score
		count++

		status := "FAIL"
		if score > 0.5 { // Threshold for "PASS"
			status = "PASS"
		}

		fmt.Printf("%-20s | %.2f       | %s\n", file.Name(), score, status)
	}

	fmt.Println("---------------------------------------------------")
	if count > 0 {
		avgScore := totalScore / float64(count)
		fmt.Printf("Average Similarity: %.2f\n", avgScore)
		if avgScore > 0.5 {
			fmt.Println("Overall Result: PASS CERTIFICATE GRANTED")
		} else {
			fmt.Println("Overall Result: FAIL")
		}
	} else {
		fmt.Println("No files benchmarked.")
	}
}

// calculateSimilarity computes a simple Jaccard similarity of words
func calculateSimilarity(s1, s2 string) float64 {
	words1 := tokenize(s1)
	words2 := tokenize(s2)

	set1 := make(map[string]bool)
	for _, w := range words1 {
		set1[w] = true
	}

	intersection := 0
	union := len(set1)

	set2 := make(map[string]bool)
	for _, w := range words2 {
		if set2[w] {
			continue
		}
		if set1[w] {
			intersection++
		} else {
			union++
		}
		set2[w] = true
	}

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
