package tests

import (
	"fmt"
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
	endpoint := "/internal/donation/step"
	donationEndpoint := "/internal/donation"
	headers := &utils.TestHeaders

	setupApp := func(t *testing.T) (*fluxgo.FluxGo, interface{}) {
		fx, app := module.Module().GetTestApp(t)
		return fx, app
	}

	createDonation := func(t *testing.T, app interface{}) string {
		status, body := fluxgo.RunTestRequest(app, "POST", donationEndpoint, nil, headers)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, body["id_donation"])

		return fmt.Sprintf("%v", body["id_donation"])
	}

	makeBody := func(idDonation string, setDate *string) dto.CreateDonationStepReq {
		return dto.CreateDonationStepReq{
			IdDonation:  idDonation,
			Name:        entities.EnumDonationStepBloodTest,
			Description: "Agendamento da proxima etapa",
			SetDate:     setDate,
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Without set_date", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			idDonation := createDonation(t, app)
			req := makeBody(idDonation, nil)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.NotEmpty(t, body["id_donation_step"])
			assert.Equal(t, idDonation, body["id_donation"])
			assert.Equal(t, string(entities.EnumDonationStepBloodTest), body["name"])
			assert.Equal(t, "pending", body["status"])
			assert.Nil(t, body["set_date"])
		})

		t.Run("With set_date", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			idDonation := createDonation(t, app)
			setDate := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
			req := makeBody(idDonation, &setDate)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.NotEmpty(t, body["id_donation_step"])
			assert.Equal(t, idDonation, body["id_donation"])
			assert.Equal(t, string(entities.EnumDonationStepBloodTest), body["name"])
			assert.Equal(t, "pending", body["status"])
			assert.NotNil(t, body["set_date"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Donation not found", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			req := makeBody("don_2veL1FPpuXxUaZcFaEC57Bfpd53", nil)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation not found", body["message"])
		})

		t.Run("Donation is not active", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			req := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKE", nil)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Donation is not active", body["message"])
		})

		t.Run("Previous donation step is not completed", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			req := makeBody("don_2veL1FPpuXxUaZcFaEC57BfpcKF", nil)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Previous donation step is not completed", body["message"])
		})

		t.Run("Set date must be in the future", func(t *testing.T) {
			fx, app := setupApp(t)
			defer fx.RequireStart().RequireStop()

			idDonation := createDonation(t, app)
			pastDate := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
			req := makeBody(idDonation, &pastDate)

			status, body := fluxgo.RunTestRequest(app, "POST", endpoint, req, headers)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Set date must be in the future", body["message"])
		})
	})
}
