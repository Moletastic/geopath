package path

import "github.com/Moletastic/geopath/models"

type StorePath interface {
	GetPathToDest(origin, dest *models.Coordenada) ([]models.Path, error)
	GetParadeByID(id string) (*models.Paradero, error)
	GetMicroBusByID(id string) (*models.MicroBus, error)
}
