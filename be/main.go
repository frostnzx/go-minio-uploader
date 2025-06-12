package main

import (
	"os"

	"github.com/frostnzx/antd-minio-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupFiber() *fiber.App {
	app := fiber.New()


	app.Use(cors.New())

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
