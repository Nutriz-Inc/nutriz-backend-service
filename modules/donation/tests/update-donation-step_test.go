package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDonationStep(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/v1/internal/donation/step"
	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin updates status", func(t *testing.T) {
			statusToUpdate := entities.EnumDonationStepStatusWarn
			setDate := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)

			body := dto.UpdateDonationStepReq{
				Description: "Status alterado para acompanhamento",
				Status:      &statusToUpdate,
				SetDate:     &setDate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfslJG",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "dst_2veL1FPpuXxUaZcFaEC57BfslJG", resp["id_donation_step"])
			assert.Equal(t, string(statusToUpdate), resp["status"])
			assert.Equal(t, body.Description, resp["description"])
		})

		t.Run("Admin updates with existing address id belonging to the user", func(t *testing.T) {
			idAddress := "adr_01JTX0H1V8N5Q3W7E2R4T6Y8ADM"
			body := dto.UpdateDonationStepReq{
				Description: "Atualizacao com endereco existente",
				IdAddress:   &idAddress,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcII",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, idAddress, resp["id_address"])
		})

		t.Run("Admin updates with new address by zip code", func(t *testing.T) {
			body := dto.UpdateDonationStepReq{
				Description: "Atualizacao com endereco criado via CEP",
				Address: &dto.AddressCreateBase{
					ZipCode: "05010000",
				},
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcJJ",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.NotEmpty(t, resp["id_address"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User does not have permission to update donation step", func(t *testing.T) {
			statusToUpdate := entities.EnumDonationStepStatusWarn
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa sem permissao",
				Status:      &statusToUpdate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update donation step", resp["message"])
		})

		t.Run("Donation step cannot be updated because it's already done or failed", func(t *testing.T) {
			statusToUpdate := entities.EnumDonationStepStatusWarn
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa de atualizar step finalizado",
				Status:      &statusToUpdate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcDD",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Donation step cannot be updated because it's already done or failed", resp["message"])
		})

		t.Run("Donation is not active", func(t *testing.T) {
			statusToUpdate := entities.EnumDonationStepStatusWarn
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa de atualizar doacao inativa",
				Status:      &statusToUpdate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcAA",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Donation is not active", resp["message"])
		})

		t.Run("You can't set a set date if the status is done or failed", func(t *testing.T) {
			statusToUpdate := entities.EnumDonationStepStatusDone
			setDate := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa invalida de status com set_date",
				Status:      &statusToUpdate,
				SetDate:     &setDate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "You can't set a set date if the status is done or failed", resp["message"])
		})

		t.Run("Set date cannot be in the past", func(t *testing.T) {
			setDate := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa com data passada",
				SetDate:     &setDate,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Set date cannot be in the past", resp["message"])
		})

		t.Run("At least one field must be sent to update", func(t *testing.T) {
			body := dto.UpdateDonationStepReq{
				Description: "Sem campos opcionais para atualizacao",
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "At least one field must be sent to update", resp["message"])
		})

		t.Run("Address does not belong to the user", func(t *testing.T) {
			idAddress := "adr_01JTG8K8N4P2R6T9V1X3Y5Z7MIP"
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa com endereco de outro usuario",
				IdAddress:   &idAddress,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Address does not belong to the user", resp["message"])
		})

		t.Run("Address not found", func(t *testing.T) {
			idAddress := "adr_2veL1FPpuXxUaZcFaEC57BfpcZZ"
			body := dto.UpdateDonationStepReq{
				Description: "Tentativa com endereco inexistente",
				IdAddress:   &idAddress,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/dst_2veL1FPpuXxUaZcFaEC57BfpcBB",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Address not found", resp["message"])
		})
	})
}
