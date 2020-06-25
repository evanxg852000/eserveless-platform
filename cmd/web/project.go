package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/evanxg852000/eserveless/internal/core"
	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/evanxg852000/eserveless/internal/helpers"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
)

// ProjectController provides all  project & function handlers
type ProjectController struct {
	store database.Datastore
}

// ListProjects ...
func (pc *ProjectController) ListProjects(c *fiber.Ctx) {
	//pc.store.Close()
	fmt.Println("test")
}

// CreateProject ...
func (pc *ProjectController) CreateProject(c *fiber.Ctx) {
	var data map[string]string
	err := json.Unmarshal([]byte(c.Body()), &data)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "unable to parse request data",
		})
		return
	}

	repoURL, projectName, err := helpers.ValidateGithubRepoURL(data["repository"])
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "repository field is not a valid github repository url",
		})
		return
	}

	//attempt to create project and functions
	hasChanged, isCreated, err := core.SetupProject(pc.store, projectName, repoURL)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error":   "Server Error",
			"message": err.Error(),
		})
		return
	}

	if hasChanged == false {
		c.Status(200).JSON(fiber.Map{
			"message": "repository has not changed since last deployment",
		})
		return
	}

	status := "created"
	if isCreated == false {
		status = "updated"
	}
	c.JSON(fiber.Map{
		"message":    fmt.Sprintf("Yeah! project %s.", status),
		"repository": repoURL,
		"project":    projectName,
	})
}

// GetProject ...
func (pc *ProjectController) GetProject(c *fiber.Ctx) {

}

// DeleteProject ...
func (pc *ProjectController) DeleteProject(c *fiber.Ctx) {
}

// InvokeFunction ...
func (pc *ProjectController) InvokeFunction(c *fiber.Ctx) {
	projectName := c.Params("project")
	functionName := c.Params("function")

	project := pc.store.GetProject(projectName)
	if project == nil {
		c.Status(404).JSON(fiber.Map{
			"message": "project not found!",
		})
		return
	}

	function := pc.store.GetFunction(functionName, project.ID)
	if function == nil || function.Handler != database.HttpHandler {
		c.Status(404).JSON(fiber.Map{
			"message": "function not found!",
		})
		return
	}

	//run container and redirect current request to it
	err := helpers.RunDockerImage(function, func(url string) {
		fmt.Println("\n request", url)
		time.Sleep(2 * time.Second)
		c.Redirect(url, 301)
	})
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Info("function invoked succesfully")
}
