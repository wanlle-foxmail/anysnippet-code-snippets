package main

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const MaxFileSizeBytes = 1024 * 1024

var AllowedContentTypes = map[string]struct{}{
	"text/plain": {},
	"text/csv":   {},
}

func UploadValidationHandler(c echo.Context) error {
	// Flow:
	//   validate filename and content type
	//      |
	//      +-> invalid upload -> return an HTTP error response
	//      `-> read body -> validate size -> return accepted file metadata
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "file is required"})
	}

	sanitizedFilename := sanitizeFilename(fileHeader)
	if sanitizedFilename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "filename is required"})
	}

	contentType := fileHeader.Header.Get(echo.HeaderContentType)
	if _, ok := AllowedContentTypes[contentType]; !ok {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported content type"})
	}

	fileSize, err := readValidatedFileSize(fileHeader)
	if err != nil {
		return c.JSON(http.StatusRequestEntityTooLarge, map[string]string{"message": err.Error()})
	}
	if fileSize == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "file is empty"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"filename":     sanitizedFilename,
		"content_type": contentType,
		"size":         fileSize,
	})
}

func sanitizeFilename(fileHeader *multipart.FileHeader) string {
	if fileHeader == nil {
		return ""
	}

	normalizedPath := strings.ReplaceAll(fileHeader.Filename, "\\", "/")
	parts := strings.Split(normalizedPath, "/")
	return strings.TrimSpace(parts[len(parts)-1])
}

func readValidatedFileSize(fileHeader *multipart.FileHeader) (int, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buffer := make([]byte, 64*1024)
	fileSize := 0

	for {
		readCount, readErr := file.Read(buffer)
		if readCount > 0 {
			fileSize += readCount
			if fileSize > MaxFileSizeBytes {
				return 0, errors.New("file is too large")
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return 0, readErr
		}
	}

	return fileSize, nil
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.POST("/upload", UploadValidationHandler)
	return e
}

func main() {
	e := NewServer()
	e.Logger.Fatal(e.Start(":8080"))
}