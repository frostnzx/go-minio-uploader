package routes

import (
	"github.com/frostnzx/antd-minio-go/controllers"
	"github.com/gofiber/fiber/v2"
)

func ImageCollectionRoute(a *fiber.App) {
	route := a.Group("/api/v1")
	route.Get("/image-collections" , controllers.GetAllImageCollections)
	route.Get("/image-collection/:name/uploads" , controllers.GetAllUploadedImages)
	route.Post("/image-collection" , controllers.UploadImageCollection)
	route.Delete("/image-collection/:name" , controllers.DeleteImageCollection)
}