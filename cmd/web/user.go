package main

import (
	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/gofiber/fiber"
)

// UserController provides all user handlers
type UserController struct {
	store database.Datastore
}

// ListUsers ...
func (uc *UserController) ListUsers(c *fiber.Ctx) {

}

// CreateUser ...
func (uc *UserController) CreateUser(c *fiber.Ctx) {

}

// GetUser ...
func (uc *UserController) GetUser(c *fiber.Ctx) {

}

// DeleteUser ...
func (uc *UserController) DeleteUser(c *fiber.Ctx) {
}
