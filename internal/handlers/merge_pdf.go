package handlers

import (
	"bytes"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dendianugerah/reubah/internal/processor/document"
	"github.com/dendianugerah/reubah/internal/validator"
	"github.com/dendianugerah/reubah/pkg/errors"
)

// seekableReader wraps a *bytes.Reader to implement io.ReadSeeker
type seekableReader struct {
	*bytes.Reader
}

func MergePDF(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Unable to parse form", err))
		return
	}

	// Get all uploaded files
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "No files uploaded", nil))
		return
	}

	// Process images
	images, err := processUploadedFiles(files)
	if err != nil {
		errors.SendError(w, err)
		return
	}

	// Get PDF options from form
	opts := document.PDFOptions{
		PageSize:      r.FormValue("pageSize"),
		Orientation:   r.FormValue("orientation"),
		ImagesPerPage: getImagesPerPage(r.FormValue("imagesPerPage")),
		Quality:       85, // Default quality
	}

	// Generate PDF
	pdfReader, err := document.MergeToPDF(images, opts)
	if err != nil {
		errors.SendError(w, errors.New(errors.ErrProcessingFailed, "Failed to generate PDF", err))
		return
	}

	// Convert io.Reader to []byte
	pdfBytes, err := io.ReadAll(pdfReader)
	if err != nil {
		errors.SendError(w, errors.New(errors.ErrProcessingFailed, "Failed to read PDF", err))
		return
	}

	// Create seekable reader from bytes
	seeker := &seekableReader{bytes.NewReader(pdfBytes)}

	// Send response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=merged.pdf")
	http.ServeContent(w, r, "merged.pdf", time.Now(), seeker)
}

func processUploadedFiles(files []*multipart.FileHeader) ([]image.Image, error) {
	images := make([]image.Image, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, errors.New(errors.ErrInvalidFormat, "Failed to open file", err)
		}
		defer file.Close()

		if err := validator.ValidateMIMEType(file); err != nil {
			return nil, errors.New(errors.ErrInvalidMIME, "Invalid file type", err)
		}

		// Reset file pointer after MIME check
		if _, err := file.Seek(0, 0); err != nil {
			return nil, errors.New(errors.ErrInvalidFormat, "Failed to process file", err)
		}

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, errors.New(errors.ErrInvalidFormat, "Invalid image file", err)
		}

		images = append(images, img)
	}

	return images, nil
}

func getImagesPerPage(value string) int {
	if value == "" {
		return 1
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 || n > 4 {
		return 1
	}
	return n
} 