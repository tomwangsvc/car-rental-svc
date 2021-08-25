package routes

import (
	"net/http"

	lib_http "github.com/tomwangsvc/lib-svc/http"
)

// @Summary health check
// @Param type query string true "health=true"
// @Produce json
// @Success 200
// @Router /car-svc [get]
func (c client) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("health") == "true" {
			lib_http.RenderResponse(r.Context(), w, Health{
				BuildDate:   "@foo_BUILD_DATE@",
				BuildNumber: "@foo_BUILD_NUMBER@",
				CommitId:    "@foo_COMMIT_ID@",
				Env:         c.config.Env.Id,
				Status:      "OK",
				Svc:         c.config.Env.SvcId,
			})
		}
	}
}

type Health struct {
	BuildDate   string `json:"build_date"`
	BuildNumber string `json:"build_number"`
	CommitId    string `json:"commit_id"`
	Env         string `json:"env"`
	Status      string `json:"status"`
	Svc         string `json:"svc"`
}
