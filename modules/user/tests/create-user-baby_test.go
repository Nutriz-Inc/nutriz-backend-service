package tests

import (
	httpCode "net/http"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

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
            ActionBy:   "usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
            Name:       &name,
            BirthDate:  time.Now().AddDate(-1, 0, 0),
        }
        
        status, body := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)
        t.Logf("body: %v", body)
        assert.Equal(t, httpCode.StatusCreated, status)
        assert.NotNil(t, body["id_user_baby"])
        assert.NotEmpty(t, body["id_user_baby"])
        assert.Equal(t, body["id_user"], "usr_2veL1FPpuXxUaZcFaEC57BfpcKE")
        assert.Equal(t, body["name"], "Bebe Teste")
        assert.NotNil(t, body["birth_date"])
        assert.NotNil(t, body["created_at"])
    })

    t.Run("Error", func(t *testing.T) {
        t.Run("Missing id_user", func(t *testing.T) {
            data := dto.CreateUserBabyReq{
                Name:      &name,
                BirthDate: time.Now().AddDate(-1, 0, 0),
            }
            status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, data, nil)
            assert.Equal(t, httpCode.StatusUnauthorized, status)
        })

        t.Run("Missing birth_date", func(t *testing.T) {
            data := dto.CreateUserBabyReq{
                ActionBy: "usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
                Name:   &name,
            }
            status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, data, headers)
            assert.Equal(t, httpCode.StatusBadRequest, status)
        })
    })
}