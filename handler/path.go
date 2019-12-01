package handler

import (
	"net/http"

	"github.com/Moletastic/geopath/models"
	"github.com/labstack/echo/v4"
)

// GetPath will be commented
func (h *Handler) GetPath(c echo.Context) error {
	originParam := c.QueryParam("origin")
	destParam := c.QueryParam("destination")
	if originParam != "" && destParam != "" {
		origin, err := models.StrToCoord(originParam)
		if err != nil {
			return echo.NewHTTPError(400, "Origen: Coordenadas inválidas")
		}
		dest, err := models.StrToCoord(destParam)
		if err != nil {
			return echo.NewHTTPError(400, "Destino: Coordenadas inválidas")
		}
		data, err := h.PathStore.GetPathToDest(&origin, &dest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		response := models.PathResponse{Data: data, Status: http.StatusOK}
		return c.JSON(http.StatusOK, &response)
	}
	return echo.NewHTTPError(422, "Se requieren coordenadas de origen y destino")
}
