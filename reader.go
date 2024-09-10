package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path"

	"github.com/compose-spec/compose-go/v2/cli"
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
	CurrentImage string
	LatestImage  string
}

func readFiles(root string) ([]*Server, error) {
	reader := &reader{
		root:    root,
		servers: make([]*Server, 0),
	}

	if err := reader.processRoot(); err != nil {
		return nil, err
	}

	return reader.servers, nil
}

type reader struct {
	root    string
	servers []*Server
}

func (r *reader) processRoot() error {
	entries, err := os.ReadDir(r.root)
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

	entries, err := os.ReadDir(path.Join(r.root, serverName))
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
	for _, candidate := range []string{"docker-compose.yaml", "docker-compose.yml"} {
		projectFile := path.Join(r.root, serverName, serviceName, candidate)
		exists, err := fileExists(projectFile)
		if err != nil {
			return err
		}

		if exists {
			if err := r.processFile(serverName, serviceName, projectFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *reader) processFile(serverName string, serviceName string, projectFile string) error {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{projectFile},
	)
	if err != nil {
		return err
	}

	project, err := cli.ProjectFromOptions(ctx, options)
	if err != nil {
		return err
	}

	for name, containerInfo := range project.Services {
		container := r.makeContainer(serverName, serviceName, name)
		container.CurrentImage = containerInfo.Image
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

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	var def bool
	return def, err
}
