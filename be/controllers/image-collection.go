package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	minioUpload "github.com/frostnzx/antd-minio-go/storages"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type ImageCollectionInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func GetAllImageCollections(c *fiber.Ctx) error {
	return nil
}
func GetImageCollection(c *fiber.Ctx) error {
	return nil
}
func UploadImageCollection(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET") // all file upload under image-collection bucket

	// decode the info about this image collection first
	imageCollectionInfoJSON := c.FormValue("info")
	var imageCollectionInfo ImageCollectionInfo
	err := json.Unmarshal([]byte(imageCollectionInfoJSON), &imageCollectionInfo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image collection info JSON",
		})
	}
	// get files of all images uploaded to this collection
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}
	files := form.File["images"]

	// save to minio
	// create minio connection first
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// start saving files
	savedFiles := []string{} // for logging
	for _, file := range files {
		// get buffer for each file
		buffer, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fileBuffer := buffer
		objectName := imageCollectionInfo.Name + "/" + file.Filename
		contentType := file.Header["Content-Type"][0]
		fileSize := file.Size

		_, err = minioClient.PutObject(
			ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		log.Printf("Successfully uploaded %s of size %d into collection %s \n", file.Filename, file.Size, imageCollectionInfo.Name)
		savedFiles = append(savedFiles, file.Filename)
		buffer.Close()
	}

	// save meta data along with the files
	metadataBytes, err := json.Marshal(imageCollectionInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal metadata",
		})
	}
	metadataReader := bytes.NewReader(metadataBytes)
	metadataObjectName := imageCollectionInfo.Name + "/metadata.json"
	_, err = minioClient.PutObject(
		ctx, bucketName, metadataObjectName, metadataReader, int64(len(metadataBytes)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload metadata: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": imageCollectionInfo.Name + " collection successfully uploaded",
		"uploaded" : savedFiles,
	})
}
