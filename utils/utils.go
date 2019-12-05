package utils

import (
	"math"
)

// StrListContains retorna si una lista contiene un string
func StrListContains(list []string, element *string) bool {
	for _, item := range list {
		if &item == element {
			return true
		}
	}
	return false
}

// ToRad retorna la conversión en radianes de un número
func ToRad(num float64) float64 {
	return num * math.Pi / 180
}
