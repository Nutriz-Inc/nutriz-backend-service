package http

import (
	"net/url"
	"strings"

	"nutriz-backend-service/config"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func GetHttp(apm *fluxgo.Apm, prom *fluxgo.Prometheus, env *config.Env) *fluxgo.Http {
	http := fluxgo.NewHttp(fluxgo.HttpOptions{
		Port:            3333,
		LogRequest:      true,
		Apm:             apm,
		Prometheus:      prom,
		AddHealthRoutes: true,
		Cors:            buildCorsConfig(env),
	})

	http.GetValidator().Validate = utils.GetValidate()

	http.CreateRouter("/public")
	http.CreateRouter("/internal", authMiddleware(env))

	return http
}

func buildCorsConfig(env *config.Env) *cors.Config {
	return &cors.Config{
		AllowOrigins:     strings.TrimSpace(env.Cors.AllowedOrigins),
		AllowOriginsFunc: isLocalhostOrigin,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,action-by",
		AllowCredentials: true,
	}
}

func isLocalhostOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	switch u.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}

const UserContextKey = "id_user"

func authMiddleware(env *config.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.ErrorUnauthorizedFiber(c)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.DecodeJwtToken(env, tokenString)
		if err != nil {
			return utils.ErrorUnauthorizedFiber(c)
		}

		isIdUserValid := utils.IdValidation(claims.IdUser)
		if !isIdUserValid {
			return utils.ErrorUnauthorizedFiber(c)
		}

		c.Locals(UserContextKey, claims.IdUser)

		c.Request().Header.Set("action-by", claims.IdUser)

		return c.Next()
	}
}
