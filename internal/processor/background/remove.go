package background

import (
	"bytes"
	"image"
	"image/png"
	"os/exec"
)

func RemoveBackground(img image.Image) (image.Image, error) {
	// Convert input image to PNG bytes
	var inputBuf bytes.Buffer
	if err := png.Encode(&inputBuf, img); err != nil {
		return nil, err
	}

	// Create rembg command
	cmd := exec.Command("rembg", "i", "-", "-") // "-" means use stdin/stdout
	cmd.Stdin = &inputBuf
	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, err
	}	

	// Decode the result back to image.Image
	result, err := png.Decode(&outputBuf)
	if err != nil {
		return nil, err
	}

	return result, nil
}
