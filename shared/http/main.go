package http

import (
	"nutriz-backend-service/config"
	"nutriz-backend-service/shared/utils"
	"strings"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

func GetHttp(apm *fluxgo.Apm, prom *fluxgo.Prometheus, env *config.Env) *fluxgo.Http {
	http := fluxgo.NewHttp(fluxgo.HttpOptions{Port: 3333, LogRequest: true, Apm: apm, Prometheus: prom, AddHealthRoutes: true})

	http.CreateRouter("/public")
	http.CreateRouter("/internal", authMiddleware(env))

	return http
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

		return c.Next()
	}
}
