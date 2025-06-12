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

// GetAllImageCollections godoc
// @Summary Get all image collections
// @Description Retrieve metadata for all uploaded image collections
// @Tags ImageCollections
// @Produce json
// @Success 200 {array} ImageCollectionInfo
// @Failure 500 {object} map[string]string
// @Router /api/image-collections [get]
func GetAllImageCollections(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bucketName := os.Getenv("MINIO_BUCKET")

	// create minio connection first
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// get obj
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	collections := []ImageCollectionInfo{}

	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error reading object:", object.Err)
			continue
		}

		// get only meta data file
		if !bytes.HasSuffix([]byte(object.Key), []byte("/metadata.json")) {
			continue
		}

		// Get the metadata object
		obj, err := minioClient.GetObject(ctx, bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			log.Println("Error getting object:", err)
			continue
		}

		var meta ImageCollectionInfo
		err = json.NewDecoder(obj).Decode(&meta)
		if err != nil {
			log.Println("Error decoding metadata:", err)
			continue
		}
		log.Printf("meta (decoded) : %s\n", meta)

		collections = append(collections, meta)
		obj.Close()
	}

	return c.JSON(collections)
}

// DeleteImageCollection godoc
// @Summary Delete an image collection
// @Description Remove all files in the specified image collection
// @Tags ImageCollections
// @Param name path string true "Collection name to delete"
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/image-collections/{name} [delete]
func DeleteImageCollection(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET")

	prefix := c.Params("name") + "/"

	// connect to minio
	minioClient, err := minioUpload.ConnectMinio()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// list object with the prefix out first
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true, // include all nested "files"
	})
	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing object:", object.Err)
			continue
		}
		err := minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Println("Error deleting object:", err)
			continue
		}
		log.Println("Deleted:", object.Key)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"msg": "Successfully removed the collection",
	})
}

// UploadImageCollection godoc
// @Summary Upload a new image collection
// @Description Upload multiple images along with metadata (name, description, date)
// @Tags ImageCollections
// @Accept multipart/form-data
// @Produce json
// @Param info formData string true "Metadata JSON for the image collection (name, description, date)"
// @Param images formData file true "Multiple image files"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/image-collections [post]
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
		"msg":      imageCollectionInfo.Name + " collection successfully uploaded",
		"uploaded": savedFiles,
	})
}
