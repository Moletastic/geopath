package handler

import (
	"net/http"

	"github.com/Moletastic/geopath/models"
	"github.com/labstack/echo/v4"
)

// GetPath retorna una respuesta con el Path con menos trasbordos
// de acuerdo a las coordenadas de origen y destino

// GetPath godoc
// @Summary Devuelve la ruta con menor cantidad de transbordos
// @tags Geopaths
// @Description GetPath retorna una respuesta con el Path con menos trasbordos de acuerdo a las coordenadas de origen y destino
// @Description Se recomienda probar con las siguientes coordenadas
// @Description Origen: -33.443873,-70.634870
// @Description Destino: -33.4163487,-70.5683719
// @Produce json
// @Param Coordenadas origin query string true "Coordenadas Origen"
// @Param Coordenadas destination query string true "Coordenadas Destino"
// @Success 200 {object} models.PathResponse
// @Failure 400 {object} HTTPError "Se requieren coordenadas de origen y destino, ambas válidas"
// @Router /path?origin={origin}&destination={destination} [get]
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

// GetMicrobus retorna el microbus con el id fijado

// GetMicrobus godoc
// @Summary Devuelve un microbus
// @tags Geopaths
// @Description GetMicrobus devuelve el microbus correspondiente al id entregado
// @Description Probar con: 210
// @Produce json
// @Param Microbus id query string true "Microbus ID"
// @Success 200 {object} models.MicroBus
// @Failure 400 {object} HTTPError "Se requiere el id del Microbus"
// @Router /microbus [get]
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

// GetParadero retorna el paradero con el id fijado

// GetParadero godoc
// @Summary Devuelve un paradero
// @tags Geopaths
// @Description GetParadero devuelve el paradero correspondiente al id entregado
// @Description Probar con: PA1
// @Produce json
// @Param Paradero id query string true "Paradero ID"
// @Success 200 {object} models.Paradero
// @Failure 400 {object} HTTPError "Se requiere el id del paradero"
// @Router /paradero [get]
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

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"`
}
