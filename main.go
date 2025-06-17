// package main

// import (
// 	"log"
// 	"net/http"

// 	"clocking/config"
// 	"clocking/handlers"
// )

// func main() {
// 	db := config.InitDB()
// 	defer db.Close()

// 	http.HandleFunc("/login", handlers.LoginHandler(db))
// 	http.HandleFunc("/me", handlers.MeHandler(db))
// 	http.HandleFunc("/logout", handlers.LogoutHandler())
// 	http.HandleFunc("/users", handlers.CreateUserHandler(db))
// 	http.HandleFunc("/release", handlers.HandleRelease(db))
// 	http.HandleFunc("/arrive", handlers.HandleArrival(db))

// 	http.Handle("/", http.FileServer(http.Dir("./static")))
// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

// 	log.Println("✅ Server running at http://localhost:8085")
// 	log.Fatal(http.ListenAndServe(":8085", nil))
// }

package main

import (
	"log"

	"clocking/config"
	"clocking/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	app := fiber.New()

	// API routes
	app.Post("/login", handlers.LoginHandler(db))
	app.Get("/me", handlers.MeHandler(db))
	// app.Post("/logout", handlers.LogoutHandler())
	app.Post("/users", handlers.CreateUserHandler(db))
	app.Post("/release", handlers.HandleRelease(db))
	app.Post("/arrive", handlers.HandleArrival(db))

	// Serve static files
	app.Static("/", "./static") // Serves index.html and others from /static
	app.Static("/static", "./static")

	log.Println("✅ Server running at http://localhost:8085")
	log.Fatal(app.Listen(":8085"))
}
