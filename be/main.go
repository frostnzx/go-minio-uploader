package main

import (
	"os"

	"github.com/frostnzx/antd-minio-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func setupFiber() *fiber.App {
	app := fiber.New()

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
