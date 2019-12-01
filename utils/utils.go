package utils

import (
	"math"
)

func StrListContains(s []string, e *string) bool {
	for _, a := range s {
		if &a == e {
			return true
		}
	}
	return false
}

func ToRad(num float64) float64 {
	return num * math.Pi / 180
}
