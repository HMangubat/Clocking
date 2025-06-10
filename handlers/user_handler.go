package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"clocking/models"
	"clocking/utils"
)

// CreateUserHandler handles POST /users
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lat, err := utils.ParseDMS(user.LatitudeDMS)
		if err != nil {
			http.Error(w, "Invalid latitude DMS", http.StatusBadRequest)
			return
		}
		lng, err := utils.ParseDMS(user.LongitudeDMS)
		if err != nil {
			http.Error(w, "Invalid longitude DMS", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(`
			INSERT INTO users (username, email, latitude_dms, longitude_dms, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, user.Username, user.Email, user.LatitudeDMS, user.LongitudeDMS, lat, lng)
		if err != nil {
			http.Error(w, "DB insertion failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	}
}
