package optimize

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/dendianugerah/reubah/internal/constants"
	"github.com/dendianugerah/reubah/pkg/errors"
)

type OptimizeOptions struct {
	Quality       int  // 1-100 for JPEG/WebP
	Progressive   bool // For JPEG
	Compression   int  // 0-9 for PNG
	StripMetadata bool // Remove EXIF and other metadata
	AutoQuality   bool // Automatically determine quality based on image content
}

type QualityLevel string

const (
	QualityLow      QualityLevel = "low"      // 60% quality
	QualityMedium   QualityLevel = "medium"   // 75% quality
	QualityHigh     QualityLevel = "high"     // 90% quality
	QualityLossless QualityLevel = "lossless" // 100% quality
)

// DefaultOptions returns recommended optimization options per format
func DefaultOptions(format string) OptimizeOptions {
	switch format {
	case "jpeg", "jpg":
		return OptimizeOptions{
			Quality:       constants.DefaultQuality,
			Progressive:   true,
			StripMetadata: true,
			AutoQuality:   true,
		}
	case "png":
		return OptimizeOptions{
			Compression:   9,
			StripMetadata: true,
		}
	case "webp":
		return OptimizeOptions{
			Quality:       constants.DefaultQuality,
			StripMetadata: true,
			AutoQuality:   true,
		}
	default:
		return OptimizeOptions{
			Quality:       constants.DefaultQuality,
			StripMetadata: true,
		}
	}
}

// GetOptionsForQuality returns optimization options for a specific quality level
func GetOptionsForQuality(format string, level QualityLevel) OptimizeOptions {
	opts := DefaultOptions(format)

	switch level {
	case QualityLow:
		opts.Quality = 60
		opts.StripMetadata = true
	case QualityMedium:
		opts.Quality = 75
		opts.StripMetadata = true
	case QualityHigh:
		opts.Quality = 90
		opts.StripMetadata = false
	case QualityLossless:
		opts.Quality = 100
		opts.Compression = 9
		opts.StripMetadata = false
	}

	return opts
}

// Optimize writes the image to the writer with optimization options
func Optimize(w io.Writer, img image.Image, format string, opts OptimizeOptions) error {
	if opts.AutoQuality {
		opts = autoAdjustQuality(img, opts)
	}

	switch format {
	case "jpeg", "jpg":
		return optimizeJPEG(w, img, opts)
	case "png":
		return optimizePNG(w, img, opts)
	case "webp":
		return optimizeWebP(w, img, opts)
	default:
		return errors.New(errors.ErrOptimizationFailed, "unsupported format for optimization", nil)
	}
}

func optimizeJPEG(w io.Writer, img image.Image, opts OptimizeOptions) error {
	options := &jpeg.Options{
		Quality: opts.Quality,
	}
	if err := jpeg.Encode(w, img, options); err != nil {
		return errors.New(errors.ErrOptimizationFailed, "failed to optimize JPEG", err)
	}
	return nil
}

func optimizePNG(w io.Writer, img image.Image, opts OptimizeOptions) error {
	encoder := png.Encoder{
		CompressionLevel: png.CompressionLevel(opts.Compression),
	}
	if err := encoder.Encode(w, img); err != nil {
		return errors.New(errors.ErrOptimizationFailed, "failed to optimize PNG", err)
	}
	return nil
}

func optimizeWebP(w io.Writer, img image.Image, opts OptimizeOptions) error {
	return errors.New(errors.ErrOptimizationFailed, "WebP optimization not implemented", nil)
}

// autoAdjustQuality analyzes image content and adjusts quality settings
func autoAdjustQuality(img image.Image, opts OptimizeOptions) OptimizeOptions {
	complexity := calculateImageComplexity(img)

	if complexity < 0.3 {
		opts.Quality = 70 // Simple images
	} else if complexity < 0.6 {
		opts.Quality = 80 // Medium complexity
	} else {
		opts.Quality = 90 // Complex images
	}

	return opts
}

// calculateImageComplexity returns a value between 0 and 1
func calculateImageComplexity(img image.Image) float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return 0
	}

	sampleSize := 100
	variance := 0.0
	lastColor := img.At(0, 0)

	for i := 0; i < sampleSize; i++ {
		x := (i * width) / sampleSize
		for j := 0; j < sampleSize; j++ {
			y := (j * height) / sampleSize
			currentColor := img.At(x, y)

			r1, g1, b1, _ := lastColor.RGBA()
			r2, g2, b2, _ := currentColor.RGBA()

			diff := abs(int(r1)-int(r2)) + abs(int(g1)-int(g2)) + abs(int(b1)-int(b2))
			variance += float64(diff) / float64(0xffff*3)

			lastColor = currentColor
		}
	}

	return variance / float64(sampleSize*sampleSize)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
