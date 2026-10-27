package extractor

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/pdf"
)

// The fixtures were written by qpdf (pikepdf) from a one-line PDF, so they
// check our decryption against an independent implementation. The user
// password is "secret" and the owner password "owner"; aes256_R6_blank has a
// blank user password.
func TestCustomExtractor_EncryptedPDF(t *testing.T) {
	const want = "Hello encrypted pdf2md"

	tests := []struct {
		name     string
		file     string
		password string
		wantErr  string // substring; empty means success
	}{
		{"AES-128 R4 user password", "aes128_R4.pdf", "secret", ""},
		{"AES-256 R6 user password", "aes256_R6.pdf", "secret", ""},
		{"AES-256 R6 blank user password", "aes256_R6_blank.pdf", "", ""},
		{"AES-256 R6 blank, unneeded password given", "aes256_R6_blank.pdf", "secret", ""},
		{"AES-128 R4 no password", "aes128_R4.pdf", "", "user password required"},
		{"AES-256 R6 no password", "aes256_R6.pdf", "", "user password required"},
		{"AES-128 R4 wrong password", "aes128_R4.pdf", "nope", "wrong password"},
		{"AES-256 R6 wrong password", "aes256_R6.pdf", "nope", "wrong password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, err := NewCustomExtractor(filepath.Join("testdata", "encrypted", tt.file), pdf.WithPassword(tt.password))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) || !errors.Is(err, pdf.ErrEncrypted) {
					t.Fatalf("NewCustomExtractor() error = %v, want ErrEncrypted containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCustomExtractor() error = %v", err)
			}
			defer func() { _ = ext.Close() }()

			content, err := ext.ExtractTextBlocks(1)
			if err != nil {
				t.Fatalf("ExtractTextBlocks() error = %v", err)
			}
			var got []string
			for _, b := range content.TextBlocks {
				got = append(got, b.Text)
			}
			if joined := strings.Join(got, " "); !strings.Contains(joined, want) {
				t.Errorf("text = %q, want it to contain %q", joined, want)
			}
		})
	}
}
