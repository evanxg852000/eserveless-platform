package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/evanxg852000/eserveless/internal/helpers"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func SetupProject(store database.Datastore, projectName, repoURL string) (bool, bool, error) {
	isCreated := false
	hasChanged := false
	projectDir, err := ioutil.TempDir("", "repos-")
	if err != nil {
		return hasChanged, isCreated, err
	}

	currentCommit, err := helpers.CloneProjectRepo(projectDir, repoURL)
	if err != nil {
		return hasChanged, isCreated, err
	}

	project := store.GetProject(projectName)
	if project != nil && project.LastCommit == currentCommit {
		return hasChanged, isCreated, nil
	}

	hasChanged = true
	if project == nil {
		isCreated = true
		project = &database.Project{
			Name:       projectName,
			RepoURL:    repoURL,
			LastCommit: currentCommit,
		}
	}

	manifest, err := helpers.ReadProjectManifest(projectDir)
	if err != nil {
		return hasChanged, isCreated, err
	}

	//create project & functions in store
	project.LastCommit = currentCommit
	store.SaveProject(project)

	var functions []*database.Function
	for _, f := range manifest.Functions {
		meta, _ := json.Marshal(f.Meta)
		imageName := fmt.Sprintf("%s-%s-%s:latest", currentCommit[0:7], projectName, strings.ToLower(f.Name))
		fn := &database.Function{
			Name:      f.Name,
			Image:     imageName,
			Runtime:   manifest.Runtime,
			Handler:   f.Type,
			Schedule:  f.Schedule,
			Meta:      string(meta),
			ProjectID: project.ID,
		}
		if foundFn := store.GetFunction(f.Name, project.ID); foundFn != nil {
			fn.ID = foundFn.ID
		}
		store.SaveFunction(fn)
		functions = append(functions, fn)
	}

	go func(projectDir string, store database.Datastore, functions []*database.Function) {
		defer os.RemoveAll(projectDir) // remove temp folder after all
		//build images & store logs if any error
		for _, f := range functions {
			logrus.Info(fmt.Sprintf("preparing project & build function image: %s", f.Image))
			buildDir, err := helpers.PrepareDockerImage("./runtimes", projectDir, f)
			if err != nil {
				logrus.Error(err.Error())
				break
			}

			logs, err := helpers.BuildDockerImage(buildDir, f)
			if err != nil {
				logrus.Error(err.Error())
			}
			logrus.Info(logs)
			os.RemoveAll(buildDir)
		}
	}(projectDir, store, functions)

	return hasChanged, isCreated, nil
}

func StartCronHandler(store database.Datastore) {
	cronExecutor := cron.New()
	makeCallback := func(fn *database.Function) func() {
		return func() {
			go func() {
				//run container
				err := helpers.RunDockerImage(fn, nil)
				if err != nil {
					logrus.Error(err.Error())
					return
				}
				logrus.Info(fmt.Sprintf("function invoked succesfully: %s", fn.Name))
			}()
		}
	}

	functions := store.GetCronFunctions()
	for _, fn := range functions {
		cronExecutor.AddFunc(fn.Schedule, makeCallback(fn))
	}
	cronExecutor.Start()
}
