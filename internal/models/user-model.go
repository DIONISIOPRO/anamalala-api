package models

import (
	"time"
)

// Role represents user roles in the system
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type UserRequest struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	Province string `bson:"province" json:"province" validate:"required"`
	Contact  string `bson:"contact" json:"contact" validate:"required,unique"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// User represents a user in the system
type User struct {
	ID              string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string    `bson:"name" json:"name" validate:"required"`
	Province        string    `bson:"province" json:"province" validate:"required"`
	Contact         string    `bson:"contact" json:"contact" validate:"required,unique"`
	Password        string    `bson:"password" json:"-" validate:"required"`
	Role            Role      `bson:"role" json:"role"`
	Active          bool      `bson:"active" json:"active"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
	LastLoginAt     time.Time `bson:"last_login_at" json:"last_login_at,omitempty"`
	IsLoggedIn      bool      `bson:"is_logged_in" json:"is_logged_in,omitempty"`
	ResetCode       string
	ResetCodeExpiry time.Time
}

// PasswordReset represents password reset data
type PasswordReset struct {
	Token     string    `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}

// UserLogin represents login credentials
type UserLogin struct {
	Contact  string `json:"contact" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserRegistration represents user registration data
type UserRegistration struct {
	Name     string `json:"name" validate:"required"`
	Province string `json:"province" validate:"required"`
	Contact  string `json:"contact" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Contact string `json:"contact" validate:"required"`
}

// PasswordResetConfirm represents a password reset confirmation
type PasswordResetConfirm struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse represents the user data returned to clients
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Province  string    `json:"province"`
	Contact   string    `json:"contact"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Province:  u.Province,
		Contact:   u.Contact,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

// Users represents a slice of User
type Users []User

// ToResponse converts a slice of Users to a slice of UserResponse
func (u Users) ToResponse() []UserResponse {
	response := make([]UserResponse, len(u))
	for i, user := range u {
		response[i] = user.ToResponse()
	}
	return response
}
