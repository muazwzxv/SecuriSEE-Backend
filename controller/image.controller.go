package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"
	"Oracle-Hackathon-BE/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ImageRepository struct {
	StorageClient service.ObjectStorageInstance
	gorm          *gorm.DB
}

func NewImageRepository() *ImageRepository {
	objectStorage := service.ConnectToObjectStorage()
	return &ImageRepository{StorageClient: *objectStorage}
}

func (imageRepository *ImageRepository) Upload(ctx *fiber.Ctx) error {

	// validate role
	claim := util.GetClaims(ctx)
	var user model.User
	user.GetUserById(imageRepository.gorm, claim["ID"].(string))

	// Check permissions
	isUser := user.IsRoleExist("user")
	isAdmin := user.IsRoleExist("admin")

	if !isUser || !isAdmin {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	file, err := ctx.FormFile("image")

	if err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	getFile, err := file.Open()
	if err != nil {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	if err := imageRepository.StorageClient.UploadFile(file.Filename, file.Size, getFile, nil); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	image := model.Image{
		FileName: file.Filename,
		UserID:   user.ID,
	}

	_ = image.Create(imageRepository.gorm)

	return ctx.Status(http.StatusConflict).JSON(fiber.Map{
		"Success": false,
		"Message": err.Error(),
	})

}
