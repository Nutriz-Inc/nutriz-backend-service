package http

import (
	"nutriz-backend-service/config"
	"nutriz-backend-service/docs"
	"nutriz-backend-service/shared/utils"
	"strings"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

func GetHttp(apm *fluxgo.Apm, prom *fluxgo.Prometheus, env *config.Env) *fluxgo.Http {
	http := fluxgo.NewHttp(fluxgo.HttpOptions{
		Port:            3333,
		LogRequest:      true,
		Apm:             apm,
		Prometheus:      prom,
		AddHealthRoutes: true,
	})

	http.GetValidator().Validate = utils.GetValidate()

	http.CreateRouter("/public")
	http.CreateRouter("/internal", authMiddleware(env))

	registerDocsRoutes(http.GetApp())

	return http
}

// registerDocsRoutes exposes the hand-written OpenAPI spec (docs/openapi.json)
// as raw JSON plus an interactive Swagger UI, so endpoints can be documented
// and tested without going through the module/route income system.
func registerDocsRoutes(app *fiber.App) {
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.Send(docs.OpenAPISpec)
	})

	app.Get("/docs/*", swagger.New(swagger.Config{
		URL:   "/openapi.json",
		Title: "Nutriz API Docs",
	}))
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
