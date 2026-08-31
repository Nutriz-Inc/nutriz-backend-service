package tests

import (
	"net/http"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDonation(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation"
	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders
	nurseHeaders := &utils.TestHeadersNurse
	scoreFeedback := int16(5)

	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})
	db, err := sqlx.Connect("postgres", env.Database.Dsn)
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	bottleCount := func(t *testing.T, idDonation string) int {
		var count int
		err := db.Get(&count, `SELECT COUNT(*) FROM bottle WHERE id_donation = $1`, idDonation)
		assert.NoError(t, err)
		return count
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin", func(t *testing.T) {
			idDonation := "don_2veL1FPpuXxUaZcFaEC57BfpcKF"
			isActive := false
			firstBottle := 4.2
			secondBottle := 0.0
			discarded := true
			description := "Frasco vazado"

			body := dto.UpdateDonationReq{
				IsActive: &isActive,
				Bottles: &[]dto.BottleUpdateBase{
					{
						QuantityDonatedMl: &firstBottle,
					},
					{
						QuantityDonatedMl: &secondBottle,
						Discarded:         &discarded,
						Description:       &description,
					},
				},
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/"+idDonation,
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, idDonation, resp["id_donation"])
			assert.Equal(t, isActive, resp["is_active"])

			status, getResp := fluxgo.RunTestRequest(
				app,
				"GET",
				endpoint+"/"+idDonation,
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			bottles, ok := getResp["bottles"].([]interface{})
			assert.True(t, ok)
			assert.Len(t, bottles, 2)
			assert.Equal(t, 2, bottleCount(t, idDonation))

			quantity := 10.0
			replaceBody := dto.UpdateDonationReq{
				Bottles: &[]dto.BottleUpdateBase{
					{
						QuantityDonatedMl: &quantity,
					},
				},
			}

			status, _ = fluxgo.RunTestRequest(app, "PUT", endpoint+"/"+idDonation, replaceBody, adminHeaders)
			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, 1, bottleCount(t, idDonation))
		})

		t.Run("Common", func(t *testing.T) {
			userFeedback := "Feedback atualizado"
			body := dto.UpdateDonationReq{
				UserFeedback:  &userFeedback,
				ScoreFeedback: &scoreFeedback,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKE",
				body,
				commonHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "don_2veL1FPpuXxUaZcFaEC57BfpcKE", resp["id_donation"])
			assert.Equal(t, userFeedback, resp["user_feedback"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User does not have permission to update donation", func(t *testing.T) {
			userFeedback := "Sem permissao"
			body := dto.UpdateDonationReq{
				UserFeedback:  &userFeedback,
				ScoreFeedback: &scoreFeedback,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKE",
				body,
				nurseHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update donation", resp["message"])
		})

		t.Run("Only admins can update bottles", func(t *testing.T) {
			quantity := 1.0
			body := dto.UpdateDonationReq{
				Bottles: &[]dto.BottleUpdateBase{
					{
						QuantityDonatedMl: &quantity,
					},
				},
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKE",
				body,
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "donation.bottles_forbidden", resp["code"])
		})

		t.Run("Bottle id_donation must match the donation", func(t *testing.T) {
			quantity := 1.0
			body := dto.UpdateDonationReq{
				Bottles: &[]dto.BottleUpdateBase{
					{
						QuantityDonatedMl: &quantity,
					},
				},
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKF",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "bottle.invalid_donation", resp["code"])
		})

		t.Run("Donation not found", func(t *testing.T) {
			isActive := false
			body := dto.UpdateDonationReq{
				IsActive: &isActive,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57Bfpd53",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation not found", resp["message"])
		})

		t.Run("You don't have permission to access this resource", func(t *testing.T) {
			userFeedback := "Sem permissao de acesso"
			body := dto.UpdateDonationReq{
				UserFeedback:  &userFeedback,
				ScoreFeedback: &scoreFeedback,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKF",
				body,
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", resp["message"])
		})

		t.Run("At least one field must be sent to update", func(t *testing.T) {
			body := dto.UpdateDonationReq{}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKF",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "At least one field must be sent to update", resp["message"])
		})
	})
}
