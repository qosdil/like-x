package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

const (
	authMiddlewareHeader = "Auth-User-ID"
)

// AuthMiddleware is a authentication middleware that extracts the user ID from the authMiddlewareHeader header.
var AuthMiddleware = func(c fiber.Ctx) error {
	authUserID := fiber.GetReqHeader[uint](c, authMiddlewareHeader)
	if authUserID == 0 {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// Store in request-local context
	c.Locals("auth_user_id", authUserID)
	return c.Next()
}
