package controllers

import (
	"context"
	"os"

	minioUpload "github.com/frostnzx/antd-minio-go/storages"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// CreateCsvFile godoc
// @Summary Create and upload a CSV file
// @Description Creates an empty CSV file and uploads it to MinIO storage
// @Tags CSV
// @Accept json
// @Produce json
// @Param name path string true "Name of the CSV file (without extension)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/csv/{name} [post]
func CreateCsvFile(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET_CSV")
	fileName := c.Params("name") + ".csv"

	// create file first
	file, err := os.Create(fileName)
	if err != nil {
		return c.Status(500).SendString("Failed to create CSV: " + err.Error())
	}
	defer file.Close()
	defer os.Remove(fileName)
	// connect to minio
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// save to minio
	_, err = minioClient.FPutObject(ctx, bucketName, fileName, fileName, minio.PutObjectOptions{
		ContentType: "text/csv",
	})
	if err != nil {
		c.Status(500).SendString("Failed to upload CSV: " + err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"msg": fileName + "Successfully uploaded",
	})
}

// DownloadCsvFile godoc
// @Summary Download a CSV file
// @Description Downloads a CSV file from MinIO storage
// @Tags CSV
// @Produce octet-stream
// @Param name path string true "Name of the CSV file (without extension)"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/csv/{name} [get]
func DownloadCsvFile(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET_CSV")
	fileName := c.Params("name") + ".csv"

	// connect to minio
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	object, err := minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(500).SendString("Failed to get object")
	}
	stat, err := object.Stat()
	if err != nil {
		return c.Status(404).SendString("File not found")
	}
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Set("Content-Type", stat.ContentType)
	return c.Status(200).SendStream(object)
}

func DeleteCsvFile(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET_CSV")
	fileName := c.Params("name") + ".csv"

	// connect to minio
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return c.Status(500).SendString("Failed to delete object")
	}

	return c.Status(200).SendString("Delete " + fileName + " Successful")
}
