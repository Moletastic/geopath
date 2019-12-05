package handler

import "github.com/Moletastic/geopath/path"

// Handler almacena aquellos store utilizados por las peticiones
type Handler struct {
	PathStore path.StorePath
}

// NewHandler retorna un nuevo Handler con los store entregados
func NewHandler(store path.StorePath) *Handler {
	return &Handler{
		PathStore: store,
	}
}
