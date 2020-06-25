package main

import (
	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/gofiber/fiber"
)

// HomeController provides home handler
type HomeController struct {
	store database.Datastore
}

// Index ...
func (hc *HomeController) Index(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"message": "Welcome eserveless platform",
		"version": "0.0.1",
		"author":  "@evanxg852000",
	})
}
