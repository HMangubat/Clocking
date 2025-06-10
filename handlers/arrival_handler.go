package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"clocking/utils"
)

// HandleArrival handles POST /arrive
func HandleArrival(db *sql.DB) http.HandlerFunc {
	type request struct {
		UserID  int `json:"userID"`
		EventID int `json:"eventID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var releaseTime time.Time
		var relLat, relLng float64
		err := db.QueryRow(`
			SELECT releaseTime, releaseLat, releaseLng FROM events WHERE eventID = $1
		`, req.EventID).Scan(&releaseTime, &relLat, &relLng)
		if err != nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		var userLat, userLng float64
		err = db.QueryRow(`
			SELECT latitude, longitude FROM users WHERE id = $1
		`, req.UserID).Scan(&userLat, &userLng)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		distKm := utils.HaversineDistance(userLat, userLng, relLat, relLng)
		distIn60ths := distKm * 1000 * 60

		loc, _ := time.LoadLocation("Asia/Manila")
		arrivedAt := time.Now().In(loc)
		log.Println("Arrival logged at:", arrivedAt)

		flyingSecs := arrivedAt.Sub(releaseTime).Seconds()
		if flyingSecs <= 0 {
			http.Error(w, "Arrival before release", http.StatusBadRequest)
			return
		}

		speed := distIn60ths / flyingSecs

		_, err = db.Exec(`
			INSERT INTO arrivals (userID, eventID, arrivedAt, speed)
			VALUES ($1, $2, $3, $4)
		`, req.UserID, req.EventID, arrivedAt, speed)
		if err != nil {
			http.Error(w, "DB insert error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"userID":     req.UserID,
			"eventID":    req.EventID,
			"arrivedAt":  arrivedAt.Format("2006-01-02 03:04:05.000000 PM"),
			"distanceKm": utils.RoundTo3Decimals(distKm),
			"speed":      utils.RoundTo3Decimals(speed),
		})
	}
}
