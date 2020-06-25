package helpers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/evanxg852000/eserveless/internal/database"
	"github.com/mholt/archiver/v3"

	// "github.com/docker/docker/pkg/archive"
	"github.com/go-git/go-git/v5"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// CloneProjectRepo clones the repository from a git repo url
func CloneProjectRepo(directory string, url string) (string, error) {
	rep, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return "", err
	}

	// retrieve the the branch pointed by HEAD
	ref, err := rep.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}

// ReadProjectManifest reads and decoded the project eserveless
// configursation manifest
func ReadProjectManifest(directory string) (*database.Manifest, error) {
	manifestFile := path.Join(directory, ".eserveless.yaml")
	content, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}

	manifest := database.Manifest{}
	err = yaml.Unmarshal(content, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

// PrepareDockerImage will prepare a fresh copy of the repository
// and create all necessary files from templates in order to build a docker image
func PrepareDockerImage(resDir string, repoDir string, f *database.Function) (string, error) {
	//create a temp directory to hold the repo
	buildDir, err := ioutil.TempDir("", "docker-builds-")
	if err != nil {
		return "", err
	}

	// make copy of the repo
	err = CopyDir(repoDir, buildDir)
	if err != nil {
		return "", err
	}

	//create main file
	mainCode, err := ioutil.ReadFile(path.Join(resDir, f.GetCodeTemplateFileName()))
	if err != nil {
		return "", err
	}
	mainCode = []byte(strings.ReplaceAll(string(mainCode), "{%functionName%}", f.Name))
	err = ioutil.WriteFile(
		path.Join(buildDir, f.GetCodeFileName()),
		mainCode,
		0644,
	)
	if err != nil {
		return "", err
	}

	//copy docker file
	err = CopyFile(
		path.Join(resDir, f.GetDockerTemplateFileName()),
		path.Join(buildDir, f.GetDockerFileName()),
	)
	if err != nil {
		return "", err
	}

	return buildDir, nil
}

// GetDockerBuildContext creates a tarball to Docker client SDK as io.Reader
func GetDockerBuildContext(src, dest string) (*os.File, error) {
	var files []string
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.Contains(path, ".git") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// archive format is determined by file extension
	err = archiver.Archive(files, dest)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(dest)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func BuildDockerImage(srcDir string, f *database.Function) (string, error) {
	buildCtx, err := GetDockerBuildContext(srcDir, path.Join(srcDir, "buildCtx.tar"))
	if err != nil {
		return "", err
	}
	defer buildCtx.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()
	buildOptions := types.ImageBuildOptions{
		SuppressOutput: false,
		PullParent:     true,
		Dockerfile:     "Dockerfile",
		Tags:           []string{f.Image},
		NoCache:        true,
		Remove:         true,
	}
	buildResponse, err := cli.ImageBuild(ctx, buildCtx, buildOptions)
	if err != nil {
		return "", err
	}
	defer buildResponse.Body.Close()

	reader := bufio.NewReader(buildResponse.Body)
	logs := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return logs, err
		}
		data := make(map[string]string)
		err = json.Unmarshal(line, &data)
		if err == nil {
			logs = logs + data["stream"]
		}
	}
	return logs, nil
}

func RunDockerImage(f *database.Function, onRunningCallback func(string)) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	//choose port
	port, err := freeport.GetFreePort()
	if err != nil {
		return err
	}
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(port),
	}
	containerPort, err := nat.NewPort("tcp", "8000")
	if err != nil {
		return errors.New("Unable to get the port")
	}
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	cont, err := cli.ContainerCreate(ctx, &container.Config{
		Image: f.Image,
		Tty:   true,
	}, &container.HostConfig{
		PortBindings: portBinding,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Container %s is started.", cont.ID)
	if onRunningCallback != nil {
		err = WaitHostPort(
			fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(port)),
			time.Duration(30)*time.Second,
		)
		if err != nil {
			return errors.New("Could not run the cloud function")
		}
		// call callback when container is ready
		onRunningCallback(fmt.Sprintf("http://0.0.0.0:%s", strconv.Itoa(port)))
	}

	go func() {
		//30 seconds is arbitrary it should run til func timeout
		time.Sleep(time.Duration(30) * time.Second)
		out, err := cli.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			logrus.Error(err.Error())
		}
		defer out.Close()

		buf := new(strings.Builder)
		_, err = io.Copy(buf, out)
		if err != nil && err != io.EOF {
			logrus.Error(err.Error())
		} else {
			logrus.Info(buf.String())
		}

		cli.ContainerStop(ctx, cont.ID, nil)
		if err != nil {
			logrus.Error(err.Error())
		}
	}()
	return nil
}

// CopyFile will copy a file
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// ValidateGithubRepoURL validates and extracts project name from repo Url
func ValidateGithubRepoURL(repoURL string) (string, string, error) {
	repoURL = strings.TrimSpace(repoURL)
	repoURL = strings.Trim(repoURL, "/")
	_, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(repoURL, "/")
	if len(parts) != 5 || parts[2] != "github.com" {
		return "", "", errors.New("No a valid github url")
	}

	return repoURL, fmt.Sprintf("%s-%s", parts[3], parts[4]), nil
}

// Wait for server to be available
func WaitHostPort(host string, timeout time.Duration) error {
	elapsedTime := time.Duration(0)
	for {
		conn, err := net.Dial("tcp", host)
		if err == nil && conn != nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
		elapsedTime += time.Duration(500) * time.Millisecond
		if elapsedTime >= timeout {
			return errors.New("Connection timeout")
		}
	}
}
