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
			return echo.NewHTTPError(http.StatusBadRequest, "Origen: Coordenadas inválidas")
		}
		dest, err := models.StrToCoord(destParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Destino: Coordenadas inválidas")
		}
		data, err := h.PathStore.GetPathToDest(&origin, &dest)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		response := models.PathResponse{Data: data, Status: http.StatusOK}
		return c.JSON(http.StatusOK, &response)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Se requieren coordenadas de origen y destino")
}

// GetMicrobus will be commented
func (h *Handler) GetMicrobus(c echo.Context) error {
	id := c.QueryParam("id")
	if id != "" {
		microbus, err := h.PathStore.GetMicroBusByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, microbus)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Se requiere el id del Microbus")
}

// GetParadero will be commented
func (h *Handler) GetParadero(c echo.Context) error {
	id := c.QueryParam("id")
	if id != "" {
		paradero, err := h.PathStore.GetParadeByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, paradero)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Se requiere el id del Paradero")
}
