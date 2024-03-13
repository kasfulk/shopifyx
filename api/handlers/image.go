package handlers

import (
	"net/http"
	"shopifyx/db/functions"

	"github.com/gofiber/fiber/v2"
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

	file, err := fileHeader.Open()
	if err != nil {
		c.
			Status(http.StatusInternalServerError).
			JSON("failed open image")
	}

	defer file.Close()

	path, err := i.Uploader.Upload(c.UserContext(), file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(map[string]string{
		"imageUrl": path,
	})
}
