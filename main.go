package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

// github.com/containers/image => pacman -S btrfs-progs

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repository := os.Getenv("GITHUB_REPOSITORY")

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: repository,
		Auth: &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	gitRef, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("head commit: %s\n", gitRef.Hash())

	wt, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	servers, err := readFiles(wt.Filesystem.Root(), wt.Filesystem)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	for _, server := range servers {
		for _, service := range server.Services {
			for _, container := range service.Containers {
				fullName := container.CurrentImage.FullName
				split := strings.Split(fullName, ":")

				if len(split) < 2 {
					// no proper version
					continue
				}

				baseName := split[0]

				container.LatestImage = &Image{
					FullName: baseName + ":latest",
				}

				wg.Add(1)

				go func(container *Container) {
					if err := fetchImageData(container); err != nil {
						log.Fatal(err)
					}

					wg.Done()
				}(container)
			}
		}
	}

	wg.Wait()

	for _, server := range servers {
		fmt.Printf("%s\n", server.Name)

		for _, service := range server.Services {
			fmt.Printf("  %s\n", service.Name)

			for _, container := range service.Containers {
				if container.CurrentImage != nil && container.LatestImage != nil {
					isLatest := container.CurrentImage.Digest == container.LatestImage.Digest
					fmt.Printf("    %s => %s (created %s, latest = %t)\n", container.Name, container.CurrentImage.FullName, container.CurrentImage.Created, isLatest)
					if !isLatest {
						fmt.Printf("      latest created %s\n", container.LatestImage.Created)
					}
				} else {
					fmt.Printf("    %s => %s\n", container.Name, container.CurrentImage.FullName)
				}
			}
		}
	}

}
