package document

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math"

	"github.com/jung-kurt/gofpdf"
)

type PDFOptions struct {
	PageSize      string // "A4", "letter", "legal"
	Orientation   string // "portrait", "landscape", "auto"
	ImagesPerPage int    // 1, 2, or 4
	Quality       int    // JPEG quality for images in PDF
}

func MergeToPDF(images []image.Image, opts PDFOptions) (io.Reader, error) {
	// Set defaults if not specified
	if opts.PageSize == "" {
		opts.PageSize = "A4"
	}
	if opts.Orientation == "" || opts.Orientation == "auto" {
		opts.Orientation = determineOrientation(images[0])
	}
	if opts.ImagesPerPage == 0 {
		opts.ImagesPerPage = 1
	}
	if opts.Quality == 0 {
		opts.Quality = 85
	}

	// Create PDF
	pdf := gofpdf.New(opts.Orientation, "mm", opts.PageSize, "")
	pdf.SetMargins(10, 10, 10)

	// Get page dimensions
	pageWidth, pageHeight := pdf.GetPageSize()
	effectiveWidth := pageWidth - 20  // Account for margins
	effectiveHeight := pageHeight - 20

	// Process images
	for i := 0; i < len(images); i += opts.ImagesPerPage {
		pdf.AddPage()

		// Calculate grid layout
		cols := 1
		rows := 1
		if opts.ImagesPerPage == 2 {
			cols = 2
		} else if opts.ImagesPerPage == 4 {
			cols = 2
			rows = 2
		}

		cellWidth := effectiveWidth / float64(cols)
		cellHeight := effectiveHeight / float64(rows)

		// Process images for current page
		for j := 0; j < opts.ImagesPerPage && (i+j) < len(images); j++ {
			img := images[i+j]
			
			// Convert image to JPEG bytes
			var imgBuf bytes.Buffer
			if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: opts.Quality}); err != nil {
				return nil, fmt.Errorf("failed to encode image: %w", err)
			}

			// Calculate position in grid
			col := j % cols
			row := j / cols
			x := 10 + (float64(col) * cellWidth)
			y := 10 + (float64(row) * cellHeight)

			// Calculate image dimensions to fit cell while maintaining aspect ratio
			imgWidth := float64(img.Bounds().Dx())
			imgHeight := float64(img.Bounds().Dy())
			scale := math.Min(
				cellWidth/imgWidth,
				cellHeight/imgHeight,
			)
			finalWidth := imgWidth * scale
			finalHeight := imgHeight * scale

			// Center image in cell
			x += (cellWidth - finalWidth) / 2
			y += (cellHeight - finalHeight) / 2

			// Add image to PDF
			imgID := fmt.Sprintf("img%d", i+j)
			pdf.RegisterImageOptionsReader(imgID, gofpdf.ImageOptions{ImageType: "JPEG"}, &imgBuf)
			pdf.Image(imgID, x, y, finalWidth, finalHeight, false, "", 0, "")
		}
	}

	// Write PDF to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}

func determineOrientation(img image.Image) string {
	bounds := img.Bounds()
	if bounds.Dx() > bounds.Dy() {
		return "L" // Landscape
	}
	return "P" // Portrait
} 