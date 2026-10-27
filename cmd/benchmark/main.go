// Package main implements a benchmark tool for comparing pdf2md output against ground truth.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// baseline records per-document scores and the READoc commit they were measured on.
type baseline struct {
	Dataset string             `json:"dataset"`
	Scores  map[string]float64 `json:"scores"`
}

func main() {
	baselinePath := flag.String("baseline", "benchmark_baseline.json", "baseline scores file")
	updateBaseline := flag.Bool("update-baseline", false, "write the current scores as the new baseline")
	tolerance := flag.Float64("tolerance", 0.02, "allowed score drop per document before failing")
	flag.Parse()

	// Define paths
	dataDir := "READoc/data"
	outputDir := "READoc/output/pdf2md"

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
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
		if err := os.MkdirAll(filepath.Dir(outPath), 0o750); err != nil {
			log.Printf("Failed to create dir for %s: %v", outPath, err)
			continue
		}

		// Run pdf2md
		cmd := exec.Command("./pdf2md", "--no-onnx", "--no-ocr", pdfPath, outPath) //nolint:gosec // Benchmark tool, paths are from local filesystem walk
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
		gtBytes, err := os.ReadFile(gtPath) //nolint:gosec // Benchmark reads local ground truth files
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

	scores := make(map[string]float64, len(results))
	for _, r := range results {
		scores[filepath.ToSlash(r.File)] = r.Score
	}
	if *updateBaseline {
		if len(scores) == 0 {
			log.Fatal("Refusing to write an empty baseline")
		}
		if err := writeBaseline(*baselinePath, scores); err != nil {
			log.Fatalf("Failed to write baseline: %v", err)
		}
		fmt.Printf("Baseline written to %s\n", *baselinePath)
		return
	}
	if regressions := checkBaseline(*baselinePath, scores, *tolerance); len(regressions) > 0 {
		fmt.Printf("\nRegressions against %s:\n  %s\n", *baselinePath, strings.Join(regressions, "\n  "))
		os.Exit(1)
	}
}

func writeBaseline(path string, scores map[string]float64) error {
	// The dataset is a separate git clone; record which commit the scores belong to.
	rev, _ := exec.Command("git", "-C", "READoc", "rev-parse", "HEAD").Output() //nolint:gosec // fixed arguments
	data, err := json.MarshalIndent(baseline{Dataset: strings.TrimSpace(string(rev)), Scores: scores}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

// checkBaseline reports documents whose score fell more than tol below the
// baseline. A missing baseline file is not an error: there is nothing to compare.
func checkBaseline(path string, scores map[string]float64, tol float64) []string {
	data, err := os.ReadFile(path) //nolint:gosec // G304: path comes from a CLI flag of a local tool
	if err != nil {
		fmt.Printf("No baseline at %s (run with -update-baseline to create one)\n", path)
		return nil
	}
	var b baseline
	if err := json.Unmarshal(data, &b); err != nil {
		log.Fatalf("Invalid baseline %s: %v", path, err)
	}
	return findRegressions(scores, b.Scores, tol)
}

func findRegressions(scores, base map[string]float64, tol float64) []string {
	var out []string
	for file, want := range base {
		got, ok := scores[file]
		switch {
		case !ok:
			out = append(out, fmt.Sprintf("%s: missing (baseline %.4f)", file, want))
		case got < want-tol:
			out = append(out, fmt.Sprintf("%s: %.4f < baseline %.4f", file, got, want))
		}
	}
	sort.Strings(out)
	return out
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
