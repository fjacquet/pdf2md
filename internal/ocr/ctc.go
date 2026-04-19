package ocr

import (
	"math"
	"strings"
)

// CTCDecode performs greedy CTC decoding on a time-major class-probability tensor.
//
// Pure function. logits is shape [T, C] flattened row-major (T time steps,
// C classes). Index 0 is the CTC blank token (convention: blank prepended to
// the dictionary, so class index k>0 maps to chars[k-1]).
//
// Decoding:
//  1. For each time step, pick argmax class (greedy).
//  2. Collapse consecutive duplicates.
//  3. Drop blanks.
//  4. Confidence is the mean of the max probabilities over kept (non-blank,
//     non-duplicate) time steps.
func CTCDecode(logits []float32, timeSteps, numClasses int, chars []string) (string, float64) {
	if len(logits) != timeSteps*numClasses || timeSteps == 0 || numClasses == 0 {
		return "", 0
	}

	var text strings.Builder
	var sumConf float64
	kept := 0
	prev := -1

	for t := range timeSteps {
		cls, prob := argmax(logits[t*numClasses : (t+1)*numClasses])
		if cls == 0 || cls == prev {
			prev = cls
			continue
		}
		if cls-1 < len(chars) {
			text.WriteString(chars[cls-1])
			sumConf += prob
			kept++
		}
		prev = cls
	}

	if kept == 0 {
		return "", 0
	}
	return text.String(), sumConf / float64(kept)
}

// argmax returns the index and value of the largest element.
// Pure function.
func argmax(row []float32) (int, float64) {
	bestIdx := 0
	bestVal := float32(math.Inf(-1))
	for i, v := range row {
		if v > bestVal {
			bestVal = v
			bestIdx = i
		}
	}
	return bestIdx, float64(bestVal)
}
