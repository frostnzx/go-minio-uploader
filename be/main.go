// @title AntD Minio API
// @version 1.0
// @description API for uploading and managing images with MinIO and AntD frontend.
// @host localhost:8080
// @BasePath /
package main

import (
	"os"

	"github.com/frostnzx/antd-minio-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/frostnzx/antd-minio-go/docs"
	"github.com/gofiber/swagger"
)

func setupFiber() *fiber.App {
	app := fiber.New()


	app.Use(cors.New())
	app.Get("/swagger/*" , swagger.HandlerDefault)

	// register routes
	routes.ImageCollectionRoute(app)


	return app
}
func main() {
	godotenv.Load()
	app := setupFiber()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
