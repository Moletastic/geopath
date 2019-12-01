package handler

import "github.com/Moletastic/geopath/path"

// Handler saves every Store
type Handler struct {
	PathStore path.StorePath
}

// NewHandler returns a new Handler pointer
func NewHandler(store path.StorePath) *Handler {
	return &Handler{
		PathStore: store,
	}
}
