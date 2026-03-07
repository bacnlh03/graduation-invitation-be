package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(dir string) *UploadHandler {
	return &UploadHandler{uploadDir: dir}
}

func (h *UploadHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File is required"})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Ensure upload dir exists
	if _, err := os.Stat(h.uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(h.uploadDir, os.ModePerm); err != nil {
			c.Logger().Errorf("Failed to create upload directory: %v", err)
			return err
		}
	}

	// Generate unique filename: base_timestamp.ext
	ext := filepath.Ext(file.Filename)
	base := file.Filename[:len(file.Filename)-len(ext)]
	sanitizedBase := slugify(base)
	filename := fmt.Sprintf("%s_%d%s", sanitizedBase, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(h.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		c.Logger().Errorf("Failed to create destination file: %v", err)
		return err
	}
	defer dst.Close()

	absPath, _ := filepath.Abs(dstPath)
	c.Logger().Infof("Saving file to: %s", absPath)

	if _, err = io.Copy(dst, src); err != nil {
		c.Logger().Errorf("Failed to copy file contents: %v", err)
		return err
	}

	// Return the relative URL
	fileURL := fmt.Sprintf("/api/v1/uploads/%s", filename)
	return c.JSON(http.StatusOK, map[string]string{
		"url":      fileURL,
		"filename": filename,
	})
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var res strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			res.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' || r == '.' {
			res.WriteRune('-')
		}
	}
	// Replace multiple consecutive dashes with a single one
	result := res.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
