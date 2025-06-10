package models

// User represents user data.
type User struct {
	ID            int     `json:"id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	LatitudeDMS   string  `json:"latitudeDms"`
	LongitudeDMS  string  `json:"longitudeDms"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	CreatedAt     string  `json:"createdAt"`
}
