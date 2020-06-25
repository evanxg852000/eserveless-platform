package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"strings"

// 	"github.com/evanxg852000/eserveless/internal/data"
// 	"github.com/evanxg852000/eserveless/internal/helpers"
// 	"github.com/gofrs/uuid"
// )

// func main() {
// 	projectDir, err := ioutil.TempDir("", "repos-")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer os.RemoveAll(projectDir)

// 	hash, err := helpers.CloneProjectRepo(projectDir, "https://github.com/evanxg852000/node-eserveless-example")
// 	if err != nil {
// 		panic(err)
// 	}

// 	manifest, err := helpers.ReadProjectManifest(projectDir)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var functions []*data.Function
// 	for _, f := range manifest.Functions {
// 		meta, _ := json.Marshal(f.Meta)
// 		functions = append(functions, &data.Function{
// 			Name:      f.Name,
// 			Image:     fmt.Sprintf("%s-%s:latest", uuid.Must(uuid.NewV4()), strings.ToLower(f.Name)),
// 			Runtime:   manifest.Runtime,
// 			Handler:   f.Type,
// 			Schedule:  f.Schedule,
// 			Meta:      string(meta),
// 			ProjectID: 0,
// 		})
// 	}

// 	f := functions[0]
// 	buildDir, err := helpers.PrepareDockerImage("./runtimes", projectDir, f)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer os.RemoveAll(buildDir)

// 	logs, err := helpers.BuildDockerImage(buildDir, f)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(logs)

// 	fmt.Println("server", hash, manifest, functions)
// }
