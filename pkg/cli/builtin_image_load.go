package cli

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/duso-org/duso/pkg/script"
)

// builtinLoadImage calls load_binary then adds image metadata
func builtinLoadImage(evaluator *script.Evaluator, args map[string]any) (any, error) {
	// Call load_binary to get the binary
	binValue, err := builtinLoadBinary(evaluator, args)
	if err != nil {
		return nil, err
	}

	// Extract the binary from the result
	var bin *script.BinaryValue
	if val, ok := binValue.(script.Value); ok && val.IsBinary() {
		bin = val.AsBinary()
	} else {
		return binValue, nil // Not a binary, return as-is
	}

	// Try to decode the image to extract metadata
	img, format, err := decodeImageData(*bin.Data)
	if err != nil {
		// If decoding fails, just return the binary without image metadata
		return binValue, nil
	}

	// Populate image metadata
	bounds := img.Bounds()
	bin.Metadata["width"] = script.NewNumber(float64(bounds.Dx()))
	bin.Metadata["height"] = script.NewNumber(float64(bounds.Dy()))
	bin.Metadata["format"] = script.NewString(format)

	contentType := "image/png"
	if format == "jpeg" {
		contentType = "image/jpeg"
	} else if format == "gif" {
		contentType = "image/gif"
	}
	bin.Metadata["content_type"] = script.NewString(contentType)

	return binValue, nil
}

// builtinSaveImage just delegates to save_binary
func builtinSaveImage(evaluator *script.Evaluator, args map[string]any) (any, error) {
	return builtinSaveBinary(evaluator, args)
}

// decodeImageData decodes binary data and returns the image and format
func decodeImageData(data []byte) (image.Image, string, error) {
	// Try PNG first
	if img, err := png.Decode(bytes.NewReader(data)); err == nil {
		return img, "png", nil
	}

	// Try JPEG
	if img, err := jpeg.Decode(bytes.NewReader(data)); err == nil {
		return img, "jpeg", nil
	}

	// Try GIF
	if img, err := gif.Decode(bytes.NewReader(data)); err == nil {
		return img, "gif", nil
	}

	return nil, "", fmt.Errorf("could not decode image: unsupported format")
}
