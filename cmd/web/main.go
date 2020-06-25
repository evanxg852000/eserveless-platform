package main

import (
	"github.com/evanxg852000/eserveless/internal/core"
	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/gofiber/fiber"
)

func main() {
	db, err := database.NewDB("./data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	core.StartCronHandler(db)

	hc := &HomeController{store: db}
	pc := &ProjectController{store: db}
	uc := &UserController{store: db}

	app := fiber.New()
	app.Get("/", hc.Index)
	app.All("/invoke/:project/:function", pc.InvokeFunction)

	api := app.Group("/api")
	api.Get("/projects", pc.ListProjects)
	api.Post("/projects", pc.CreateProject)
	api.Get("/projects/:name", pc.GetProject)
	api.Delete("/projects/:name", pc.DeleteProject)

	api.Get("/users", uc.ListUsers)
	api.Post("/users", uc.CreateUser)
	api.Get("/users/:name", uc.GetUser)
	api.Delete("/users/:name", uc.DeleteUser)

	app.Listen(8000)
}
