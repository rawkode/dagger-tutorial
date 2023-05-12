package main

import (
	"context"
	"devtool"
	"fmt"
	"path/filepath"
	"runtime"

	backendBuilder "backend/build"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	openAiToken, err := devtool.GetSecret("op://Personal/OpenAI/Demo")
	if err != nil {
		return err
	}

	openAiTokenSecret := client.SetSecret("openAiToken", openAiToken)

	// Get Backend Container
	backend := backendBuilder.BuildContainerImage(client).
		WithEnvVariable("DATABASE_URI", "postgres://postgres:postgres@database:5432/postgres").
		WithSecretVariable("OPENAI_TOKEN", openAiTokenSecret).
		WithEnvVariable("PORT", "9090").
		WithExposedPort(9090)

	// Get DB
	database := devtool.GetDatabase(client)

	// Run Tests
	_, filePath, _, _ := runtime.Caller(0)
	sourceDirectory := filepath.Dir(filePath)

	src := client.Host().Directory(sourceDirectory)

	testRun := client.
		Container().
		From("ghcr.io/orange-opensource/hurl:3.0.0").
		WithServiceBinding("backend", backend).
		WithServiceBinding("database", database).
		WithDirectory("/tests", src).
		WithEntrypoint([]string{"hurl"}).
		WithExec([]string{"/tests/tests.hurl"})

	exitCode, err := testRun.ExitCode(ctx)
	if err != nil {
		return err
	}

	eOutput, err := testRun.Stderr(ctx)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		fmt.Printf("👹 Test failed with exit code %d\n", exitCode)
		fmt.Printf("Logs:\n%s\n\n", eOutput)
	} else {
		fmt.Printf("✅ PASSED")
	}

	return nil
}
