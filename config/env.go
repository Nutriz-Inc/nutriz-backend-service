package config

import fluxgo "github.com/MMortari/FluxGo"

type Env struct {
	fluxgo.Env
	Database struct {
		Dsn string `env:"DATABASE_DSN" validate:"required"`
	}
	// Redis struct {
	// 	Addr string `env:"REDIS_ADDR" validate:"required"`
	// }
	Secret struct {
		IV  string `env:"SECRET_IV" validate:"required"`
		Key string `env:"SECRET_KEY" validate:"required"`
	}
	JWT struct {
		Secret string `env:"AUTH_JWT_SECRET" validate:"required"`
	}
}
