package main

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/layout"
)

func TestRecoverPage(t *testing.T) {
	_, err := recoverPage(func() ([]layout.Element, error) { panic("boom") })
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("panic not turned into error: %v", err)
	}
	els, err := recoverPage(func() ([]layout.Element, error) { return []layout.Element{{Content: "ok"}}, nil })
	if err != nil || len(els) != 1 {
		t.Fatalf("normal page lost: %v %v", els, err)
	}
}
