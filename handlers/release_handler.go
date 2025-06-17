package handlers

import (
	"clocking/utils"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleRelease(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type request struct {
			EventName     string `json:"eventName"`
			ReleaseLatDMS string `json:"releaseLatDMS"`
			ReleaseLngDMS string `json:"releaseLngDMS"`
		}

		var req request
		if err := c.BodyParser(&req); err != nil {
			log.Println("[ERROR] Failed to parse JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid JSON",
			})
		}

		if req.EventName == "" || req.ReleaseLatDMS == "" || req.ReleaseLngDMS == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Missing required fields",
			})
		}

		lat, err := utils.ParseDMS(req.ReleaseLatDMS)
		if err != nil {
			log.Println("[ERROR] Invalid latitude DMS:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid latitude DMS",
			})
		}

		lng, err := utils.ParseDMS(req.ReleaseLngDMS)
		if err != nil {
			log.Println("[ERROR] Invalid longitude DMS:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid longitude DMS",
			})
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
			log.Println("[ERROR] DB insert failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "DB insert error",
			})
		}

		return c.JSON(fiber.Map{
			"eventID":     eventID,
			"eventName":   req.EventName,
			"releaseTime": now,
			"releaseLat":  lat,
			"releaseLng":  lng,
		})
	}
}
