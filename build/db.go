package devtool

import (
	"dagger.io/dagger"
)

func GetDatabase(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("postgres:15").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_DB", "postgres").
		WithExposedPort(5432)
}
