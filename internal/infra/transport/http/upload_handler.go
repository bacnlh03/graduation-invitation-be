package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		os.MkdirAll(h.uploadDir, os.ModePerm)
	}

	// Generate unique filename: base_timestamp.ext
	ext := filepath.Ext(file.Filename)
	base := file.Filename[:len(file.Filename)-len(ext)]
	filename := fmt.Sprintf("%s_%d%s", base, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(h.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Return the relative URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	return c.JSON(http.StatusOK, map[string]string{
		"url":      fileURL,
		"filename": filename,
	})
}
