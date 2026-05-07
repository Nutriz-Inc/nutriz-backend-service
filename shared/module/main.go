package module

import (
	"nutriz-backend-service/config"
	"nutriz-backend-service/modules/donation"
	"nutriz-backend-service/shared/http"
	"nutriz-backend-service/shared/repositories"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxGo {
	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})

	flux := fluxgo.New(fluxgo.FluxGoConfig{Name: "Nutriz-Backend-Service", Version: "1", Env: &env.Env, Debugger: true, FullDebugger: true})

	flux.AddDependency(func() *config.Env { return &env })
	flux.AddDatabase(fluxgo.DatabaseOptions{Instances: []fluxgo.DatabaseConn{{Dsn: env.Database.Dsn}}})
	// flux.AddRedis(fluxgo.RedisOptions{Options: redis.Options{Addr: env.Redis.Addr}})
	flux.AddCron()
	flux.AddHttp(http.GetHttp(flux.GetApm(), nil))
	flux.AddTools()

	//Repositories
	flux.AddDependency(repositories.DonationPointRepositoryStart)
	flux.AddDependency(repositories.DonationRepositoryStart)

	//Modules
	flux.AddModule(donation.Module())

	return flux
}
