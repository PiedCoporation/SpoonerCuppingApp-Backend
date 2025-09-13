package main

import "backend/internal/app"

// @title Coffee Cupping API
// @version 1.0
// @description API documentation for Coffee Cupping backend.
// @BasePath /v1
// @schemes http
// @accept json
// @produce json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide a Bearer token: "Bearer {token}"
func main() {
	app.Run()
}
