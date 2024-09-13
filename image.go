package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/types"
)

type Image struct {
	FullName string
	Digest   string
	Created  *time.Time
}

func fetchImageData(container *Container) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	sys := &types.SystemContext{}

	images := []*Image{container.CurrentImage, container.LatestImage}

	wg := sync.WaitGroup{}
	wg.Add(len(images))

	for _, image := range images {
		go func(image *Image) {

			defer wg.Done()

			imageRef, err := docker.ParseReference("//" + image.FullName)
			if err != nil {
				// TODO
				fmt.Println(err)
				return
			}

			rawDigest, err := docker.GetDigest(ctx, sys, imageRef)
			if err != nil {
				// TODO
				fmt.Println(err)
				return
			}

			image.Digest = string(rawDigest)

			created, err := findImageCreated(ctx, sys, imageRef)
			if err != nil {
				// TODO
				fmt.Println(err)
				return
			}

			image.Created = created

		}(image)
	}

	wg.Wait()

	return nil
}

func findImageCreated(ctx context.Context, sys *types.SystemContext, imageRef types.ImageReference) (*time.Time, error) {
	imageSrc, err := imageRef.NewImageSource(ctx, sys)
	if err != nil {
		return nil, err
	}

	img, err := image.FromUnparsedImage(ctx, sys, image.UnparsedInstance(imageSrc, nil))
	if err != nil {
		return nil, err
	}

	info, err := img.Inspect(ctx)
	if err != nil {
		return nil, err
	}

	return info.Created, nil
}
