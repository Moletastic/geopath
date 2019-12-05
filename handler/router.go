package handler

import "github.com/labstack/echo/v4"

// Register registra rutas a un grupo de rutas
func (h *Handler) Register(v4 *echo.Group) {
	path := v4.Group("/path")
	path.GET("", h.GetPath)
	microbus := v4.Group("/microbus")
	microbus.GET("", h.GetMicrobus)
	paradero := v4.Group("/paradero")
	paradero.GET("", h.GetParadero)
}
