package module

import (
	"nutriz-backend-service/config"
	"nutriz-backend-service/modules/auth"
	"nutriz-backend-service/modules/donation"
	"nutriz-backend-service/shared/http"
	"nutriz-backend-service/shared/repositories"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

func Module() *fluxgo.FluxGo {
	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})

	flux := fluxgo.New(fluxgo.FluxGoConfig{Name: "Nutriz-Backend-Service", Version: "1", Env: &env.Env, Debugger: true, FullDebugger: true})
	flux.AddApm(fluxgo.ApmOptions{CollectorURL: env.Apm.CollectorUrl, Exporter: env.Apm.Exporter})
	flux.ConfigLogger(fluxgo.LoggerOptions{Type: env.Logger.Type, Level: env.Logger.Level, LogFilePath: env.Logger.FilePath})

	prom := flux.AddPrometheus()
	prom.NewCounterVec(prometheus.CounterOpts{
		Name: "get_user",
		Help: "Quantidade de requests para buscar usuários",
	}, []string{"user"})

	flux.AddDependency(func() *config.Env { return &env })
	flux.AddDatabase(fluxgo.DatabaseOptions{Instances: []fluxgo.DatabaseConn{{Dsn: env.Database.Dsn}}})
	flux.AddRedis(fluxgo.RedisOptions{Options: redis.Options{Addr: env.Redis.Addr}})
	flux.AddKafka(env.Kafka.GetConfig())
	flux.AddCron()
	flux.AddHttp(http.GetHttp(flux.GetApm(), prom))
	flux.AddTools()

	//Repositories
	flux.AddDependency(repositories.DonationPointRepositoryStart)
	flux.AddDependency(repositories.DonationRepositoryStart)
	flux.AddDependency(repositories.UserRepositoryStart)

	//Modules
	flux.AddModule(donation.Module())
	flux.AddModule(auth.Module())

	return flux
}
