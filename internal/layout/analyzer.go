// Package layout provides PDF document layout analysis functionality.
package layout

import (
	"strings"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

// Element is an alias for types.Element for backward compatibility.
type Element = types.Element

// ElementType is an alias for types.ElementType for backward compatibility.
type ElementType = types.ElementType

// Re-export element type constants
const (
	ElementTypeHeader     = types.ElementTypeHeader
	ElementTypeParagraph  = types.ElementTypeParagraph
	ElementTypeCodeBlock  = types.ElementTypeCodeBlock
	ElementTypeList       = types.ElementTypeList
	ElementTypeTable      = types.ElementTypeTable
	ElementTypeAdmonition = types.ElementTypeAdmonition
	ElementTypeImage      = types.ElementTypeImage
)

// Rule defines a classification rule for the layout analyzer
type Rule struct {
	Name      string
	Condition func(text string, fontSize, bodyFontSize float64) bool
	Type      ElementType
}

// ExclusionZone defines areas to ignore (e.g., headers, footers)
type ExclusionZone struct {
	Top    float64 // Height from top to ignore
	Bottom float64 // Height from bottom to ignore
}

// Analyzer analyzes the layout of text blocks
type Analyzer struct {
	ColumnGapThreshold float64
	HeaderSizeRatio    float64
	Rules              []Rule
	Exclusion          ExclusionZone
	Config             *Config
}

// NewAnalyzer creates a new Analyzer with default configuration
func NewAnalyzer() *Analyzer {
	return NewAnalyzerWithConfig(DefaultConfig())
}

// NewAnalyzerWithConfig creates a new Analyzer with the specified configuration
func NewAnalyzerWithConfig(config *Config) *Analyzer {
	if config == nil {
		config = DefaultConfig()
	}
	a := &Analyzer{
		ColumnGapThreshold: config.ColumnGapThreshold,
		HeaderSizeRatio:    config.HeaderSizeRatio,
		Config:             config,
	}

	// Initialize Rules
	a.Rules = []Rule{
		{
			Name: "Admonition",
			Condition: func(text string, _, _ float64) bool {
				cleanText := text
				if strings.HasPrefix(cleanText, "**") {
					cleanText = strings.TrimPrefix(cleanText, "**")
				} else if strings.HasPrefix(cleanText, "*") {
					cleanText = strings.TrimPrefix(cleanText, "*")
				}
				keywords := []string{"IMPORTANT", "WARNING", "NOTE", "TIP", "CAUTION"}
				for _, kw := range keywords {
					if strings.HasPrefix(cleanText, kw) {
						return true
					}
				}
				return false
			},
			Type: ElementTypeAdmonition,
		},
		{
			Name: "Code Block (Keywords)",
			Condition: func(text string, _, _ float64) bool {
				return a.isCodeBlock(text)
			},
			Type: ElementTypeCodeBlock,
		},
		{
			Name: "Table Row (Wide Gaps)",
			Condition: func(text string, _, _ float64) bool {
				// First check if this looks like an equation
				if a.isLikelyEquationText(text) {
					return false
				}

				// Check if this looks like an author block
				if a.isLikelyAuthorBlock(text) {
					return false
				}

				// Count 4-space gaps (columns)
				gapCount := strings.Count(text, "    ")
				if gapCount == 0 {
					if strings.Count(text, "   ") < 2 {
						return false
					}
					// 3-space gaps found, continue with segment analysis
				}

				segments := strings.Split(text, "    ")
				if len(segments) < 2 {
					segments = strings.Split(text, "   ")
				}

				if len(segments) < 2 {
					return false
				}

				minLen := len(segments[0])
				maxLen := len(segments[0])
				for _, seg := range segments[1:] {
					segLen := len(strings.TrimSpace(seg))
					if segLen < minLen {
						minLen = segLen
					}
					if segLen > maxLen {
						maxLen = segLen
					}
				}

				if minLen > 0 && maxLen > minLen*10 {
					return false
				}

				hasNumericContent := false
				hasStructuredData := false
				for _, seg := range segments {
					seg = strings.TrimSpace(seg)
					if containsSignificantNumbers(seg) {
						hasNumericContent = true
					}
					if len(seg) > 0 && len(seg) < 30 {
						hasStructuredData = true
					}
				}

				if !hasNumericContent && !hasStructuredData {
					allShort := true
					for _, seg := range segments {
						if len(strings.TrimSpace(seg)) > 40 {
							allShort = false
							break
						}
					}
					if !allShort {
						return false
					}
				}

				return true
			},
			Type: ElementTypeTable,
		},
		{
			Name: "Header (Numbered)",
			Condition: func(text string, _, _ float64) bool {
				return a.isNumberedHeader(text)
			},
			Type: ElementTypeHeader,
		},
		{
			Name: "List Item (Bullet/Number)",
			Condition: func(text string, _, _ float64) bool {
				return a.isListItem(text)
			},
			Type: ElementTypeList,
		},
	}

	return a
}

// Analyze converts raw text blocks into structured elements
func (a *Analyzer) Analyze(content *extractor.PageContent) []Element {
	if content == nil {
		return nil
	}
	blocks := content.TextBlocks
	images := content.Images
	graphics := content.Graphics

	if len(blocks) == 0 && len(images) == 0 && len(graphics) == 0 {
		return nil
	}

	// Step 0: Filter blocks based on exclusion zones
	if a.Exclusion.Top > 0 || a.Exclusion.Bottom > 0 {
		filteredBlocks := make([]extractor.TextBlock, 0, len(blocks))
		pageHeight := content.PageHeight

		for _, block := range blocks {
			if a.Exclusion.Bottom > 0 && block.Y < a.Exclusion.Bottom {
				continue
			}
			if a.Exclusion.Top > 0 && pageHeight > 0 && block.Y > (pageHeight-a.Exclusion.Top) {
				continue
			}
			filteredBlocks = append(filteredBlocks, block)
		}
		blocks = filteredBlocks
	}

	// Step 0.5: Apply links to blocks
	blocks = ApplyLinksToBlocks(blocks, content.Links)

	// Step 1: Sort blocks using Recursive XY-Cut
	orderedBlocks := SortBlocks(blocks)

	// Step 3: Detect body font size
	bodyFontSize := a.detectBodyFontSize(orderedBlocks)

	// Step 4: Convert blocks to initial elements
	elements := a.blocksToElements(orderedBlocks)

	// Step 4.5: Add images
	for _, img := range images {
		elements = append(elements, Element{
			Type:    ElementTypeImage,
			Content: img.ID,
			X:       img.X,
			Y:       img.Y,
			Width:   img.Width,
			Height:  img.Height,
		})
	}

	// Step 4.6: Add vector graphics
	for _, vg := range graphics {
		minSize := 50.0
		if a.Config != nil {
			minSize = a.Config.MinImageSize
		}
		if vg.Width < minSize || vg.Height < minSize {
			continue
		}
		elements = append(elements, Element{
			Type:    ElementTypeImage,
			Content: vg.ID,
			X:       vg.X,
			Y:       vg.Y,
			Width:   vg.Width,
			Height:  vg.Height,
		})
	}

	// Step 4.7: Merge consecutive elements with same LinkURI
	elements = a.MergeSameLinkElements(elements)

	// Step 5: Merge consecutive elements on same line
	elements = a.MergeElements(elements)

	// Step 6: Classify elements
	var prev *Element
	for i := range elements {
		if elements[i].Type == ElementTypeImage {
			continue
		}
		a.classifyElement(&elements[i], prev, bodyFontSize)
		prev = &elements[i]
	}

	// Step 5.5: Merge consecutive paragraph lines
	elements = a.MergeParagraphLines(elements)

	// Step 7: Merge consecutive code blocks
	elements = a.MergeCodeBlocks(elements)

	// Step 8: Merge consecutive table rows
	elements = a.MergeTableRows(elements)

	// Step 9: Clean math symbols and merge fragmented math
	elements = a.CleanMathSymbols(elements)
	elements = a.CleanFragmentedMath(elements)

	// Step 10: Filter page numbers
	elements = a.FilterPageNumbers(elements)

	// Step 11: Normalize header levels
	elements = a.NormalizeHeaderLevels(elements)

	return elements
}

// detectBodyFontSize finds the most common font size
func (a *Analyzer) detectBodyFontSize(blocks []extractor.TextBlock) float64 {
	if len(blocks) == 0 {
		return 12.0
	}

	fontSizes := make(map[float64]int)
	for _, block := range blocks {
		if block.FontSize > 0 {
			fontSizes[block.FontSize]++
		}
	}

	var maxCount int
	bodySize := 12.0
	for size, count := range fontSizes {
		if count > maxCount {
			maxCount = count
			bodySize = size
		}
	}

	return bodySize
}

// blocksToElements converts text blocks to initial elements
func (a *Analyzer) blocksToElements(blocks []extractor.TextBlock) []Element {
	elements := make([]Element, 0, len(blocks))

	for _, block := range blocks {
		content := block.Text
		if strings.TrimSpace(content) != "" {
			switch {
			case isMathFont(block.FontName):
				content = "$" + content + "$"
			case isBold(block.FontName):
				content = "**" + content + "**"
			case isItalic(block.FontName):
				content = "*" + content + "*"
			}
		}

		element := Element{
			Content:  content,
			X:        block.X,
			Y:        block.Y,
			Width:    block.Width,
			Height:   block.Height,
			FontSize: block.FontSize,
			Type:     ElementTypeParagraph,
			LinkURI:  block.LinkURI,
		}
		elements = append(elements, element)
	}

	return elements
}
