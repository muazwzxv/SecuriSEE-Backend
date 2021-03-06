package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"

	// "Oracle-Hackathon-BE/util"
	// "bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"

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

func NewImageRepository() *ImageRepository {
	db := service.GetGormInstance()
	objectStorage := service.ConnectToObjectStorage()
	return &ImageRepository{
		StorageClient: objectStorage,
		gorm:          db.Orm,
	}
}

func (r *ImageRepository) Download(ctx *fiber.Ctx) error {
	var image model.Image

	if err := image.GetById(r.gorm, ctx.Params("imageId")); err != nil {
		return NotFound(ctx, err.Error(), err)
	}

	if res, err := r.StorageClient.DownloadFile(image.FileName); err != nil {
		return Conflict(ctx, err.Error(), err)
	} else {
		content, _ := ioutil.ReadAll(res.Content)
		// download := downloadPart{
		// 	size:     int64(len(content)),
		// 	partBody: content,
		// }

		//ctx.Set("Content-Type", "image/jpg")
		//ctx.Set("Content-Type", "base64")
		//ctx.Set("Content-Type", "multipart/form-data")

		toBase64 := base64.StdEncoding.EncodeToString(content)
		return Ok(ctx, "Image successfully download", toBase64)
		//	return ctx.Status(http.StatusOK).SendStream(bytes.NewReader(content))
	}
}

func (r *ImageRepository) Upload(ctx *fiber.Ctx) error {
	// validate role
	userId := ctx.Locals("userId").(string)
	var user model.User

	// Check permissions
	user.GetUserById(r.gorm, userId)
	if !user.IsRoleUser() && !user.IsRoleAdmin() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	// Handle file
	file, err := ctx.FormFile("image")
	if err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	getFile, err := file.Open()
	if err != nil {
		return Forbidden(ctx, err.Error(), nil)
	}

	// Validate report
	var report model.Report
	if err := report.GetById(r.gorm, ctx.Params("reportId")); err != nil {
		return Conflict(ctx, "Failed to fetch associate report", nil)
	}

	// Upload to object storage
	if err := r.StorageClient.UploadFile(file.Filename, file.Size, getFile, nil); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	image := model.Image{
		FileName: file.Filename,
		ReportID: report.ID.String(),
	}

	if err := image.Create(r.gorm); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	return Ok(ctx, "Image successfully uploaded", image)
}
