package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"golang.org/x/crypto/bcrypt"
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
	hist := metrics.GetHistogramVec("post_api_register")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	// TODO Consider if this functionality can be shared with ui controller. Logic should probably be in service.
	var registerRequest models.RegisterRequest
	_ = json.NewDecoder(r.Body).Decode(&registerRequest)

	updateLatest(r)
	var registerError string
	if r.Method == "POST" {
		if len(registerRequest.Username) == 0 {
			registerError = "You have to enter a username"
		} else if len(registerRequest.Email) == 0 || !strings.Contains(registerRequest.Email, "@") {
			registerError = "You have to enter a valid email address"
		} else if len(registerRequest.Password) == 0 {
			registerError = "You have to enter a password"
		} else if _, err := services.GetUserID(registerRequest.Username); err == nil {
			registerError = "The username is already taken"
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Print(err)
			}

			user := models.User{Username: registerRequest.Username, Email: registerRequest.Email, PwHash: string(hash)}
			err = services.CreateUser(user)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if registerError != "" {
			var errorCode = 400
			error := map[string]interface{}{
				"status":    errorCode,
				"error_msg": registerError,
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
}
