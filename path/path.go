package path

import "github.com/Moletastic/geopath/models"

type StorePath interface {
	GetPathToDest(origin, dest *models.Coordenada) ([]models.Path, error)
}
