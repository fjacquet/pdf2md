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
	file, err := os.Open(pdfPath)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}
	defer file.Close()

	// Step 2: Create and initialize the PDF extractor
	ext := extractor.NewLedongthucExtractor()
	if err := ext.Open(file); err != nil {
		log.Fatalf("Failed to initialize PDF extractor: %v", err)
	}
	defer ext.Close()

	fmt.Printf("PDF has %d pages\n", ext.GetPageCount())

	// Step 3: Create layout analyzer with custom settings
	analyzer := layout.NewAnalyzer()
	analyzer.ColumnGapThreshold = 50.0 // Adjust for your PDFs
	analyzer.HeaderSizeRatio = 1.3     // Adjust header detection sensitivity

	// Step 4: Create markdown builder
	builder := markdown.NewBuilder()

	// Step 5: Process each page
	var allElements []layout.Element

	for page := 1; page <= ext.GetPageCount(); page++ {
		// Extract text blocks with positioning
		blocks, err := ext.ExtractTextBlocks(page)
		if err != nil {
			log.Printf("Warning: failed to extract page %d: %v", page, err)
			continue
		}

		fmt.Printf("Page %d: extracted %d text blocks\n", page, len(blocks))

		// Analyze layout (detect columns, headers, etc.)
		elements := analyzer.Analyze(blocks)

		// Merge consecutive text on same line
		elements = analyzer.MergeElements(elements)

		allElements = append(allElements, elements...)

		// Add page separator
		if page < ext.GetPageCount() {
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
