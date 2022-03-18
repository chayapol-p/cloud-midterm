package main

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/chayapol-p/cloud-midterm-server/controllers"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload"
)

func hello(c *fiber.Ctx) error {
	return c.SendString("Welcome Super Cool Midterm Project!!")
}

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault) // get one user by ID
}

func MessagesRoute(app *fiber.App) {
	route := app.Group("/api/messages")

	route.Get("/", controllers.GetMessages)
	route.Post("/", controllers.CreateMessage)
	route.Delete("/:uuid", controllers.DeleteMessage)
	route.Put("/:uuid", controllers.UpdateMessage)
}

func main() {
	app := fiber.New()
	app.Get("/", hello)
	MessagesRoute(app)
	SwaggerRoute(app)
	app.Listen(":3000")
}
