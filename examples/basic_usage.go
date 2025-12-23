// Package main provides an example of using pdf2md as a library.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/layout"
	"github.com/fjacquet/pdf2md/internal/markdown"
)

// This example demonstrates how to use pdf2md as a library
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run basic_usage.go <input.pdf>")
		os.Exit(1)
	}

	pdfPath := os.Args[1]

	// Step 1: Open the PDF file
	// Step 2: Create and initialize the PDF extractor
	// Initialize the extractor
	ext, err := extractor.NewCustomExtractor(pdfPath)
	if err != nil {
		log.Fatalf("Failed to initialize PDF extractor: %v", err)
	}
	defer func() { _ = ext.Close() }()

	pageCount, err := ext.GetPageCount()
	if err != nil {
		log.Fatalf("Failed to get page count: %v", err) //nolint:gocritic // exitAfterDefer: intentional early exit on error
	}
	pc, err := ext.GetPageCount()
	if err != nil {
		log.Fatalf("Failed to get page count: %v", err) //nolint:gocritic // exitAfterDefer: intentional early exit on error
	}
	fmt.Printf("PDF has %d pages\n", pc)

	// Step 3: Create layout analyzer
	analyzer := layout.NewAnalyzer()

	// Step 4: Create markdown builder
	builder := markdown.NewBuilder()

	// Step 5: Process each page
	var allElements []layout.Element

	for page := 1; page <= pageCount; page++ {
		// Extract text from page 1 (0-indexed)
		content, err := ext.ExtractTextBlocks(page - 1)
		if err != nil {
			log.Fatalf("Failed to extract text: %v", err)
		}

		fmt.Printf("Page %d: extracted %d text blocks\n", page, len(content.TextBlocks))

		// Analyze layout (detect columns, headers, etc.)
		elements := analyzer.Analyze(content)

		// Merge consecutive text on same line
		elements = analyzer.MergeElements(elements)

		allElements = append(allElements, elements...)

		// Add page separator
		if page < pageCount {
			allElements = append(allElements, layout.Element{
				Type:    layout.ElementTypeParagraph,
				Content: "---",
			})
		}
	}

	// Step 6: Build markdown
	md := builder.Build(allElements)

	// Step 7: Output result
	fmt.Println("Successfully converted PDF to Markdown!")
	fmt.Println(md)
}
