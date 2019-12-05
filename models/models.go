package models

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Moletastic/geopath/utils"
)

// Coordenada almacena información de latitud y longitud
type Coordenada struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
}

// EqualsOrClose retorna true, si la coordenada es igual
// o próxima (bajo un delta) a otra
func (c *Coordenada) EqualsOrClose(coord Coordenada, delta float64) bool {
	if c.Latitud == coord.Latitud && c.Longitud == coord.Longitud {
		return true
	}
	d := GetDistance(*c, coord)
	if d <= delta {
		return true
	}
	return false
}

// ToStr retorna la conversión a string de Coordenada
func (c *Coordenada) ToStr() string {
	return fmt.Sprintf("%f,%f", c.Latitud, c.Longitud)
}

// GetDistance retorna la distancia entre 2 coordenadas
func GetDistance(a, b Coordenada) float64 {
	R := 6373.0
	lat1 := utils.ToRad(a.Latitud)
	lat2 := utils.ToRad(b.Latitud)
	lon1 := utils.ToRad(a.Longitud)
	lon2 := utils.ToRad(b.Longitud)
	dlon := lon2 - lon1
	dlat := lat2 - lat1
	fact1 := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	fact2 := 2 * math.Atan2(math.Sqrt(fact1), math.Sqrt(1-fact1))
	return fact2 * R
}

// StrToCoord retorna una coordenada mediante un
// string o un error
func StrToCoord(str string) (Coordenada, error) {
	var geo Coordenada
	coords := strings.Split(str, ",")
	if len(coords) == 1 {
		return geo, errors.New("Coordenada inválida")
	}
	lat, err := strconv.ParseFloat(coords[0], 64)

	if err != nil {
		return geo, err
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return geo, err
	}
	geo.Latitud = lat
	geo.Longitud = lon
	return geo, nil
}
