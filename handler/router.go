package handler

import "github.com/labstack/echo/v4"

// Register adds a group for router
func (h *Handler) Register(v4 *echo.Group) {
	routes := v4.Group("/path")
	routes.GET("", h.GetPath)
}
