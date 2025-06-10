package main

import (
	"log"
	"net/http"

	"clocking/config"
	"clocking/handlers"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	http.HandleFunc("/users", handlers.CreateUserHandler(db))
	http.HandleFunc("/release", handlers.HandleRelease(db))
	http.HandleFunc("/arrive", handlers.HandleArrival(db))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	log.Println("✅ Server running at http://localhost:8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}




// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"math"
// 	"net/http"
// 	"regexp"
// 	"strconv"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// var db *sql.DB

// type User struct {
// 	ID            int     `json:"id"`
// 	Username      string  `json:"username"`
// 	Email         string  `json:"email"`
// 	LatitudeDMS   string  `json:"latitudeDms"`
// 	LongitudeDMS  string  `json:"longitudeDms"`
// 	Latitude      float64 `json:"latitude"`
// 	Longitude     float64 `json:"longitude"`
// 	CreatedAt     string  `json:"createdAt"`
// }


// func main() {
//     var err error
//     db, err = sql.Open("postgres", "postgres://postgres:123@10.9.2.30:5432/kafka?sslmode=disable")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer db.Close()

// 	http.HandleFunc("/users", createUserHandler)
//  http.HandleFunc("/release", handleRelease)
// 	http.HandleFunc("/arrive", handleArrival)

//     http.Handle("/", http.FileServer(http.Dir("./static")))

//     log.Println("✅ Server running at http://localhost:8085")
//     log.Fatal(http.ListenAndServe(":8085", nil))
// }

// func InsertUser(db *sql.DB, user User) error {
// 	// Convert DMS to decimal
// 	latDecimal, err := parseDMS(user.LatitudeDMS)
// 	if err != nil {
// 		return fmt.Errorf("invalid latitude DMS: %v", err)
// 	}
// 	lngDecimal, err := parseDMS(user.LongitudeDMS)
// 	if err != nil {
// 		return fmt.Errorf("invalid longitude DMS: %v", err)
// 	}

// 	query := `
// 		INSERT INTO users (username, email, latitude_dms, longitude_dms, latitude, longitude)
// 		VALUES ($1, $2, $3, $4, $5, $6)
// 	`
// 	_, err = db.Exec(query, user.Username, user.Email, user.LatitudeDMS, user.LongitudeDMS, latDecimal, lngDecimal)
// 	return err
// }

// func createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User
// 	json.NewDecoder(r.Body).Decode(&user)

// 	err := InsertUser(db, user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
// }


// func handleRelease(w http.ResponseWriter, r *http.Request) {
// 	type ReleaseRequest struct {
// 		EventName     string `json:"eventName"`
// 		ReleaseLatDMS string `json:"releaseLatDMS"`
// 		ReleaseLngDMS string `json:"releaseLngDMS"`
// 	}

// 	var req ReleaseRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
// 		return
// 	}

// 	// Validate required fields
// 	if req.EventName == "" || req.ReleaseLatDMS == "" || req.ReleaseLngDMS == "" {
// 		http.Error(w, "Missing required fields", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse DMS strings to decimal degrees
// 	lat, err := parseDMS(req.ReleaseLatDMS)
// 	if err != nil {
// 		http.Error(w, "Invalid releaseLatDMS format", http.StatusBadRequest)
// 		return
// 	}

// 	lng, err := parseDMS(req.ReleaseLngDMS)
// 	if err != nil {
// 		http.Error(w, "Invalid releaseLngDMS format", http.StatusBadRequest)
// 		return
// 	}

// 	// Round off to 6 decimals
// 	lat = roundTo6Decimals(lat)
// 	lng = roundTo6Decimals(lng)

// 	// Use current time truncated to minute
// 	now := time.Now().UTC().Truncate(time.Minute)

// 	// Insert into DB
// 	var eventID int
// 	err = db.QueryRow(
// 		`INSERT INTO events (eventName, releaseTime, releaseLat, releaseLng, releaseLatDMS, releaseLngDMS) 
// 		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING eventID`,
// 		req.EventName, now, lat, lng, req.ReleaseLatDMS, req.ReleaseLngDMS).Scan(&eventID)

// 	if err != nil {
// 		http.Error(w, "Database insert error", http.StatusInternalServerError)
// 		return
// 	}

// 	resp := map[string]interface{}{
// 		"eventID":     eventID,
// 		"eventName":   req.EventName,
// 		"releaseTime": now,
// 		"releaseLat":  lat,
// 		"releaseLng":  lng,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }
// func roundTo6Decimals(f float64) float64 {
//     return math.Round(f*1_000_000) / 1_000_000
// }
// func parseDMS(dms string) (float64, error) {
// 	// Example input: 12°36′15.47″ N
// 	// Regex to extract degrees, minutes, seconds, and direction (N,S,E,W)
// 	re := regexp.MustCompile(`(\d+)°(\d+)′([\d.]+)″\s*([NSEW])`)
// 	matches := re.FindStringSubmatch(dms)
// 	if len(matches) != 5 {
// 		return 0, errors.New("Invalid DMS format")
// 	}

// 	deg, _ := strconv.ParseFloat(matches[1], 64)
// 	min, _ := strconv.ParseFloat(matches[2], 64)
// 	sec, _ := strconv.ParseFloat(matches[3], 64)
// 	dir := matches[4]

// 	decimal := deg + min/60 + sec/3600

// 	// Negative for South and West
// 	if dir == "S" || dir == "W" {
// 		decimal = -decimal
// 	}

// 	return decimal, nil
// }

// func handleArrival(w http.ResponseWriter, r *http.Request) {
// 	type ArrivalRequest struct {
// 		UserID  int `json:"userID"`
// 		EventID int `json:"eventID"`
// 	}

// 	var req ArrivalRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
// 		return
// 	}

// 	// Get event release time and coordinates
// 	var releaseTime time.Time
// 	var releaseLat, releaseLng float64
// 	err := db.QueryRow(`
// 		SELECT releaseTime, releaseLat, releaseLng FROM events WHERE eventID = $1
// 	`, req.EventID).Scan(&releaseTime, &releaseLat, &releaseLng)
// 	if err != nil {
// 		http.Error(w, "Event not found", http.StatusNotFound)
// 		return
// 	}

// 	// Get user coordinates
// 	var userLat, userLng float64
// 	err = db.QueryRow(`
// 		SELECT latitude, longitude FROM users WHERE id = $1
// 	`, req.UserID).Scan(&userLat, &userLng)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	// Calculate distance in km using Haversine
// 	distKm := haversineDistance(userLat, userLng, releaseLat, releaseLng)

// 	// Convert to 60ths: meters × 60
// 	distIn60ths := distKm * 1000 * 60

// 	// Get current time truncated to minute
// 	loc, _ := time.LoadLocation("Asia/Manila")
// 	arrivedAt := time.Now().In(loc)

// 	log.Printf(arrivedAt.String())

// 	// Compute time difference in seconds (in 60ths)
// 	flyingSecs := arrivedAt.Sub(releaseTime).Seconds()
// 	if flyingSecs <= 0 {
// 		http.Error(w, "Arrival time must be after release time", http.StatusBadRequest)
// 		return
// 	}

// 	// Compute velocity (m/min)
// 	speed := distIn60ths / flyingSecs

// 	// Save to DB
// 	_, err = db.Exec(`
// 		INSERT INTO arrivals (userID, eventID, arrivedAt, speed)
// 		VALUES ($1, $2, $3, $4)
// 	`, req.UserID, req.EventID, arrivedAt, speed)
// 	if err != nil {
// 		http.Error(w, "Database insert error", http.StatusInternalServerError)
// 		return
// 	}

// 	resp := map[string]interface{}{
// 		"userID":     req.UserID,
// 		"eventID":    req.EventID,
// 		"arrivedAt":  arrivedAt.Format("2006-01-02 15:04:05.000000"),
// 		"distanceKm": roundTo3Decimals(distKm),
// 		"speed":      roundTo3Decimals(speed), // m/min
// 	}


// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }

// func roundTo3Decimals(f float64) float64 {
// 	return math.Round(f*1000) / 1000
// }

// // Haversine formula to calculate great-circle distance
// func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
// 	const R = 6371 // Earth radius in km
// 	lat1Rad := lat1 * math.Pi / 180
// 	lon1Rad := lon1 * math.Pi / 180
// 	lat2Rad := lat2 * math.Pi / 180
// 	lon2Rad := lon2 * math.Pi / 180

// 	dlat := lat2Rad - lat1Rad
// 	dlon := lon2Rad - lon1Rad

// 	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
// 		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
// 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
// 	return R * c
// }





































// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// var db *sql.DB

// func main() {
// 	var err error
// 	db, err = sql.Open("postgres", "postgres://postgres:123@10.9.2.30:5432/kafka?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	http.HandleFunc("/release", handleRelease)
// 	http.HandleFunc("/arrival", handleArrival)
// 	http.Handle("/", http.FileServer(http.Dir("./static")))

// 	log.Println("✅ Server running at http://localhost:8085")
// 	log.Fatal(http.ListenAndServe(":8085", nil))
// }

// func handleRelease(w http.ResponseWriter, r *http.Request) {
// 	//now := time.Now().UTC() // read seconds
// 	now := time.Now().UTC().Truncate(time.Minute) //seconds always 00
// 	var id int
// 	err := db.QueryRow(`INSERT INTO clockings (release_time) VALUES ($1) RETURNING id`, now).Scan(&id)
// 	if err != nil {
// 		http.Error(w, "Database error", 500)
// 		return
// 	}

// 	log.Printf("🕒 Release: ID=%d at %v", id, now)

// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"id":           id,
// 		"release_time": now,
// 	})
// }

// func handleArrival(w http.ResponseWriter, r *http.Request) {
// 	type ArrivalRequest struct {
// 		ID       int     `json:"id"`
// 		Distance float64 `json:"distance"` // in kilometers
// 	}

// 	var req ArrivalRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid JSON", 400)
// 		return
// 	}

// 	var releaseTime time.Time
// 	err := db.QueryRow(`SELECT release_time FROM clockings WHERE id = $1`, req.ID).Scan(&releaseTime)
// 	if err != nil {
// 		http.Error(w, "Release not found", 404)
// 		return
// 	}

// 	arrivalTime := time.Now().UTC()
// 	duration := arrivalTime.Sub(releaseTime).Seconds()

// 	if duration <= 0 {
// 		http.Error(w, "Invalid duration (arrival before release)", 400)
// 		return
// 	}

// 	// Convert distance to meters then to 60ths
// 	distanceMeters := req.Distance * 1000
// 	distance60ths := distanceMeters * 60

// 	// Speed in meters per 60th of second (m/m)
// 	speed := distance60ths / duration

// 	_, err = db.Exec(`UPDATE clockings SET arrival_time = $1, distance_meters = $2, speed_m_s = $3 WHERE id = $4`,
// 		arrivalTime, distanceMeters, speed, req.ID)
// 	if err != nil {
// 		http.Error(w, "Failed to update", 500)
// 		return
// 	}

// 	log.Printf("✅ Arrival recorded: ID=%d Distance=%.3f km (%.0f 60ths) Duration=%.0f s Speed=%.3f m/m",
// 		req.ID, req.Distance, distance60ths, duration, speed)

// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"id":               req.ID,
// 		"release_time":     releaseTime,
// 		"arrival_time":     arrivalTime,
// 		"distance_km":      req.Distance,
// 		"distance_60ths":   distance60ths,
// 		"time_seconds":     duration,
// 		"speed_m_per_60th": speed,
// 	})
// }
