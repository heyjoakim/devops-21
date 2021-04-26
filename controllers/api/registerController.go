package api

import (
	"encoding/json"
	"net/http"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

// Messages godoc
// @Summary Registers a user
// @Description Registers a user, provided that the given info passes all checks.
// @Accept json
// @Produce json
// @Success 203
// @Failure 400 {string} string "unauthorized"
// @Router /api/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)
	hist := metrics.GetHistogramVec("post_api_register")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	// TODO Consider if this functionality can be shared with ui controller. Logic should probably be in service.
	var registerRequest models.RegisterRequest
	_ = json.NewDecoder(r.Body).Decode(&registerRequest)
	user := models.UserCreateRequest{
		Username:  registerRequest.Username,
		Email:     registerRequest.Email,
		Password:  registerRequest.Password,
		Password2: registerRequest.Password,
	}
	err := services.CreateUser(user)
	if err != nil {
		var errorCode = 400
		error := map[string]interface{}{
			"status":    errorCode,
			"error_msg": err.Error(),
		}
		jsonData, _ := helpers.Serialize(error)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(jsonData)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
