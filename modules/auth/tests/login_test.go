package tests

import (
	"nutriz-backend-service/modules/auth/dtos"
	"nutriz-backend-service/shared/module"
	"testing"

	httpCode "net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/public/auth/login"

	t.Run("Success", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			data := dtos.LoginReq{
				Email:    "maria@email.com",
				Password: "12345678",
			}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)

			assert.Equal(t, httpCode.StatusOK, status)
			assert.NotNil(t, body["token"])
			assert.NotEmpty(t, body["token"])
			assert.Equal(t, body["id_user"], "usr_2veL1FPpuXxUaZcFaEC57BfpcKE")
			assert.Equal(t, body["name"], "Maria Silva")
			assert.Equal(t, body["type"], "common")
		})

		t.Run("Is recurrent donor when the blood exam is valid", func(t *testing.T) {
			data := dtos.LoginReq{
				Email:    "marta@email.com",
				Password: "12345678",
			}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)

			assert.Equal(t, httpCode.StatusOK, status)
			assert.Equal(t, body["id_user"], "usr_2veL1FPpuXxUaZcFaEC57BfpcKL")
			assert.Equal(t, true, body["is_recurrent_donor"])
		})

		t.Run("Is not recurrent donor when the blood exam is expired", func(t *testing.T) {
			data := dtos.LoginReq{
				Email:    "beatriz@email.com",
				Password: "12345678",
			}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)

			assert.Equal(t, httpCode.StatusOK, status)
			assert.Equal(t, body["id_user"], "usr_2veL1FPpuXxUaZcFaEC57BfpcBX")
			assert.Equal(t, false, body["is_recurrent_donor"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User not found", func(t *testing.T) {
			data := dtos.LoginReq{
				Email:    "nonexistent@example.com",
				Password: "password123",
			}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)

			assert.Equal(t, httpCode.StatusBadRequest, status)
			assert.Equal(t, "Invalid email or password", body["message"])
		})
		t.Run("Invalid password", func(t *testing.T) {
			data := dtos.LoginReq{
				Email:    "maria@email.com",
				Password: "123sfdgtjrskr65rrst",
			}

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)

			assert.Equal(t, httpCode.StatusBadRequest, status)
			assert.Equal(t, "Invalid email or password", body["message"])
		})
	})
}
