package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Define paths
	dataDir := "READoc/data"
	outputDir := "READoc/output/pdf2md"

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Find all PDFs
	var pdfFiles []string
	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk data directory: %v", err)
	}

	fmt.Printf("Found %d PDF files to benchmark.\n", len(pdfFiles))

	var totalScore float64
	var count int

	type Result struct {
		File  string
		Score float64
	}
	var results []Result

	for _, pdfPath := range pdfFiles {
		// Determine relative path to maintain structure in output
		relPath, err := filepath.Rel(dataDir, pdfPath)
		if err != nil {
			log.Printf("Failed to get relative path for %s: %v", pdfPath, err)
			continue
		}

		// Construct output path
		outPath := filepath.Join(outputDir, relPath)
		outPath = strings.TrimSuffix(outPath, ".pdf") + ".md"

		// Ensure output subdirectory exists
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			log.Printf("Failed to create dir for %s: %v", outPath, err)
			continue
		}

		// Run pdf2md
		cmd := exec.Command("./pdf2md", pdfPath, outPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to convert %s: %v\nStderr: %s", pdfPath, err, stderr.String())
			continue
		}

		// Determine ground truth path
		// Structure is data/<subset>/pdf/<file>.pdf -> data/<subset>/markdown/<file>.md
		// relPath is <subset>/pdf/<file>.pdf
		parts := strings.Split(relPath, string(os.PathSeparator))
		if len(parts) < 3 || parts[1] != "pdf" {
			log.Printf("Unexpected path structure for %s, skipping comparison", relPath)
			continue
		}
		// Replace 'pdf' dir with 'markdown'
		parts[1] = "markdown"
		gtRelPath := filepath.Join(parts...)
		gtPath := filepath.Join(dataDir, gtRelPath)
		gtPath = strings.TrimSuffix(gtPath, ".pdf") + ".md"

		// Read files
		genBytes, err := os.ReadFile(outPath)
		if err != nil {
			log.Printf("Failed to read generated file %s: %v", outPath, err)
			continue
		}
		gtBytes, err := os.ReadFile(gtPath)
		if err != nil {
			log.Printf("Failed to read ground truth file %s: %v", gtPath, err)
			continue
		}

		// Normalize and compare
		score := calculateSimilarity(string(genBytes), string(gtBytes))
		fmt.Printf("[%s] Similarity: %.4f\n", parts[len(parts)-1], score)

		totalScore += score
		count++
		results = append(results, Result{File: relPath, Score: score})
	}

	if count > 0 {
		avg := totalScore / float64(count)
		fmt.Printf("\nBenchmark Complete.\nProcessed: %d\nAverage Similarity: %.4f\n", count, avg)

		// Find worst
		minScore := 1.0
		var worstFile string
		for _, r := range results {
			if r.Score < minScore {
				minScore = r.Score
				worstFile = r.File
			}
		}
		fmt.Printf("Worst performing file: %s (Score: %.4f)\n", worstFile, minScore)
	} else {
		fmt.Println("No files processed.")
	}
}

func calculateSimilarity(s1, s2 string) float64 {
	r1 := []rune(normalize(s1))
	r2 := []rune(normalize(s2))

	if len(r1) == 0 && len(r2) == 0 {
		return 1.0
	}
	if len(r1) == 0 || len(r2) == 0 {
		return 0.0
	}

	dist := levenshtein(r1, r2)
	maxLen := len(r1)
	if len(r2) > maxLen {
		maxLen = len(r2)
	}

	return 1.0 - float64(dist)/float64(maxLen)
}

func normalize(s string) string {
	// Simple normalization: remove extra whitespace, lowercase
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

func levenshtein(s1, s2 []rune) int {
	len1 := len(s1)
	len2 := len(s2)

	// Optimize for space: use 2 rows
	row := make([]int, len2+1)
	for i := 0; i <= len2; i++ {
		row[i] = i
	}

	for i := 1; i <= len1; i++ {
		prev := i
		var current int
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			current = min(
				row[j]+1,      // deletion
				prev+1,        // insertion
				row[j-1]+cost, // substitution
			)
			row[j-1] = prev
			prev = current
		}
		row[len2] = prev
	}
	return row[len2]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
