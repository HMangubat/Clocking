package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"clocking/models"
	"clocking/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register user: includes password and DMS parser
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		lat, err := utils.ParseDMS(u.LatitudeDMS)
		if err != nil {
			http.Error(w, "Invalid latitude format", http.StatusBadRequest)
			return
		}
		lng, err := utils.ParseDMS(u.LongitudeDMS)
		if err != nil {
			http.Error(w, "Invalid longitude format", http.StatusBadRequest)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Password hash failed", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
            INSERT INTO users (username, email, latitude_dms, longitude_dms, latitude, longitude, password_hash)
            VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			u.Username, u.Email, u.LatitudeDMS, u.LongitudeDMS, lat, lng, hash,
		)
		if err != nil {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
	}
}

// Login user and set session cookie
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var userID int
		var passwordHash string

		err := db.QueryRow(
			`SELECT id, password_hash FROM users WHERE username = $1`, input.Username,
		).Scan(&userID, &passwordHash)

		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Set cookie to track session
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    strconv.Itoa(userID),
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login successful",
		})
	}
}

// Retrieves current user's info
func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cid, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id, _ := strconv.Atoi(cid.Value)

		var u struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}
		err = db.QueryRow("SELECT username, email FROM users WHERE id=$1", id).
			Scan(&u.Username, &u.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// Logout: delete cookie
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Optional: Only allow POST if you want stricter security
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Expire the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // Set to true if using HTTPS
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
	}
}

// package handlers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"

// 	"clocking/models"
// 	"clocking/utils"

// 	"golang.org/x/crypto/bcrypt"
// )

// func CreateUserHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user models.User
// 		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}

// 		lat, err := utils.ParseDMS(user.LatitudeDMS)
// 		if err != nil {
// 			http.Error(w, "Invalid latitude DMS", http.StatusBadRequest)
// 			return
// 		}
// 		lng, err := utils.ParseDMS(user.LongitudeDMS)
// 		if err != nil {
// 			http.Error(w, "Invalid longitude DMS", http.StatusBadRequest)
// 			return
// 		}

// 		passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			http.Error(w, "Password hashing failed", http.StatusInternalServerError)
// 			return
// 		}

// 		_, err = db.Exec(`
// 			INSERT INTO users (username, email, latitude_dms, longitude_dms, latitude, longitude, password_hash)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		`, user.Username, user.Email, user.LatitudeDMS, user.LongitudeDMS, lat, lng, passwordHash)

// 		if err != nil {
// 			http.Error(w, "DB insertion failed", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
// 	}
// }

// func LoginHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var input struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}

// 		var storedHash string
// 		err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", input.Username).Scan(&storedHash)
// 		if err != nil {
// 			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 			return
// 		}

// 		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(input.Password))
// 		if err != nil {
// 			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
// 	}
// }
