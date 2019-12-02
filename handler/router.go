package handler

import "github.com/labstack/echo/v4"

// Register adds a group for router
func (h *Handler) Register(v4 *echo.Group) {
	path := v4.Group("/path")
	path.GET("", h.GetPath)
	microbus := v4.Group("/microbus")
	microbus.GET("", h.GetMicrobus)
	paradero := v4.Group("/paradero")
	paradero.GET("", h.GetParadero)
}
