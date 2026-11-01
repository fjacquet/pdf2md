package layout

import "testing"

func TestIsItalic(t *testing.T) {
	tests := map[string]bool{
		"Times-Italic":                       true,
		"ITVBSG+HelveticaNeue-BoldItalic":    true,
		"DejaVuSans-Oblique":                 true,
		"YQXQOX+NimbusRomNo9L-ReguItal":      true,
		"LREKMX+NimbusRomNo9L-MediItal":      true,
		"CQFHYB+CMTI10":                      true,
		"MUEWEX+SourceSansPro-It-Identity-H": true,
		"MinionPro-BoldIt":                   true,
		"LIFSIT+CMBX12":                      false, // "IT" only in the subset tag
		"Times-Roman":                        false,
		"EUIAKY+NimbusRomNo9L-Regu":          false,
		"Digital-Regular":                    false,
	}
	for name, want := range tests {
		if got := isItalic(name); got != want {
			t.Errorf("isItalic(%q) = %v, want %v", name, got, want)
		}
	}
}
