package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleXObject handles Do
func (in *Interpreter) handleXObject(op string) error {
	if op == "Do" { // Invoke named XObject: name Do
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Do")
		}
		name, ok := in.Stack[len(in.Stack)-1].(Name)
		if !ok {
			return fmt.Errorf("invalid operand for Do: expected Name")
		}
		in.Stack = in.Stack[:len(in.Stack)-1]
		if err := in.processXObject(name); err != nil {
			return err
		}
	}
	return nil
}

func (in *Interpreter) processXObject(name Name) error {
	// Look up XObject in Resources
	if in.Resources == nil {
		// Some PDFs might call Do without resources if it's a standard name? Unlikely.
		return fmt.Errorf("no resources dictionary")
	}

	key := Name("XObject")
	xObjectDict, ok := in.Resources[key].(Dictionary)
	if !ok {
		// If Do is called but no XObject dict, it's an error in the PDF or our parsing.
		// But we shouldn't crash or fail the page if possible.
		// Just log warning and return.
		fmt.Printf("Warning: Do operator called for %s but no XObject dictionary found\n", name)
		return nil
	}

	// Resolve indirect reference if needed
	obj := xObjectDict[name]
	if obj == nil {
		return fmt.Errorf("XObject %s not found in resources", name)
	}

	if ref, ok := obj.(IndirectRef); ok {
		// Need access to reader to resolve reference.
		if in.FontManager != nil && in.FontManager.reader != nil {
			resolved, err := in.FontManager.reader.ReadObject(ref.ObjectNumber)
			if err != nil {
				return err
			}
			obj = resolved
		} else {
			return fmt.Errorf("cannot resolve XObject reference: no reader available")
		}
	}

	var stream *Stream
	if s, ok := obj.(Stream); ok {
		stream = &s
	} else if s, ok := obj.(*Stream); ok {
		stream = s
	} else {
		return fmt.Errorf("XObject %s is not a stream (type %T)", name, obj)
	}

	subtype, _ := stream.Dictionary[Name("Subtype")].(Name)

	switch subtype {
	case "Image":
		return in.extractImage(name, stream)
	case "Form":
		return in.processFormXObject(name, stream)
	default:
		return nil
	}
}

// processFormXObject handles Form XObjects which contain nested content streams
func (in *Interpreter) processFormXObject(name Name, stream *Stream) error {
	// Get the Form's Resources dictionary (may inherit from parent)
	var formResources Dictionary
	if res, ok := stream.Dictionary[Name("Resources")].(Dictionary); ok {
		formResources = res
	} else if ref, ok := stream.Dictionary[Name("Resources")].(IndirectRef); ok {
		// Resolve indirect reference
		if in.FontManager != nil && in.FontManager.reader != nil {
			resolved, err := in.FontManager.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if dict, ok := resolved.(Dictionary); ok {
					formResources = dict
				}
			}
		}
	}

	// If no Form Resources, use parent Resources
	if formResources == nil {
		formResources = in.Resources
	}

	// Get the Form's Matrix (transformation matrix)
	formMatrix := IdentityMatrix()
	if matrixArr, ok := stream.Dictionary[Name("Matrix")].(Array); ok && len(matrixArr) == 6 {
		formMatrix = Matrix{
			toFloat(matrixArr[0]), toFloat(matrixArr[1]),
			toFloat(matrixArr[2]), toFloat(matrixArr[3]),
			toFloat(matrixArr[4]), toFloat(matrixArr[5]),
		}
	}

	// Decode the Form's content stream
	data, err := DecodeStreamFromDict(stream.Data, stream.Dictionary)
	if err != nil {
		fmt.Printf("Warning: Failed to decode Form XObject %s: %v\n", name, err)
		return nil // Don't fail the page for a single XObject
	}

	// Save current graphics state
	in.StateStack = append(in.StateStack, in.State)

	// Apply the Form's Matrix to CTM
	// CTM = CTM * FormMatrix
	in.State.CTM = in.State.CTM.Multiply(formMatrix)

	// Create a sub-interpreter for the Form content
	// Note: We reuse the same FontManager as it contains all font definitions
	subInterpreter := NewInterpreter(in.FontManager, formResources)

	// Copy current CTM so Form content is positioned correctly
	subInterpreter.State.CTM = in.State.CTM

	// Process the Form content
	textBlocks, images, graphics, err := subInterpreter.Process(data)
	if err != nil {
		fmt.Printf("Warning: Error processing Form XObject %s: %v\n", name, err)
		// Restore state and continue
		if len(in.StateStack) > 0 {
			in.State = in.StateStack[len(in.StateStack)-1]
			in.StateStack = in.StateStack[:len(in.StateStack)-1]
		}
		return nil
	}

	// Merge results from sub-interpreter
	in.TextBlocks = append(in.TextBlocks, textBlocks...)
	in.Images = append(in.Images, images...)
	in.Graphics = append(in.Graphics, graphics...)

	// Restore graphics state
	if len(in.StateStack) > 0 {
		in.State = in.StateStack[len(in.StateStack)-1]
		in.StateStack = in.StateStack[:len(in.StateStack)-1]
	}

	return nil
}

func (in *Interpreter) extractImage(name Name, stream *Stream) error {
	// Check if it's already JPEG (DCTDecode) - these can be saved directly
	filter := stream.Dictionary[Name("Filter")]
	isJPEG := false
	if f, ok := filter.(Name); ok && f == "DCTDecode" {
		isJPEG = true
	} else if arr, ok := filter.(Array); ok {
		for _, f := range arr {
			if fn, ok := f.(Name); ok && fn == "DCTDecode" {
				isJPEG = true
				break
			}
		}
	}

	// Decode the stream data
	data, err := DecodeStreamFromDict(stream.Data, stream.Dictionary)
	if err != nil {
		return err
	}

	// Get position from CTM
	x, y := in.State.CTM.Transform(0, 0)
	ctmWidth := in.State.CTM[0]
	ctmHeight := in.State.CTM[3]

	if isJPEG {
		// JPEG data is already properly encoded
		in.Images = append(in.Images, types.Image{
			ID:     string(name),
			Data:   data,
			Format: "jpeg",
			X:      x,
			Y:      y,
			Width:  ctmWidth,
			Height: ctmHeight,
		})
		return nil
	}

	// For non-JPEG, we need to encode raw pixel data as PNG
	// Get image dimensions from stream dictionary
	imgWidth := getImageInt(stream.Dictionary, "Width")
	imgHeight := getImageInt(stream.Dictionary, "Height")
	bpc := getImageInt(stream.Dictionary, "BitsPerComponent")
	if bpc == 0 {
		bpc = 8
	}

	if imgWidth == 0 || imgHeight == 0 {
		// Can't create image without dimensions
		return nil
	}

	// Get color space info including palette for indexed colors
	csInfo := getColorSpaceInfo(stream.Dictionary, in.FontManager)

	// Skip trivial solid-color images (decorative elements like bullets)
	if isTrivialImage(data, imgWidth, imgHeight, csInfo.Components) {
		return nil
	}

	// Create PNG from raw pixel data
	pngData, err := encodePNGWithPalette(data, imgWidth, imgHeight, csInfo, bpc)
	if err != nil {
		// If encoding fails, skip this image
		fmt.Printf("Warning: Failed to encode image %s as PNG: %v\n", name, err)
		return nil
	}

	in.Images = append(in.Images, types.Image{
		ID:     string(name),
		Data:   pngData,
		Format: "png",
		X:      x,
		Y:      y,
		Width:  ctmWidth,
		Height: ctmHeight,
	})

	return nil
}

// getImageInt gets an integer value from image dictionary
func getImageInt(dict Dictionary, key string) int {
	if v, ok := dict[Name(key)]; ok {
		switch val := v.(type) {
		case Integer:
			return int(val)
		case Real:
			return int(val)
		}
	}
	return 0
}

// ColorSpaceInfo holds color space information including palette
type ColorSpaceInfo struct {
	Components int
	IsIndexed  bool
	Palette    []color.RGBA // RGB palette for indexed colors
}

// getColorSpaceInfo extracts color space information including palette
func getColorSpaceInfo(dict Dictionary, fm *FontManager) ColorSpaceInfo {
	cs := dict[Name("ColorSpace")]
	if cs == nil {
		return ColorSpaceInfo{Components: 3}
	}

	switch v := cs.(type) {
	case Name:
		switch v {
		case "DeviceGray", "G":
			return ColorSpaceInfo{Components: 1}
		case "DeviceRGB", "RGB":
			return ColorSpaceInfo{Components: 3}
		case "DeviceCMYK", "CMYK":
			return ColorSpaceInfo{Components: 4}
		}
	case Array:
		if len(v) > 0 {
			if n, ok := v[0].(Name); ok {
				switch n {
				case "Indexed", "I":
					return parseIndexedColorSpace(v, fm)
				case "ICCBased":
					return ColorSpaceInfo{Components: 3}
				case "DeviceN", "Separation":
					return ColorSpaceInfo{Components: 1}
				}
			}
		}
	}
	return ColorSpaceInfo{Components: 3}
}

// parseIndexedColorSpace parses [/Indexed base hival lookup]
func parseIndexedColorSpace(arr Array, fm *FontManager) ColorSpaceInfo {
	info := ColorSpaceInfo{Components: 1, IsIndexed: true}

	if len(arr) < 4 {
		return info
	}

	// Get base color space components
	baseComponents := 3 // default RGB
	if baseName, ok := arr[1].(Name); ok {
		switch baseName {
		case "DeviceGray", "G":
			baseComponents = 1
		case "DeviceRGB", "RGB":
			baseComponents = 3
		case "DeviceCMYK", "CMYK":
			baseComponents = 4
		}
	}

	// Get hival (max index)
	hival := 255
	if h, ok := arr[2].(Integer); ok {
		hival = int(h)
	}

	// Get lookup table
	var lookupData []byte
	switch lookup := arr[3].(type) {
	case StringLiteral:
		lookupData = []byte(lookup)
	case HexString:
		lookupData = []byte(lookup)
	case IndirectRef:
		// Resolve the reference
		if fm != nil && fm.reader != nil {
			obj, err := fm.reader.ReadObject(lookup.ObjectNumber)
			if err == nil {
				if s, ok := obj.(Stream); ok {
					decoded, err := DecodeStreamFromDict(s.Data, s.Dictionary)
					if err == nil {
						lookupData = decoded
					}
				} else if str, ok := obj.(StringLiteral); ok {
					lookupData = []byte(str)
				}
			}
		}
	}

	// Build palette
	if len(lookupData) > 0 {
		numColors := hival + 1
		info.Palette = make([]color.RGBA, numColors)
		for i := 0; i < numColors && i*baseComponents < len(lookupData); i++ {
			idx := i * baseComponents
			switch baseComponents {
			case 1:
				g := lookupData[idx]
				info.Palette[i] = color.RGBA{R: g, G: g, B: g, A: 255}
			case 3:
				if idx+2 < len(lookupData) {
					info.Palette[i] = color.RGBA{
						R: lookupData[idx],
						G: lookupData[idx+1],
						B: lookupData[idx+2],
						A: 255,
					}
				}
			case 4:
				if idx+3 < len(lookupData) {
					c, m, y, k := lookupData[idx], lookupData[idx+1], lookupData[idx+2], lookupData[idx+3]
					r, g, b := cmykToRGB(c, m, y, k)
					info.Palette[i] = color.RGBA{R: r, G: g, B: b, A: 255}
				}
			}
		}
	}

	return info
}

// encodePNGWithPalette converts raw pixel data to PNG format with palette support
func encodePNGWithPalette(data []byte, width, height int, csInfo ColorSpaceInfo, bpc int) ([]byte, error) {
	// Handle indexed colors
	if csInfo.IsIndexed && len(csInfo.Palette) > 0 {
		return encodeIndexedPNG(data, width, height, csInfo.Palette, bpc)
	}

	// Only support 8-bit images for now
	if bpc != 8 {
		return nil, fmt.Errorf("unsupported bits per component: %d", bpc)
	}

	components := csInfo.Components
	expectedLen := width * height * components
	if len(data) < expectedLen {
		return nil, fmt.Errorf("insufficient image data: have %d, need %d", len(data), expectedLen)
	}

	var img image.Image

	switch components {
	case 1:
		// Grayscale
		gray := image.NewGray(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				if idx < len(data) {
					gray.SetGray(x, y, color.Gray{Y: data[idx]})
				}
			}
		}
		img = gray
	case 3:
		// RGB
		rgba := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := (y*width + x) * 3
				if idx+2 < len(data) {
					rgba.SetRGBA(x, y, color.RGBA{
						R: data[idx],
						G: data[idx+1],
						B: data[idx+2],
						A: 255,
					})
				}
			}
		}
		img = rgba
	case 4:
		// CMYK - convert to RGB
		rgba := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := (y*width + x) * 4
				if idx+3 < len(data) {
					c, m, yy, k := data[idx], data[idx+1], data[idx+2], data[idx+3]
					r, g, b := cmykToRGB(c, m, yy, k)
					rgba.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
				}
			}
		}
		img = rgba
	default:
		return nil, fmt.Errorf("unsupported color components: %d", components)
	}

	// Encode as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// encodeIndexedPNG handles palette-based images
func encodeIndexedPNG(data []byte, width, height int, palette []color.RGBA, bpc int) ([]byte, error) {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	switch bpc {
	case 8:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				if idx < len(data) {
					colorIdx := int(data[idx])
					if colorIdx < len(palette) {
						rgba.SetRGBA(x, y, palette[colorIdx])
					}
				}
			}
		}
	case 4:
		// 4 bits per pixel - 2 pixels per byte
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				byteIdx := (y*width + x) / 2
				if byteIdx < len(data) {
					b := data[byteIdx]
					var colorIdx int
					if x%2 == 0 {
						colorIdx = int(b >> 4)
					} else {
						colorIdx = int(b & 0x0F)
					}
					if colorIdx < len(palette) {
						rgba.SetRGBA(x, y, palette[colorIdx])
					}
				}
			}
		}
	case 2:
		// 2 bits per pixel - 4 pixels per byte
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				byteIdx := (y*width + x) / 4
				if byteIdx < len(data) {
					b := data[byteIdx]
					shift := uint(6 - 2*(x%4))
					colorIdx := int((b >> shift) & 0x03)
					if colorIdx < len(palette) {
						rgba.SetRGBA(x, y, palette[colorIdx])
					}
				}
			}
		}
	case 1:
		// 1 bit per pixel - 8 pixels per byte
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				byteIdx := (y*width + x) / 8
				if byteIdx < len(data) {
					b := data[byteIdx]
					shift := uint(7 - x%8)
					colorIdx := int((b >> shift) & 0x01)
					if colorIdx < len(palette) {
						rgba.SetRGBA(x, y, palette[colorIdx])
					}
				}
			}
		}
	default:
		return nil, fmt.Errorf("unsupported bits per component for indexed: %d", bpc)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// encodePNG converts raw pixel data to PNG format (legacy, for inline images)
func encodePNG(data []byte, width, height, components, bpc int) ([]byte, error) {
	csInfo := ColorSpaceInfo{Components: components}
	return encodePNGWithPalette(data, width, height, csInfo, bpc)
}

// isTrivialImage checks if an image is mostly a solid color (decorative element)
// Returns true if more than 95% of pixels are the same color
func isTrivialImage(data []byte, width, height, components int) bool {
	if width == 0 || height == 0 || len(data) == 0 {
		return true
	}

	// For small images, skip the check - they might be icons
	totalPixels := width * height
	if totalPixels < 100 {
		return false
	}

	// Sample pixels to determine if mostly solid color
	// Count occurrences of pixel values
	colorCounts := make(map[string]int)
	sampleRate := 1
	if totalPixels > 10000 {
		sampleRate = totalPixels / 1000 // Sample ~1000 pixels for large images
	}

	sampledPixels := 0
	for i := 0; i < totalPixels; i += sampleRate {
		idx := i * components
		if idx+components-1 >= len(data) {
			break
		}

		// Create a key from the pixel color
		var key string
		switch components {
		case 1:
			key = string(data[idx])
		case 3:
			key = string([]byte{data[idx], data[idx+1], data[idx+2]})
		case 4:
			key = string([]byte{data[idx], data[idx+1], data[idx+2], data[idx+3]})
		default:
			return false
		}

		colorCounts[key]++
		sampledPixels++
	}

	if sampledPixels == 0 {
		return true
	}

	// Find the most common color
	maxCount := 0
	for _, count := range colorCounts {
		if count > maxCount {
			maxCount = count
		}
	}

	// If more than 95% of sampled pixels are the same color, it's trivial
	ratio := float64(maxCount) / float64(sampledPixels)
	return ratio > 0.95
}

// cmykToRGB converts CMYK values to RGB
func cmykToRGB(c, m, y, k byte) (r, g, b byte) {
	// CMYK to RGB conversion
	// In PDF, CMYK values are typically 0-255 where 0 = no ink, 255 = full ink
	cf := float64(c) / 255.0
	mf := float64(m) / 255.0
	yf := float64(y) / 255.0
	kf := float64(k) / 255.0

	r = byte((1 - cf) * (1 - kf) * 255)
	g = byte((1 - mf) * (1 - kf) * 255)
	b = byte((1 - yf) * (1 - kf) * 255)
	return
}
