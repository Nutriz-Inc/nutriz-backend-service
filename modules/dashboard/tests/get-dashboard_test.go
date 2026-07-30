package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestGetDashboard(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/dashboard"
	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin gets full dashboard", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)

			totalMilkCollected, ok := body["total_milk_collected"].(float64)
			assert.True(t, ok, "total_milk_collected should be a number")
			assert.GreaterOrEqual(t, totalMilkCollected, 3.9, "seeded donations alone already sum to 3.9L")

			milkByMonth := fluxgo.ConvertToList(body["milk_collected_by_month"])
			assert.GreaterOrEqual(t, len(milkByMonth), 1)

			var sumByMonth float64
			for _, item := range milkByMonth {
				row := fluxgo.ConvertToMap(item)
				assert.NotEmpty(t, row["month"])
				total, ok := row["total"].(float64)
				assert.True(t, ok)
				sumByMonth += total
			}
			assert.InDelta(t, totalMilkCollected, sumByMonth, 0.001, "monthly breakdown should add up to the total")

			feedback := fluxgo.ConvertToList(body["feedback_by_score"])
			assert.Len(t, feedback, 5, "feedback breakdown should always report all 5 possible scores")
			for i, item := range feedback {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, float64(i+1), row["score"])
				count, ok := row["count"].(float64)
				assert.True(t, ok)
				assert.GreaterOrEqual(t, count, float64(0))
			}

			donationsWithError, ok := body["donations_with_error"].(float64)
			assert.True(t, ok)
			assert.GreaterOrEqual(t, donationsWithError, float64(0))

			recurrenceRate, ok := body["donor_recurrence_rate"].(float64)
			assert.True(t, ok)
			assert.GreaterOrEqual(t, recurrenceRate, float64(0))
			assert.LessOrEqual(t, recurrenceRate, float64(100))

			if body["average_service_time_hours"] != nil {
				_, ok := body["average_service_time_hours"].(float64)
				assert.True(t, ok, "average_service_time_hours should be a number when present")
			}

			activeDonationsByStep := fluxgo.ConvertToList(body["active_donations_by_step"])
			assert.Len(t, activeDonationsByStep, 4, "should always report all 4 possible donation steps")
			var sumPercentage float64
			for _, item := range activeDonationsByStep {
				row := fluxgo.ConvertToMap(item)
				assert.NotEmpty(t, row["step"])
				count, ok := row["count"].(float64)
				assert.True(t, ok)
				assert.GreaterOrEqual(t, count, float64(0))

				percentage, ok := row["percentage"].(float64)
				assert.True(t, ok)
				assert.GreaterOrEqual(t, percentage, float64(0))
				assert.LessOrEqual(t, percentage, float64(100))
				sumPercentage += percentage
			}
			assert.LessOrEqual(t, sumPercentage, float64(100.01), "percentages across the 4 steps should never exceed 100%")
		})

		t.Run("Date range with no data returns empty metrics", func(t *testing.T) {
			futureStart := utils.GetTodayFormattedDate(365)
			futureEnd := utils.GetTodayFormattedDate(366)
			route := fmt.Sprintf("%s?start_date=%s&end_date=%s", endpoint, futureStart, futureEnd)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, float64(0), body["total_milk_collected"])
			assert.Equal(t, float64(0), body["donations_with_error"])
			assert.Equal(t, float64(0), body["donor_recurrence_rate"])
			assert.Nil(t, body["average_service_time_hours"])

			milkByMonth := fluxgo.ConvertToList(body["milk_collected_by_month"])
			assert.Len(t, milkByMonth, 0)

			feedback := fluxgo.ConvertToList(body["feedback_by_score"])
			assert.Len(t, feedback, 5)
			for _, item := range feedback {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, float64(0), row["count"])
			}

			activeDonationsByStep := fluxgo.ConvertToList(body["active_donations_by_step"])
			assert.Len(t, activeDonationsByStep, 4)
			for _, item := range activeDonationsByStep {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, float64(0), row["count"])
				assert.Equal(t, float64(0), row["percentage"])
			}
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user cannot access dashboard", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only adms can access the dashboard", body["message"])
			assert.Equal(t, "dashboard.forbidden", body["code"])
		})

		t.Run("End date before start date", func(t *testing.T) {
			route := fmt.Sprintf(
				"%s?start_date=%s&end_date=%s",
				endpoint,
				utils.GetTodayFormattedDate(0),
				utils.GetTodayFormattedDate(-10),
			)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "End date cannot be before start date", body["message"])
			assert.Equal(t, "dashboard.invalid_date_range", body["code"])
		})
	})
}
