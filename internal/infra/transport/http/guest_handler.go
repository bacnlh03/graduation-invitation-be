package http

import (
	"graduation-invitation/internal/app/guest"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GuestHandler struct {
	service guest.GuestService
}

func NewGuestHandler(svc guest.GuestService) *GuestHandler {
	return &GuestHandler{service: svc}
}

func (h *GuestHandler) BulkCreate(c echo.Context) error {
	var req guest.BulkCreateGuestRequest

	// Bind dữ liệu từ JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
	}

	// Gọi xuống tầng App
	ctx := c.Request().Context()
	if err := h.service.BulkRegister(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Thêm danh sách thành công"})
}

func (h *GuestHandler) UpdateGuest(c echo.Context) error {
	var req guest.UpdateGuestRequest
	id := c.Param("id")

	// Bind dữ liệu từ JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
	}

	// Gọi xuống tầng App
	ctx := c.Request().Context()
	if err := h.service.UpdateGuest(ctx, id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cập nhật thành công"})
}

func (h *GuestHandler) GetGuests(c echo.Context) error {
	var req guest.FilterGuestsRequest

	// Bind query params vào struct
	// Echo sẽ tự hiểu status=true là bool true, status=false là bool false
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filter"})
	}

	ctx := c.Request().Context()
	// Gọi xuống Service (Service sẽ gọi Repo)
	guests, err := h.service.ListGuests(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, guests)
}
func (h *GuestHandler) GetGuest(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	guest, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if guest == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Không tìm thấy khách mời"})
	}

	return c.JSON(http.StatusOK, guest)
}
func (h *GuestHandler) ConfirmAttendance(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	if err := h.service.ConfirmAttendance(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Xác nhận tham dự thành công"})
}

func (h *GuestHandler) DeleteGuest(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	if err := h.service.DeleteGuest(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Xóa khách mời thành công"})
}

func (h *GuestHandler) UpdateStatus(c echo.Context) error {
	id := c.Param("id")
	var req guest.UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
	}
	// Validate status 0, 1, or 2
	if req.Status < 0 || req.Status > 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status phải là 0, 1 hoặc 2"})
	}
	ctx := c.Request().Context()
	if err := h.service.UpdateStatus(ctx, id, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cập nhật trạng thái thành công"})
}
