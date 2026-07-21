package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
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

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin", func(t *testing.T) {
			isActive := false
			quantityDonated := 4.2
			body := dto.UpdateDonationReq{
				IsActive:        &isActive,
				QuantityDonated: &quantityDonated,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/don_2veL1FPpuXxUaZcFaEC57BfpcKF",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "don_2veL1FPpuXxUaZcFaEC57BfpcKF", resp["id_donation"])
			assert.Equal(t, isActive, resp["is_active"])
			assert.Equal(t, quantityDonated, resp["quantity_donated"])
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
