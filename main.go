package main

import (
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	root, err := filepath.Abs("../lan-scripts")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("root = %s\n", root)

	servers, err := readFiles(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		fmt.Printf("%s\n", server.Name)

		for _, service := range server.Services {
			fmt.Printf("  %s\n", service.Name)

			for _, container := range service.Containers {
				fmt.Printf("    %s => %s\n", container.Name, container.CurrentImage)

			}
		}
	}
}
