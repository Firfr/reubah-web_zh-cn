package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dendianugerah/reubah/internal/processor/document"
	"github.com/dendianugerah/reubah/pkg/errors"
)

func ConvertDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Unable to parse form", err))
		return
	}

	// Get the uploaded file
	file, header, err := r.FormFile("document")
	if err != nil {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "No file uploaded", err))
		return
	}
	defer file.Close()

	// Check file size (32MB limit)
	if header.Size > 32<<20 {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "File size exceeds 32MB limit", nil))
		return
	}

	// Get the output format
	outputFormat := r.FormValue("format")
	if outputFormat == "" {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "No output format specified", nil))
		return
	}

	// Get input format from file extension
	inputFormat := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	
	// Validate formats
	if !document.IsFormatSupported(inputFormat, outputFormat) {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, 
			fmt.Sprintf("Conversion from %s to %s is not supported", inputFormat, outputFormat), nil))
		return
	}

	// Convert document directly from the uploaded file
	convertedContent, err := document.ConvertDocument(file, inputFormat, outputFormat)
	if err != nil {
		fmt.Printf("Document conversion error: %v\n", err)
		errors.SendError(w, errors.New(errors.ErrProcessingFailed, "Document conversion failed", err))
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", getContentType(outputFormat))
	w.Header().Set("Content-Disposition", 
		fmt.Sprintf(`attachment; filename="%s.%s"`, 
			strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)), 
			outputFormat))

	// Write the converted content
	if _, err := w.Write(convertedContent); err != nil {
		errors.SendError(w, errors.New(errors.ErrProcessingFailed, "Failed to send converted document", err))
		return
	}
}

func getContentType(format string) string {
	contentTypes := map[string]string{
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"odt":  "application/vnd.oasis.opendocument.text",
		"txt":  "text/plain",
		"rtf":  "application/rtf",
	}
	
	if ct, ok := contentTypes[format]; ok {
		return ct
	}
	return "application/octet-stream"
} 