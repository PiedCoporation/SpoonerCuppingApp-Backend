package controller

// MessageResponse represents a simple message payload.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents a simple error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenData holds access/refresh tokens.
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse wraps tokens with a message.
type TokenResponse struct {
	Message string    `json:"message"`
	Data    TokenData `json:"data"`
}
