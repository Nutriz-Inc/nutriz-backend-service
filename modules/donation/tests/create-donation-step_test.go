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

func TestCreateDonationStep(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation/step"
	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders

	makeBody := func(idDonation string, setDate *string) dto.CreateDonationStepReq {
		return dto.CreateDonationStepReq{
			IdDonation:  idDonation,
			Name:        entities.EnumDonationStepCollectMilk,
			Description: "Agendamento da coleta",
			SetDate:     setDate,
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Without set date", func(t *testing.T) {
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKG", nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotEmpty(t, resp["id_donation_step"])
			assert.Equal(t, body.IdDonation, resp["id_donation"])
			assert.Equal(t, string(body.Name), resp["name"])
			assert.Equal(t, body.Description, resp["description"])
			assert.Equal(t, string(entities.EnumDonationStepStatusPending), resp["status"])
			assert.Nil(t, resp["set_date"])
		})

		t.Run("With set date", func(t *testing.T) {
			setDate := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKH", &setDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotEmpty(t, resp["id_donation_step"])
			assert.Equal(t, body.IdDonation, resp["id_donation"])
			assert.Equal(t, string(body.Name), resp["name"])
			assert.Equal(t, body.Description, resp["description"])
			assert.Equal(t, string(entities.EnumDonationStepStatusPending), resp["status"])
			assert.NotNil(t, resp["set_date"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User does not have permission to create donation step", func(t *testing.T) {
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKG", nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create donation step", resp["message"])
		})

		t.Run("Donation not found", func(t *testing.T) {
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57Bfpd53", nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation not found", resp["message"])
		})

		t.Run("Donation is not active", func(t *testing.T) {
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKE", nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Donation is not active", resp["message"])
		})

		t.Run("Previous donation step is not completed", func(t *testing.T) {
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKF", nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Previous donation step is not completed", resp["message"])
		})

		t.Run("Set date must be in the future", func(t *testing.T) {
			setDate := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
			body := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKI", &setDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Set date must be in the future", resp["message"])
		})
	})
}
