package controller

import (
	"Oracle-Hackathon-BE/service"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ImageRepository struct {
	StorageClient service.ObjectStorageInstance
}

func NewImageRepository() *ImageRepository {
	objectStorage := service.ConnectToObjectStorage()
	return &ImageRepository{StorageClient: *objectStorage}
}

func (imageRepository *ImageRepository) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("image")

	if err != nil {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	getFile, err := file.Open()
	if err != nil {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Something wrong happened",
		})
	}

	imageRepository.StorageClient.UploadFile(file.Filename, file.Size, getFile, nil)

	fmt.Println(file)
	return nil
}
