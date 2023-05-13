package devtool

import (
	"fmt"

	"dagger.io/dagger"
)

type Config struct {
	DBVersion string
}

func GetDatabaseWithConfig(client *dagger.Client, config Config) *dagger.Container {
	return client.
		Container().
		From(fmt.Sprintf("postgres:%s", config.DBVersion)).
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_DB", "postgres").
		WithExposedPort(5432)
}

func GetDatabase(client *dagger.Client) *dagger.Container {
	return GetDatabaseWithConfig(client, Config{
		DBVersion: "15",
	})
}
