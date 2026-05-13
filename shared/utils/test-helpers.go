package utils

import (
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

var TestHeaders = fluxgo.Headers{
	"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoidXNyXzJ2ZUwxRlBwdVh4VWFaY0ZhRUM1N0JmcGNLRSIsImlhdCI6MTc2ODYwNzY5NX0.nu3pikK27PucdmX3ZobAvKCuEwZehbHtSnk4AweJix0",
}

var InvalidTestHeaders = fluxgo.Headers{
	"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoidXNyXzJ2ZUwxRlBwdVh4VWFaY0ZhRUM1N0JmcGNMViIsImlhdCI6MTc2ODU4NzA2OX0.eMkzXrc9lJ_Ah_VLvj0YgTyzrgp42SJ8rtlVm7vg_Ps",
}

func GetTodayFormattedDate(daysOffset int) string {
	return time.Now().AddDate(0, 0, daysOffset).Format("2006-01-02")
}
