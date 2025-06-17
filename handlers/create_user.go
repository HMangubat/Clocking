package handlers

import (
	"clocking/models"
	"clocking/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var u models.User

		if err := c.BodyParser(&u); err != nil {
			log.Println("[ERROR] Invalid JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid JSON",
			})
		}

		log.Printf("[DEBUG] Decoded user struct: %+v\n", u)
		log.Printf("[INFO] Creating user: %s (%s %s %s)\n", u.Username, u.Firstname, u.Middlename, u.Lastname)

		// Parse DMS to float
		lat, err := utils.ParseDMS(u.LatitudeDMS)
		if err != nil {
			log.Println("[ERROR] Invalid latitude:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid latitude format",
			})
		}
		lng, err := utils.ParseDMS(u.LongitudeDMS)
		if err != nil {
			log.Println("[ERROR] Invalid longitude:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid longitude format",
			})
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[ERROR] Password hash failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Password hash failed",
			})
		}

		// Insert user
		_, err = db.Exec(`
			INSERT INTO users 
				(username, email, password_hash, firstname, middlename, lastname, latitude_dms, longitude_dms, latitude, longitude)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, u.Username, u.Email, hash, u.Firstname, u.Middlename, u.Lastname, u.LatitudeDMS, u.LongitudeDMS, lat, lng)

		if err != nil {
			log.Println("[ERROR] Database insert failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Registration failed",
			})
		}

		log.Printf("[INFO] User created successfully: %s\n", u.Username)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered",
		})
	}
}
