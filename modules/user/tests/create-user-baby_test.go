package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserBaby(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/user/baby"

	headers := &utils.TestHeaders

	name := "Bebe Teste"
	t.Run("Success", func(t *testing.T) {
		data := dto.CreateUserBabyReq{
			UserBabyCreateBase: dto.UserBabyCreateBase{
				Name:      &name,
				BirthDate: "2024-01-01",
			},
		}

		status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)

		assert.Equal(t, http.StatusCreated, status)
		assert.NotNil(t, body["id_user_baby"])
		assert.NotEmpty(t, body["id_user_baby"])
		assert.Equal(t, body["id_user"], "usr_2veL1FPpuXxUaZcFaEC57BfpcKE")
		assert.Equal(t, body["name"], "Bebe Teste")
		assert.NotNil(t, body["birth_date"])
		assert.NotNil(t, body["created_at"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Birthdate in the future", func(t *testing.T) {
			data := dto.CreateUserBabyReq{
				UserBabyCreateBase: dto.UserBabyCreateBase{
					Name:      &name,
					BirthDate: "2099-01-01 00:00:00",
				},
			}
			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)
			assert.Equal(t, http.StatusBadRequest, status)
		})
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin
			data := dto.CreateUserBabyReq{
				UserBabyCreateBase: dto.UserBabyCreateBase{
					Name:      &name,
					BirthDate: "2024-01-01",
				},
			}

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, data, invalidHeader)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create baby", resp["message"])
		})
	})
}
