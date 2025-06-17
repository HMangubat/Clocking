package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("📥 Received login request")

		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&input); err != nil {
			fmt.Println("❌ Failed to parse JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		var userID int
		var passwordHash string

		err := db.QueryRow(
			`SELECT id, password_hash FROM users WHERE username = $1`, input.Username,
		).Scan(&userID, &passwordHash)

		if err != nil {
			fmt.Printf("⚠️ Login failed for user '%s': user not found or DB error\n", input.Username)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
			fmt.Printf("⚠️ Login failed for user '%s': invalid password\n", input.Username)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "user_id",
			Value:    strconv.Itoa(userID),
			Path:     "/",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
		})
		fmt.Printf("✅ User '%s' logged in successfully (ID: %d)\n", input.Username, userID)

		return c.JSON(fiber.Map{
			"message": "Login successful",
		})
	}
}
