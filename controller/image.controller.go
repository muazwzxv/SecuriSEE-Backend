package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"
	"Oracle-Hackathon-BE/util"
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ImageRepository struct {
	StorageClient *service.ObjectStorageInstance
	gorm          *gorm.DB
	File          *multipart.FileHeader
}
type downloadPart struct {
	size     int64
	partBody []byte
	// offset   int64
	// partNum  int
	// err      error
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	objectStorage := service.ConnectToObjectStorage()
	return &ImageRepository{
		StorageClient: objectStorage,
		gorm:          db,
	}
}

func (imageRepository *ImageRepository) Download(ctx *fiber.Ctx) error {
	var image model.Image

	if err := image.GetById(imageRepository.gorm, ctx.Params("imageId")); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	if res, err := imageRepository.StorageClient.DownloadFile(image.FileName); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	} else {
		content, _ := ioutil.ReadAll(res.Content)
		// download := downloadPart{
		// 	size:     int64(len(content)),
		// 	partBody: content,
		// }

		ctx.Set("Content-Type", "multipart/form-data")
		//return ctx.SendStream(bytes.NewReader(content))
		return ctx.Status(http.StatusOK).SendStream(bytes.NewReader(content))
	}
}

func (imageRepository *ImageRepository) Upload(ctx *fiber.Ctx) error {
	// validate role
	claim := util.GetClaims(ctx)
	var user model.User
	user.GetUserById(imageRepository.gorm, claim["ID"].(string))

	// Check permissions
	isUser := user.IsRoleExist("user")
	isAdmin := user.IsRoleExist("admin")

	if !isUser && !isAdmin {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	// Handle file
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

	// Validate report
	var report model.Report
	if err := report.GetById(imageRepository.gorm, ctx.Params("reportId")); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": "Failed to fetch associate report",
		})
	}

	// Upload to object storage
	if err := imageRepository.StorageClient.UploadFile(file.Filename, file.Size, getFile, nil); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	image := model.Image{
		FileName: file.Filename,
		ReportID: report.ID,
	}

	if err := image.Create(imageRepository.gorm); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Success": true,
		"Message": "Uploaded",
	})

}
