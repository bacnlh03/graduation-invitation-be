package http

import (
	"encoding/json"
	"graduation-invitation/internal/app/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	service config.ConfigService
}

func NewConfigHandler(svc config.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: svc}
}

func (h *ConfigHandler) GetInvitationConfig(c echo.Context) error {
	ctx := c.Request().Context()
	cfg, err := h.service.GetInvitationConfig(ctx)
	if err != nil {
		// Just log error and return empty? Or 500?
		// For now, let's just return empty object if error or not found
		return c.JSON(http.StatusOK, map[string]interface{}{})
	}
	if cfg == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{})
	}

	// Need to unmarshal to ensure it sends as JSON object, not string
	var res interface{}
	if err := json.Unmarshal(cfg, &res); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *ConfigHandler) UpdateInvitationConfig(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()
	if err := h.service.UpdateInvitationConfig(ctx, jsonData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Valid updated"})
}
