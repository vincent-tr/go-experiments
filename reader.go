package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-billy/v5"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
)

type Server struct {
	Name     string
	Services []*Service
}

type Service struct {
	Name       string
	Containers []*Container
}

type Container struct {
	Name         string
	CurrentImage *Image
	LatestImage  *Image
}

func readFiles(root string, fs billy.Filesystem) ([]*Server, error) {
	reader := &reader{
		root:    root,
		fs:      fs,
		servers: make([]*Server, 0),
	}

	if err := reader.processRoot(); err != nil {
		return nil, err
	}

	return reader.servers, nil
}

type reader struct {
	root    string
	fs      billy.Filesystem
	servers []*Server
}

func (r *reader) processRoot() error {
	entries, err := r.fs.ReadDir(r.root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		serverName := entry.Name()

		err = r.processServer(serverName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *reader) processServer(serverName string) error {

	entries, err := r.fs.ReadDir(path.Join(r.root, serverName))
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		serviceName := entry.Name()

		if err := r.processService(serverName, serviceName); err != nil {
			return err
		}
	}

	return nil
}

func (r *reader) processService(serverName string, serviceName string) error {
	projectName := ""

	for _, candidate := range cli.DefaultFileNames {
		projectFile := path.Join(r.root, serverName, serviceName, candidate)
		exists, err := r.fileExists(projectFile)
		if err != nil {
			return err
		}

		if exists {
			projectName = candidate
			break
		}
	}

	if projectName != "" {
		if err := r.processComposeDir(serverName, serviceName, projectName); err != nil {
			return fmt.Errorf("error processing dir '%s': %w", path.Join(r.root, serverName, serviceName), err)
		}
	}

	return nil
}

func (r *reader) processComposeDir(serverName string, serviceName string, projectName string) error {

	projectPath := path.Join(r.root, serverName, serviceName, projectName)

	file, err := r.fs.Open(projectPath)
	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return err
	}

	fmt.Printf("read %s\n", projectPath)

	details := types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: projectPath,
				Content:  buffer.Bytes(),
			},
		},
	}

	nameOption := func(opt *loader.Options) {
		opt.SetProjectName(projectName, false)
		opt.SkipResolveEnvironment = true
	}

	project, err := loader.Load(details, nameOption)
	if err != nil {
		return err
	}

	for name, containerInfo := range project.Services {
		container := r.makeContainer(serverName, serviceName, name)
		container.CurrentImage = &Image{
			FullName: containerInfo.Image,
		}
	}

	return nil
}

func (r *reader) makeContainer(serverName string, serviceName string, containerName string) *Container {
	server := r.getServer(serverName)
	service := r.getService(server, serviceName)

	container := &Container{
		Name: containerName,
	}

	// Note: consider we won't have duplicates across project files in a service
	service.Containers = append(service.Containers, container)

	return container
}

func (r *reader) getServer(serverName string) *Server {
	for _, server := range r.servers {
		if server.Name == serverName {
			return server
		}
	}

	server := &Server{
		Name:     serverName,
		Services: make([]*Service, 0),
	}

	r.servers = append(r.servers, server)

	return server
}

func (r *reader) getService(server *Server, serviceName string) *Service {
	for _, service := range server.Services {
		if service.Name == serviceName {
			return service
		}
	}

	service := &Service{
		Name:       serviceName,
		Containers: make([]*Container, 0),
	}

	server.Services = append(server.Services, service)

	return service
}

func (r *reader) fileExists(path string) (bool, error) {
	_, err := r.fs.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	var def bool
	return def, err
}
