package domain

// RegisterDTO represents the registration request payload
type RegisterDTO struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"securePassword123" binding:"required,min=8"`
	Username string `json:"username" example:"johndoe" binding:"omitempty,min=3"`
}

// LoginDTO represents the login request payload
type LoginDTO struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"securePassword123" binding:"required"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID       int    `json:"id" example:"1"`
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username,omitempty" example:"johndoe"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// RefreshTokenDTO represents the token refresh request
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
}

// LogoutDTO represents the logout request
type LogoutDTO struct {
	RefreshToken string `json:"refresh_token" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
