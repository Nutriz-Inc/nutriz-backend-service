package tests

import (
	"net/http"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

// Headers por perfil — autenticação via JWT (Authorization) + identidade via action-by.
var (
	// admHeaders: token ADM + action-by do ADM (usr_2veL1FPpuXxUaZcFaEC57BfpcAD)
	admHeaders = fluxgo.Headers{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoidXNyXzJ2ZUwxRlBwdVh4VWFaY0ZhRUM1N0JmcGNBRCIsImV4cCI6MTc3OTM5NzgzNCwiaWF0IjoxNzc4NzkzMDM0fQ.USyo6psNKQ2YNx0BIHUu8QmsW1KMFN5rmbzL6MhSR7I",
		"action-by":     "usr_2veL1FPpuXxUaZcFaEC57BfpcAD",
	}

	// commonHeaders: token comum (KE) + action-by do KE — para testar forbidden
	commonHeaders = fluxgo.Headers{
		"Authorization": utils.TestHeaders["Authorization"],
		"action-by":     "usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
	}
)

// IDs do seed
const (
	idNurse      = "usr_2veL1FPpuXxUaZcFaEC57BfpcNF"
	idCommonUser = "usr_2veL1FPpuXxUaZcFaEC57BfpcKE"

	// Step exclusivo para o teste de success — sem job atrelado no seed.
	// Requer inserção no seed: donation ativa KF + step CC sem job.
	idStepForSuccess = "dst_2veL1FPpuXxUaZcFaEC57BfpcCC"

	// dst_BfpcAA → donation KE → is_active = false
	idStepInactive = "dst_2veL1FPpuXxUaZcFaEC57BfpcAA"

	// IDs com formato válido mas inexistentes no banco
	idNonExistentUser = "usr_2veL1FPpuXxUaZcFaEC57BfpcZZ"
	idNonExistentStep = "dst_2veL1FPpuXxUaZcFaEC57BfpcZZ"
)

func TestCreateJob(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/job"

	baseBody := func() map[string]interface{} {
		return map[string]interface{}{
			"id_user":     idNurse,
			"id_step":     idStepForSuccess,
			"name":        "Coleta de leite materno",
			"description": "Realizar coleta na residência da doadora",
			"date_set":    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
		}
	}

	// -------------------------------------------------------------------------
	// SUCCESS
	// -------------------------------------------------------------------------

	t.Run("Success", func(t *testing.T) {
		t.Run("Create job successfully", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, baseBody(), &admHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, resp["id_job"])
			assert.Equal(t, idNurse, resp["id_user"])
			assert.Equal(t, idStepForSuccess, resp["id_step"])
			assert.Equal(t, "Coleta de leite materno", resp["name"])
			assert.NotNil(t, resp["created_at"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — FORBIDDEN
	// -------------------------------------------------------------------------

	t.Run("Error - Forbidden", func(t *testing.T) {
		t.Run("Non-admin user tries to create job", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, baseBody(), &commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only adms can create jobs", resp["message"])
			assert.Equal(t, "job.forbidden", resp["code"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — INVALID ASSIGNEE
	// -------------------------------------------------------------------------

	t.Run("Error - Invalid assignee", func(t *testing.T) {
		t.Run("Assign job to non-nurse user", func(t *testing.T) {
			body := baseBody()
			body["id_user"] = idCommonUser

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Jobs can only be assigned to nurses", resp["message"])
			assert.Equal(t, "job.invalid_assignee", resp["code"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — ASSIGNEE NOT FOUND
	// -------------------------------------------------------------------------

	t.Run("Error - Assignee not found", func(t *testing.T) {
		t.Run("Assign job to non-existent user", func(t *testing.T) {
			body := baseBody()
			body["id_user"] = idNonExistentUser

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Assignee user not found", resp["message"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — INVALID DATE
	// -------------------------------------------------------------------------

	t.Run("Error - Invalid date", func(t *testing.T) {
		t.Run("Set date in the past", func(t *testing.T) {
			body := baseBody()
			body["date_set"] = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Set date cannot be in the past", resp["message"])
			assert.Equal(t, "job.invalid_set_date", resp["code"])
		})

		t.Run("Invalid date format", func(t *testing.T) {
			body := baseBody()
			body["date_set"] = "20-05-2026" // reprovado pelo validate:"datetime=..." do DTO

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — DONATION STEP NOT FOUND
	// -------------------------------------------------------------------------

	t.Run("Error - Donation step", func(t *testing.T) {
		t.Run("Donation step not found", func(t *testing.T) {
			body := baseBody()
			body["id_step"] = idNonExistentStep

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation step not found", resp["message"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — DONATION NOT ACTIVE
	// -------------------------------------------------------------------------

	t.Run("Error - Donation not active", func(t *testing.T) {
		t.Run("Donation is inactive", func(t *testing.T) {
			body := baseBody()
			// dst_BfpcAA → donation KE → is_active = false
			body["id_step"] = idStepInactive

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Donation is not active", resp["message"])
			assert.Equal(t, "DONATION_NOT_ACTIVE", resp["code"])
		})
	})

	// -------------------------------------------------------------------------
	// ERROR — MISSING REQUIRED FIELDS
	// -------------------------------------------------------------------------

	t.Run("Error - Missing required fields", func(t *testing.T) {
		t.Run("Missing id_user", func(t *testing.T) {
			body := baseBody()
			delete(body, "id_user")

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)
			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("Missing id_step", func(t *testing.T) {
			body := baseBody()
			delete(body, "id_step")

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)
			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("Missing name", func(t *testing.T) {
			body := baseBody()
			delete(body, "name")

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)
			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("Missing description", func(t *testing.T) {
			body := baseBody()
			delete(body, "description")

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)
			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("Missing date_set", func(t *testing.T) {
			body := baseBody()
			delete(body, "date_set")

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, &admHeaders)
			assert.Equal(t, http.StatusBadRequest, status)
		})
	})
}
