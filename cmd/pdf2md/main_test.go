package main

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

// TestCLIFlags freezes the public flag set: renaming or removing a flag is a
// breaking change and must be a deliberate edit of this list.
func TestCLIFlags(t *testing.T) {
	want := []string{
		"debug", "exclude-bottom", "exclude-top", "help", "model-path",
		"no-ocr", "no-onnx", "ocr-det-model-path", "ocr-dict-path",
		"ocr-rec-model-path", "password", "runtime-path", "version",
	}

	var got []string
	for _, f := range kong.Must(&CLI).Model.Flags {
		got = append(got, f.Name)
	}
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("CLI flags changed\n got: %v\nwant: %v", got, want)
	}
}

func TestConvertPDFToMarkdown(t *testing.T) {
	pdf := writePDF(t, []run{{"Helvetica", 12, 72, 700, "Hello pdf2md"}})

	md, err := convertPDFToMarkdown(pdf, "", convertOptions{})
	if err != nil {
		t.Fatalf("convertPDFToMarkdown() error = %v", err)
	}
	if !strings.Contains(md, "Hello pdf2md") {
		t.Errorf("markdown = %q, want it to contain %q", md, "Hello pdf2md")
	}
}

func TestConvertPDFToMarkdown_MissingFile(t *testing.T) {
	_, err := convertPDFToMarkdown(filepath.Join(t.TempDir(), "nope.pdf"), "", convertOptions{})
	if err == nil {
		t.Fatal("convertPDFToMarkdown() error = nil, want error for missing file")
	}
}

func TestConvertPDFToMarkdown_EncryptedPDF(t *testing.T) {
	pdf := filepath.Join("..", "..", "internal", "extractor", "testdata", "encrypted", "aes256_R6.pdf")

	md, err := convertPDFToMarkdown(pdf, "", convertOptions{password: "secret"})
	if err != nil {
		t.Fatalf("convertPDFToMarkdown() error = %v", err)
	}
	if !strings.Contains(md, "Hello encrypted pdf2md") {
		t.Errorf("markdown = %q, want the decrypted text", md)
	}

	if _, err := convertPDFToMarkdown(pdf, "", convertOptions{}); err == nil {
		t.Error("convertPDFToMarkdown() without password: error = nil, want error")
	}
}

func TestPasswordFlagAndEnv(t *testing.T) {
	t.Cleanup(func() { CLI.Password = "" })
	t.Setenv("PDF2MD_PASSWORD", "from-env")

	k := kong.Must(&CLI)
	if _, err := k.Parse([]string{"in.pdf"}); err != nil {
		t.Fatal(err)
	}
	if CLI.Password != "from-env" {
		t.Errorf("Password = %q, want the environment value", CLI.Password)
	}

	if _, err := k.Parse([]string{"--password=from-flag", "in.pdf"}); err != nil {
		t.Fatal(err)
	}
	if CLI.Password != "from-flag" {
		t.Errorf("Password = %q, want the flag to win over the environment", CLI.Password)
	}
}

// Without an output file the Markdown goes to stdout, so conversion itself
// must not print there: a warning would land inside the Markdown.
func TestStdoutIsOnlyMarkdown(t *testing.T) {
	pdf := writePDFStream(t, "/Fm1 Do BT /F1 12 Tf 72 700 Td (hello) Tj ET")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	md, convErr := convertPDFToMarkdown(pdf, "", convertOptions{})
	os.Stdout = stdout
	_ = w.Close()
	printed, _ := io.ReadAll(r)
	if convErr != nil {
		t.Fatalf("convertPDFToMarkdown() error = %v", convErr)
	}
	if !strings.Contains(md, "hello") {
		t.Errorf("Markdown = %q, want the page text", md)
	}
	if len(printed) != 0 {
		t.Errorf("conversion wrote to stdout: %q", printed)
	}
}
