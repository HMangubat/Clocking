package models

// User represents user data.
type User struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Password     string  `json:"password,omitempty"`
	Firstname    string  `json:"firstname"`
	Middlename   string  `json:"middlename"`
	Lastname     string  `json:"lastname"`
	PasswordHash string  `json:"-"`
	LatitudeDMS  string  `json:"latitudeDms"`
	LongitudeDMS string  `json:"longitudeDms"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	CreatedAt    string  `json:"createdAt"`
}
