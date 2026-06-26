package main

import (
	c "context"
	"nutriz-backend-service/config"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func main() {
	ctx, cancel := c.WithTimeout(c.Background(), time.Minute)
	defer cancel()

	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})

	flux := fluxgo.
		New(fluxgo.FluxGoConfig{Name: "Migrations"}).
		AddDatabase(fluxgo.DatabaseOptions{Instances: []fluxgo.DatabaseConn{{Dsn: env.Database.Dsn}}})

	if err := flux.RunMigrations(ctx, fluxgo.DatabaseMigrationsOptions{Dir: "shared/database/migrations", Seeds: &seed}); err != nil {
		panic(err)
	}
}
