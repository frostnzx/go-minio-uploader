package routes

import (
	"github.com/frostnzx/antd-minio-go/controllers"
	"github.com/gofiber/fiber/v2"
)

func CsvRoute(a *fiber.App) {
	route := a.Group("/api/v1")
	route.Get("/csv/:name" , controllers.DownloadCsvFile)
	route.Post("/csv/:name" , controllers.CreateCsvFile)
}