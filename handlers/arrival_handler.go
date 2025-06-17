package handlers

import (
	"clocking/utils"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleArrival(db *sql.DB) fiber.Handler {
	type request struct {
		EventID int `json:"eventID"`
	}

	return func(c *fiber.Ctx) error {
		log.Println("📥 [ARRIVAL] Request received")

		// Get user ID from cookie
		userIDStr := c.Cookies("user_id")
		if userIDStr == "" {
			log.Println("❌ [AUTH] No user_id cookie found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not logged in",
			})
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println("❌ [AUTH] Invalid user_id cookie value:", userIDStr)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid user ID",
			})
		}
		log.Printf("✅ [AUTH] User ID from cookie: %d\n", userID)

		// Parse JSON body
		var req request
		if err := c.BodyParser(&req); err != nil {
			log.Println("❌ [BODY] Failed to parse request body:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid JSON",
			})
		}
		log.Printf("📄 [BODY] Parsed request: EventID = %d\n", req.EventID)

		// Get event info
		var releaseTime time.Time
		var relLat, relLng float64
		err = db.QueryRow(`
			SELECT releaseTime, releaseLat, releaseLng FROM events WHERE eventID = $1
		`, req.EventID).Scan(&releaseTime, &relLat, &relLng)
		if err != nil {
			log.Printf("❌ [DB] EventID %d not found: %v\n", req.EventID, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Event not found",
			})
		}
		log.Printf("✅ [DB] Fetched event: releaseTime=%v, releaseLat=%.6f, releaseLng=%.6f\n", releaseTime, relLat, relLng)

		// Get user coordinates
		var userLat, userLng float64
		err = db.QueryRow(`
			SELECT latitude, longitude FROM users WHERE id = $1
		`, userID).Scan(&userLat, &userLng)
		if err != nil {
			log.Printf("❌ [DB] UserID %d not found: %v\n", userID, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		log.Printf("✅ [DB] Fetched user location: lat=%.6f, lng=%.6f\n", userLat, userLng)

		// Compute distance and speed
		distKm := utils.HaversineDistance(userLat, userLng, relLat, relLng)
		distIn60ths := distKm * 1000 * 60
		loc, _ := time.LoadLocation("Asia/Manila")
		arrivedAt := time.Now().In(loc)

		log.Printf("📏 [CALC] Distance: %.3f km (%.2f meters/60s)\n", distKm, distIn60ths)
		log.Println("⏰ [TIME] Arrived at:", arrivedAt.Format(time.RFC3339Nano))

		flyingSecs := arrivedAt.Sub(releaseTime).Seconds()
		if flyingSecs <= 0 {
			log.Println("❌ [TIME] Arrival before release")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Arrival before release",
			})
		}

		speed := distIn60ths / flyingSecs
		log.Printf("🚀 [SPEED] Computed speed: %.3f m/min\n", speed)

		// Insert into DB
		_, err = db.Exec(`
			INSERT INTO arrivals (userID, eventID, arrivedAt, speed)
			VALUES ($1, $2, $3, $4)
		`, userID, req.EventID, arrivedAt, speed)
		if err != nil {
			log.Println("❌ [DB] Failed to insert arrival:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "DB insert error",
			})
		}
		log.Println("✅ [DB] Arrival recorded successfully")

		// Return response
		return c.JSON(fiber.Map{
			"userID":     userID,
			"eventID":    req.EventID,
			"arrivedAt":  arrivedAt.Format("2006-01-02 03:04:05.000000 PM"),
			"distanceKm": utils.RoundTo3Decimals(distKm),
			"speed":      utils.RoundTo3Decimals(speed),
		})
	}
}
