package tests

import (
	"fmt"
	"net/http"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserBaby(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/v1/internal/user/baby"

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
		t.Run("User can have up to %d babies", func(t *testing.T) {
			route := "/v1/internal/user/usr_2veL1FPpuXxUaZcFaEC57BfpcKE?show_baby=true"
			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)
			assert.Equal(t, http.StatusOK, status)

			babies, ok := body["babies"].([]interface{})
			assert.True(t, ok)

			for i := len(babies); i < entities.MAX_BABY_QUANTITY_PER_USER; i++ {
				babyName := fmt.Sprintf("Baby %d", i+1)
				data := dto.CreateUserBabyReq{
					UserBabyCreateBase: dto.UserBabyCreateBase{
						Name:      &babyName,
						BirthDate: "2024-01-01",
					},
				}
				status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)
				assert.Equal(t, http.StatusCreated, status)
			}

			extraBabyName := fmt.Sprintf("Baby %d", entities.MAX_BABY_QUANTITY_PER_USER+1)
			data := dto.CreateUserBabyReq{
				UserBabyCreateBase: dto.UserBabyCreateBase{
					Name:      &extraBabyName,
					BirthDate: "2024-01-01",
				},
			}

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, fmt.Sprintf("User can have up to %d babies", entities.MAX_BABY_QUANTITY_PER_USER), resp["message"])
		})
	})
}
