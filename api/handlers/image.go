package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"shopifyx/db/functions"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ImageUploader struct {
	Uploader *functions.ImageUploader
}

func (i *ImageUploader) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON("failed get image")
	}

	// check if file size is greater between 10kb and 2mb
	if fileHeader.Size > 2_000_000 || fileHeader.Size < 10_000 {
		return c.
			Status(http.StatusBadRequest).
			JSON("file size is too large or too small")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON("failed open image")
	}

	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(fmt.Sprintf("failed get file mimetype: %v", err.Error()))
	}

	if !(mtype.Is("image/jpeg") || mtype.Is("image/jpg")) {
		return c.
			Status(http.StatusBadRequest).
			JSON("unsupported mimetype")
	}

	filename := fmt.Sprintf("%s.%s", uuid.NewString(), filepath.Ext(fileHeader.Filename))

	path, err := i.Uploader.Upload(c.UserContext(), file, filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(map[string]string{
		"imageUrl": path,
	})
}
