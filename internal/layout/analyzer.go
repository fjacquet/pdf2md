// Package layout provides PDF document layout analysis functionality.
package layout

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

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
	// BodyFontSize is the document's body text size (see BodyFontSize); 0
	// measures each page alone, where references or figure labels can outvote
	// the body text and turn its lines into headings.
	BodyFontSize float64
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
			// Runs of numbered body-style lines are lists: demoteNumberedRuns.
			Condition: func(text string, _, _ float64) bool {
				return a.isNumberedHeader(stripBold(text))
			},
			Type: ElementTypeHeader,
		},
		{
			Name: "Header (Lettered italic)",
			// IEEE subsections: "*C. Datasets*" alone, italic at body size.
			Condition: func(text string, _, _ float64) bool { return letteredItalic.MatchString(text) },
			Type:      ElementTypeHeader,
		},
		{
			Name: "Header (Capitals)",
			// IEEE small-caps section titles ("REFERENCES") sit at body size or below.
			Condition: func(text string, _, _ float64) bool { return isCapsHeading(stripBold(text)) },
			Type:      ElementTypeHeader,
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

	// Step 1: Reading order: content stream, XY-cut fallback
	orderedBlocks := OrderBlocks(blocks)

	// Step 3: Detect body font size
	bodyFontSize := a.BodyFontSize
	if bodyFontSize <= 0 {
		bodyFontSize = a.detectBodyFontSize(orderedBlocks)
	}

	minSize := 50.0
	if a.Config != nil {
		minSize = a.Config.MinImageSize
	}

	// Step 3.5: Drop plot labels inside figures (before any merge)
	orderedBlocks = dropFigureText(orderedBlocks, graphics, bodyFontSize, minSize)

	// Step 4: Convert blocks to initial elements
	elements := a.blocksToElements(orderedBlocks)

	// Step 4.5: Add images. An empty ID means there is no file to link.
	for _, img := range images {
		if img.ID == "" {
			continue
		}
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
		if vg.ID == "" {
			continue
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
	elements = a.demoteBoldRuns(elements, bodyFontSize)
	elements = a.demoteNumberedRuns(elements, bodyFontSize)

	// Step 6.5: Merge list item continuations (paragraphs following list items)
	elements = a.MergeListContinuations(elements)

	// Step 6.6: Merge consecutive paragraph lines
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

	// NOTE: Header level normalization is now done globally in the main pipeline
	// after all pages are collected, not per-page

	return elements
}

// detectBodyFontSize is BodyFontSize, or 12 when no block has a size.
func (a *Analyzer) detectBodyFontSize(blocks []extractor.TextBlock) float64 {
	if size := BodyFontSize(blocks); size > 0 {
		return size
	}
	return 12.0
}

// BodyFontSize is the font size carrying the most characters (ties go to the
// smaller size), 0 when no block has a size. Counting characters, not blocks,
// keeps math glyphs and table cells from outvoting the body text.
func BodyFontSize(blocks []extractor.TextBlock) float64 {
	chars := map[float64]int{}
	for _, b := range blocks {
		if b.FontSize > 0 {
			chars[b.FontSize] += utf8.RuneCountInString(b.Text)
		}
	}
	best, most := 0.0, 0
	for size, n := range chars {
		if n > most || (n == most && size < best) {
			best, most = size, n
		}
	}
	return best
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
			Content:    content,
			X:          block.X,
			Y:          block.Y,
			Width:      block.Width,
			Height:     block.Height,
			FontSize:   block.FontSize,
			SpaceWidth: block.SpaceWidth,
			Type:       ElementTypeParagraph,
			LinkURI:    block.LinkURI,
		}
		elements = append(elements, element)
	}

	return elements
}

// isCapsHeading reports a short title in capitals: 1 to 3 words of at least
// two capital letters each, 4 letters or more ("REFERENCES", "RELATED WORK").
// A word that is a Roman numeral ("TABLE II") marks a caption label.
func isCapsHeading(text string) bool {
	words := strings.Fields(text)
	if len(words) == 0 || len(words) > 3 || utf8.RuneCountInString(strings.Join(words, "")) < 4 {
		return false
	}
	for _, w := range words {
		if utf8.RuneCountInString(w) < 2 || romanNumeral.MatchString(w) || strings.IndexFunc(w, func(r rune) bool { return !unicode.IsUpper(r) }) >= 0 {
			return false
		}
	}
	return true
}

var letteredItalic = regexp.MustCompile(`^\*[A-Z]\. [^*]{1,80}[^.*]\*$`)

var romanNumeral = regexp.MustCompile(`^X{0,3}(IX|IV|V?I{0,3})$`)
