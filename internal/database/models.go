package database

import (
	"github.com/jinzhu/gorm"
)

// HandlerType resperesents the handler enum type
type HandlerType = string

const (
	HttpHandler    HandlerType = "http"
	CronJobHandler             = "cron"
)

// RuntimeType represent the function runtime
type RuntimeType = string

const (
	NodeRuntime   RuntimeType = "node"
	GolangRuntime             = "golang"
)

// Project respesents the project data model
type Project struct {
	gorm.Model
	Name       string     `gorm:"column:name" json:"name"`
	RepoURL    string     `gorm:"column:repo_url" json:"repo_url"`
	LastCommit string     `gorm:"column:last_commit" json:"last_commit"`
	Functions  []Function `json:"-"`
}

// Function represents the function data model
type Function struct {
	gorm.Model
	Name      string      `gorm:"column:name" json:"name"`
	Image     string      `gorm:"colum:image" json:"image"`
	Runtime   RuntimeType `gorm:"column:runtime" json:"runtime"`
	Handler   HandlerType `gorm:"column:handler" json:"handler"`
	Schedule  string      `gorm:"column:schedule" json:"schedule"`
	Meta      string      `gorm:"column:meta" json:"meta"`
	ProjectID uint        `json:"-"`
}

//GetCodeTemplateFileName returns the name of the main code template file
func (f *Function) GetCodeTemplateFileName() string {
	if f.Runtime == NodeRuntime {
		if f.Handler == HttpHandler {
			return "http.js"
		}
		return "cron.js"
	}

	if f.Handler == HttpHandler {
		return "http.go"
	}
	return "cron.go"
}

// GetCodeFileName return the name of the main code file
func (f *Function) GetCodeFileName() string {
	if f.Runtime == NodeRuntime {
		return "index.js"
	}
	return "main.go"
}

// GetDockerTemplateFileName returns the dockerfile template
func (f *Function) GetDockerTemplateFileName() string {
	if f.Runtime == NodeRuntime {
		return "node.Dockerfile"
	}
	return "golang.Dockerfile"
}

//GetDockerFileName returns the name of the docker file
func (f *Function) GetDockerFileName() string {
	return "Dockerfile"
}

// Manifest represent the eserveless project configuration yaml
type Manifest struct {
	RepoURL   string `yaml:"repo"`
	Runtime   string `yaml:"runtime"`
	Functions []struct {
		Name     string            `yaml:"name"`
		Type     string            `yaml:"type"`
		Schedule string            `yaml:"schedule"`
		Meta     map[string]string `yaml:"meta"`
	} `yaml:"functions"`
}
