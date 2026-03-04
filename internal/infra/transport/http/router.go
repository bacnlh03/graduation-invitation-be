package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, guestH *GuestHandler, configH *ConfigHandler, authH *AuthHandler, uploadH *UploadHandler, jwtSecret string) {
	api := e.Group("/api")
	apiV1 := api.Group("/v1")

	// --- Public routes ---
	apiV1.POST("/auth/login", authH.Login)
	apiV1.GET("/guests/:id", guestH.GetGuest)                   // public: invitation page reads guest info
	apiV1.POST("/guests/:id/confirm", guestH.ConfirmAttendance) // public: guest confirms from invitation link
	apiV1.GET("/config/invitation", configH.GetInvitationConfig)

	// --- Protected routes (require JWT) ---
	admin := apiV1.Group("", JWTMiddleware(jwtSecret))
	admin.POST("/guests/bulk", guestH.BulkCreate)
	admin.GET("/guests", guestH.GetGuests)
	admin.PATCH("/guests/:id", guestH.UpdateGuest)
	admin.PATCH("/guests/:id/status", guestH.UpdateStatus)
	admin.DELETE("/guests/:id", guestH.DeleteGuest)

	admin.PUT("/config/invitation", configH.UpdateInvitationConfig)
	admin.POST("/upload", uploadH.UploadFile)
}
