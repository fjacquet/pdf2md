package ocr

import "testing"

func TestCTCDecode(t *testing.T) {
	chars := []string{"a", "b", "c"} // classes: 0=blank, 1=a, 2=b, 3=c

	tests := []struct {
		name     string
		logits   []float32
		T, C     int
		wantText string
		// wantConf is a lower bound (greedy should pick the max; assertions
		// use ≥ to keep tests robust against float formatting).
		wantConf float64
	}{
		{
			name:     "empty input returns empty",
			logits:   []float32{},
			T:        0,
			C:        4,
			wantText: "",
			wantConf: 0,
		},
		{
			name: "all blank produces empty text",
			logits: []float32{
				1, 0, 0, 0,
				1, 0, 0, 0,
			},
			T:        2,
			C:        4,
			wantText: "",
			wantConf: 0,
		},
		{
			name: "collapse consecutive duplicates",
			// classes: a, a, blank, b -> "ab"
			logits: []float32{
				0, 1, 0, 0,
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 0, 1, 0,
			},
			T:        4,
			C:        4,
			wantText: "ab",
			wantConf: 0.9,
		},
		{
			name: "blank between duplicates keeps both",
			// classes: a, blank, a -> "aa"
			logits: []float32{
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 1, 0, 0,
			},
			T:        3,
			C:        4,
			wantText: "aa",
			wantConf: 0.9,
		},
		{
			name: "mismatched shape returns empty",
			logits: []float32{
				0, 1, 0, 0,
			},
			T:        2, // declared but only 1 step of data
			C:        4,
			wantText: "",
			wantConf: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			text, conf := CTCDecode(tc.logits, tc.T, tc.C, chars)
			if text != tc.wantText {
				t.Errorf("text = %q, want %q", text, tc.wantText)
			}
			if conf < tc.wantConf {
				t.Errorf("confidence = %v, want ≥ %v", conf, tc.wantConf)
			}
		})
	}
}

func TestArgmax(t *testing.T) {
	idx, val := argmax([]float32{0.1, 0.5, 0.3, 0.9, 0.2})
	if idx != 3 {
		t.Errorf("argmax idx = %d, want 3", idx)
	}
	if val < 0.89 || val > 0.91 {
		t.Errorf("argmax val = %v, want ≈ 0.9", val)
	}
}
