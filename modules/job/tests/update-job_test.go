package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestUpdateJob(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/job"
	adminHeaders := &utils.TestHeadersAdmin
	nurseHeaders := &utils.TestHeadersNurse
	commonHeaders := &utils.TestHeaders

	// job_3EEMMlZS3VexkBHkGHPjq0Qt86V → pending, usado para testes de admin e erros
	// job_3EEP6m8y0peA0J6lrNABYPhhJzG → pending, reservado para testes de nurse que alteram status
	// job_2veL1FPpuXxUaZcFaEC57Bfpc01 → done, usado para testar job não-pending

	// Setup: reseta o estado do job antes dos testes, pois TestListJobs pode ter alterado seus campos
	fluxgo.RunTestRequest(
		app,
		"PUT",
		endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
		dto.UpdateJobReq{
			IdUser:      utils.StringPtr("usr_2veL1FPpuXxUaZcFaEC57BfpcNF"),
			Description: utils.StringPtr("Realizar coleta na residência da doadora"),
			DateSet:     utils.StringPtr(time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339)),
		},
		adminHeaders,
	)

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin updates description", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Nova descrição atualizada pelo adm"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "job_3EEMMlZS3VexkBHkGHPjq0Qt86V", resp["id_job"])
			assert.Equal(t, "Nova descrição atualizada pelo adm", resp["description"])
			assert.NotNil(t, resp["updated_at"])
		})

		t.Run("Admin updates id_user", func(t *testing.T) {
			body := dto.UpdateJobReq{
				IdUser: utils.StringPtr("usr_2veL1FPpuXxUaZcFaEC57BfplNV"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfplNV", resp["id_user"])
			assert.NotNil(t, resp["updated_at"])
		})

		t.Run("Admin updates date_set", func(t *testing.T) {
			body := dto.UpdateJobReq{
				DateSet: utils.StringPtr(time.Now().AddDate(0, 0, 14).Format(time.RFC3339)),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.NotNil(t, resp["date_set"])
			assert.NotNil(t, resp["updated_at"])
		})

		t.Run("Nurse updates status to done with feedback", func(t *testing.T) {
			statusDone := entities.EnumJobStatusDone
			body := dto.UpdateJobReq{
				Status:       &statusDone,
				UserFeedback: utils.StringPtr("Coleta realizada com sucesso"),
			}

			httpStatus, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEP6m8y0peA0J6lrNABYPhhJzG",
				body,
				nurseHeaders,
			)

			assert.Equal(t, http.StatusOK, httpStatus)
			assert.Equal(t, string(entities.EnumJobStatusDone), resp["status"])
			assert.Equal(t, "Coleta realizada com sucesso", resp["user_feedback"])
			assert.NotNil(t, resp["updated_at"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission to update job", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Tentativa sem permissão"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only admins and nurses can update jobs", resp["message"])
			assert.Equal(t, "job.forbidden", resp["code"])
		})

		t.Run("Job not found", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Descrição qualquer"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_2veL1FPpuXxUaZcFaEC57BfpcZZ",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Job not found", resp["message"])
		})

		t.Run("Job is not pending", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Tentativa em job finalizado"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_2veL1FPpuXxUaZcFaEC57Bfpc01",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Only pending jobs can be updated", resp["message"])
			assert.Equal(t, "job.not_pending", resp["code"])
		})

		t.Run("Admin tries to update status", func(t *testing.T) {
			statusDone := entities.EnumJobStatusDone
			body := dto.UpdateJobReq{
				Status: &statusDone,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Admins can only update id_user, description and date_set", resp["message"])
			assert.Equal(t, "job.invalid_fields_for_admin", resp["code"])
		})

		t.Run("Admin tries to update user_feedback", func(t *testing.T) {
			body := dto.UpdateJobReq{
				UserFeedback: utils.StringPtr("Tentativa indevida"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Admins can only update id_user, description and date_set", resp["message"])
			assert.Equal(t, "job.invalid_fields_for_admin", resp["code"])
		})

		t.Run("Admin assigns job to non-nurse user", func(t *testing.T) {
			body := dto.UpdateJobReq{
				IdUser: utils.StringPtr("usr_2veL1FPpuXxUaZcFaEC57BfpcKG"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Jobs can only be assigned to nurses", resp["message"])
			assert.Equal(t, "job.invalid_assignee", resp["code"])
		})

		t.Run("Admin assigns job to non-existent user", func(t *testing.T) {
			body := dto.UpdateJobReq{
				IdUser: utils.StringPtr("usr_2veL1FPpuXxUaZcFaEC57BfpcZZ"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Assignee user not found", resp["message"])
		})

		t.Run("Admin sets date_set in the past", func(t *testing.T) {
			body := dto.UpdateJobReq{
				DateSet: utils.StringPtr(time.Now().AddDate(0, 0, -1).Format(time.RFC3339)),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Set date cannot be in the past", resp["message"])
			assert.Equal(t, "job.invalid_set_date", resp["code"])
		})

		t.Run("Nurse tries to update description", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Tentativa indevida"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				nurseHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Nurses can only update status and user_feedback", resp["message"])
			assert.Equal(t, "job.invalid_fields_for_nurse", resp["code"])
		})

		t.Run("Nurse updates status without user_feedback", func(t *testing.T) {
			statusDone := entities.EnumJobStatusDone
			body := dto.UpdateJobReq{
				Status: &statusDone,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				nurseHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "User feedback is required when updating status", resp["message"])
			assert.Equal(t, "job.feedback_required", resp["code"])
		})

		t.Run("At least one field must be different to update", func(t *testing.T) {
			body := dto.UpdateJobReq{
				Description: utils.StringPtr("Nova descrição atualizada pelo adm"),
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpoint+"/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
				body,
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "At least one field must be different to update", resp["message"])
			assert.Equal(t, "job.no_fields_to_update", resp["code"])
		})
	})
}
