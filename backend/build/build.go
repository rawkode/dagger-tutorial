package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
)

func ExportContainerImage(ctx context.Context, client *dagger.Client) error {
	BuildContainerImage(client).Publish(ctx, "ghcr.io/rawkode/dagger-example-backend:latest")

	return nil
}

func BuildContainerImage(client *dagger.Client) *dagger.Container {
	fmt.Println("Build Container Image")

	_, filePath, _, _ := runtime.Caller(0)
	sourceDirectory, err := filepath.Abs(fmt.Sprintf("%s/..", filepath.Dir(filePath)))
	if err != nil {
		panic(err)
	}

	src := client.Host().Directory(sourceDirectory)

	buildOutput := client.
		Container().
		From("golang:latest").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "backend", "./main.go"})

	return client.Container().
		From("ubuntu:22.10").
		WithFile("/entrypoint", buildOutput.File("/src/backend")).
		WithEntrypoint([]string{"/entrypoint"}).
		WithDefaultArgs()
}
