package layout

import (
	"testing"
)

func TestMergeParagraphLines_Aggressive(t *testing.T) {
	analyzer := NewAnalyzer()

	// Case 1: Distinct paragraphs (should NOT merge)
	// - Close vertical gap (simulating tight spacing)
	// - First ends with terminator
	// - Second starts with Uppercase
	elements := []Element{
		{Type: ElementTypeParagraph, Content: "End of sentence.", Y: 100, FontSize: 10, Height: 10},
		{Type: ElementTypeParagraph, Content: "Start of new paragraph.", Y: 85, FontSize: 10, Height: 10}, // Gap = 15 (1.5 * FontSize) < 1.8 * FontSize
	}

	merged := analyzer.MergeParagraphLines(elements)

	if len(merged) != 2 {
		t.Errorf("Case 1: Expected 2 elements, got %d. Content: '%s'", len(merged), merged[0].Content)
	}

	// Case 2: Continuation (should merge)
	// - Close vertical gap
	// - First does NOT end with terminator
	elements2 := []Element{
		{Type: ElementTypeParagraph, Content: "Sentence continues", Y: 100, FontSize: 10, Height: 10},
		{Type: ElementTypeParagraph, Content: "on next line.", Y: 85, FontSize: 10, Height: 10},
	}

	merged2 := analyzer.MergeParagraphLines(elements2)

	if len(merged2) != 1 {
		t.Errorf("Case 2: Expected 1 element, got %d", len(merged2))
	}
}

func TestReproJSONCodeBlock(t *testing.T) {
	// Simulate JSON content that was being misclassified as Table
	// | "ignition": { | "version": "3.2.0" |

	analyzer := NewAnalyzer()

	// Case 1: JSON line with indentation (simulated)
	// "    \"ignition\": {"
	text1 := "    \"ignition\": {"
	if !analyzer.isCodeBlock(text1) {
		t.Errorf("Expected '%s' to be detected as CodeBlock", text1)
	}

	// Case 2: JSON line with key-value
	// "version": "3.2.0"
	text2 := "\"version\": \"3.2.0\""
	if !analyzer.isCodeBlock(text2) {
		t.Errorf("Expected '%s' to be detected as CodeBlock", text2)
	}

	// Case 3: JSON line with closer and opener
	// }, "passwd": {
	text3 := "}, \"passwd\": {"
	if !analyzer.isCodeBlock(text3) {
		t.Errorf("Expected '%s' to be detected as CodeBlock", text3)
	}

	// Case 4: SSH Key inside list
	// "sshAuthorizedKeys": [ "ssh-rsa AAA..." ]
	text4 := "    \"sshAuthorizedKeys\": [ \"ssh-rsa AAAAB3NzaC1yc....\" ]"
	if !analyzer.isCodeBlock(text4) {
		t.Errorf("Expected '%s' to be detected as CodeBlock", text4)
	}

	// Case 5: Structural line (brackets only)
	// |  |   } |  ] | -> "      }   ]"
	text5 := "      }   ]"
	if !analyzer.isCodeBlock(text5) {
		t.Errorf("Expected '%s' to be detected as CodeBlock", text5)
	}
}

func TestReproTableDetection(t *testing.T) {
	analyzer := NewAnalyzer()

	// Simulate "Component" and "Description" on the same line with a gap
	// FontSize = 12
	// Block 1: X=10, Width=50 (End=60)
	// Block 2: X=90, Width=50 (Start=90)
	// Gap = 30
	// Current WideGapThreshold = 12 * 3 = 36.
	// Gap 30 < 36, so it merges with space -> "Component Description" -> Paragraph.

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Component", X: 10, Y: 100, Width: 50, FontSize: 12, Height: 12},
		{Type: ElementTypeParagraph, Content: "Description", X: 100, Y: 100, Width: 50, FontSize: 12, Height: 12},
	}

	merged := analyzer.MergeElements(elements)

	// Check if it's a table
	// We need to run classifyElement on the result
	for i := range merged {
		analyzer.classifyElement(&merged[i], nil, 12.0)
	}

	if len(merged) != 1 {
		t.Errorf("Expected 1 merged element, got %d", len(merged))
	} else {
		if merged[0].Type != ElementTypeTable {
			t.Errorf("Expected ElementTypeTable, got %v. Content: '%s'", merged[0].Type, merged[0].Content)
		}
	}
}
