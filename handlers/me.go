package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func MeHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("📥 Received /me request")

		cookie := c.Cookies("user_id")
		if cookie == "" {
			fmt.Println("[WARN] Missing user_id cookie in /me request")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userID, err := strconv.Atoi(cookie)
		if err != nil {
			fmt.Printf("[ERROR] Invalid user_id cookie value: %s\n", cookie)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}
		fmt.Printf("🔍 Fetching user info for user_id: %d\n", userID)

		var u struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			Firstname  string `json:"firstname"`
			Middlename string `json:"middlename"`
			Lastname   string `json:"lastname"`
		}

		err = db.QueryRow(`
			SELECT username, email, firstname, middlename, lastname 
			FROM users WHERE id = $1
		`, userID).Scan(&u.Username, &u.Email, &u.Firstname, &u.Middlename, &u.Lastname)

		if err != nil {
			fmt.Printf("[ERROR] No user found for user_id %d: %v\n", userID, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		fmt.Printf("✅ User data retrieved: %s \n", u)
		return c.JSON(u)
	}
}
