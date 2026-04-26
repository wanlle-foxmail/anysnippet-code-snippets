package main

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
)

const MaxFileSizeBytes = 1024 * 1024

var (
	errFileTooLarge            = errors.New("file is too large")
	errContentTypeDoesNotMatch = errors.New("file content does not match declared content type")
)

var AllowedContentTypes = map[string]struct{}{
	"text/plain": {},
	"text/csv":   {},
}

func UploadValidationHandler(c echo.Context) error {
	// Flow:
	//   validate filename and declared content type
	//      |
	//      +-> invalid upload -> return an HTTP error response
	//      `-> read body -> validate size and detected content type -> return accepted file metadata
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "file is required"})
	}

	sanitizedFilename := sanitizeFilename(fileHeader)
	if sanitizedFilename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "filename is required"})
	}

	contentType := normalizeContentType(fileHeader.Header.Get(echo.HeaderContentType))
	if _, ok := AllowedContentTypes[contentType]; !ok {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported content type"})
	}

	fileSize, err := readValidatedFileSize(fileHeader, contentType)
	if err != nil {
		if errors.Is(err, errFileTooLarge) {
			return c.JSON(http.StatusRequestEntityTooLarge, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, errContentTypeDoesNotMatch) {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
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
	basename := strings.TrimSpace(path.Base(normalizedPath))
	if basename == "" || basename == "." || basename == ".." || hasUnsafeFilenameRune(basename) {
		return ""
	}
	return basename
}

func readValidatedFileSize(fileHeader *multipart.FileHeader, declaredContentType string) (int, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buffer := make([]byte, 64*1024)
	sniffBuffer := make([]byte, 0, 512)
	fileSize := 0

	for {
		readCount, readErr := file.Read(buffer)
		if readCount > 0 {
			if len(sniffBuffer) < 512 {
				remaining := 512 - len(sniffBuffer)
				takeCount := minInt(remaining, readCount)
				sniffBuffer = append(sniffBuffer, buffer[:takeCount]...)
			}

			fileSize += readCount
			if fileSize > MaxFileSizeBytes {
				return 0, errFileTooLarge
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return 0, readErr
		}
	}

	if fileSize > 0 && !detectedContentTypeMatches(sniffBuffer, declaredContentType) {
		return 0, errContentTypeDoesNotMatch
	}

	return fileSize, nil
}

func normalizeContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil {
		return ""
	}
	return strings.ToLower(mediaType)
}

func detectedContentTypeMatches(sample []byte, declaredContentType string) bool {
	detectedContentType := normalizeContentType(http.DetectContentType(sample))
	if detectedContentType == declaredContentType {
		return true
	}
	return declaredContentType == "text/csv" && detectedContentType == "text/plain"
}

func hasUnsafeFilenameRune(filename string) bool {
	return strings.ContainsFunc(filename, func(r rune) bool {
		return r == 0 || unicode.IsControl(r)
	})
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
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
