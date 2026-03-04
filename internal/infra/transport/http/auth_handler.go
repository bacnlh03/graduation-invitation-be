package http

import (
	"graduation-invitation/internal/app/auth"
	nethttp "net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service auth.AuthService
}

func NewAuthHandler(svc auth.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(nethttp.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
	}

	resp, err := h.service.Login(req)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			return c.JSON(nethttp.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(nethttp.StatusInternalServerError, map[string]string{"error": "Lỗi server"})
	}

	return c.JSON(nethttp.StatusOK, resp)
}
