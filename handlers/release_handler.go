package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"clocking/utils"
)

// HandleRelease handles POST /release
func HandleRelease(db *sql.DB) http.HandlerFunc {
	type request struct {
		EventName     string `json:"eventName"`
		ReleaseLatDMS string `json:"releaseLatDMS"`
		ReleaseLngDMS string `json:"releaseLngDMS"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.EventName == "" || req.ReleaseLatDMS == "" || req.ReleaseLngDMS == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		lat, err := utils.ParseDMS(req.ReleaseLatDMS)
		if err != nil {
			http.Error(w, "Invalid latitude DMS", http.StatusBadRequest)
			return
		}
		lng, err := utils.ParseDMS(req.ReleaseLngDMS)
		if err != nil {
			http.Error(w, "Invalid longitude DMS", http.StatusBadRequest)
			return
		}

		lat = utils.RoundTo6Decimals(lat)
		lng = utils.RoundTo6Decimals(lng)

		now := time.Now().UTC().Truncate(time.Minute)
		var eventID int
		err = db.QueryRow(`
			INSERT INTO events (eventName, releaseTime, releaseLat, releaseLng, releaseLatDMS, releaseLngDMS)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING eventID
		`, req.EventName, now, lat, lng, req.ReleaseLatDMS, req.ReleaseLngDMS).Scan(&eventID)
		if err != nil {
			http.Error(w, "DB insert error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"eventID":     eventID,
			"eventName":   req.EventName,
			"releaseTime": now,
			"releaseLat":  lat,
			"releaseLng":  lng,
		})
	}
}
