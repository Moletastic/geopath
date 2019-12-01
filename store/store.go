package store

import (
	"errors"

	"github.com/Moletastic/geopath/models"
)

type PathStore struct {
	Paraderos  models.Paraderos
	IParaderos models.IndParaderos
	Buses      models.MicroBuses
	IBuses     models.IndMicroBuses
}

// NewPathStore will be commented
func NewPathStore(paraderos models.Paraderos, buses models.MicroBuses) *PathStore {
	return &PathStore{
		Paraderos:  paraderos,
		Buses:      buses,
		IParaderos: paraderos.ToIndParaderos(),
		IBuses:     buses.ToIndMicroBuses(),
	}
}

// GetPathToDest will be commented
func (store *PathStore) GetPathToDest(origin, dest *models.Coordenada) ([]models.Path, error) {
	store.Paraderos.SortByCoordDistance(origin)
	originParades := store.Paraderos.GetNextParaderos(origin)
	store.Paraderos.SortByCoordDistance(dest)
	destParades := store.Paraderos.GetNextParaderos(dest)
	if len(originParades) == 0 || len(destParades) == 0 {
		return nil, errors.New("No existen paradas cercanas")
	}
	useful := make([]models.Path, 0)
	noavalaible := make(models.RouteCodes, 0)
	for _, parade := range originParades {
		bcodes := parade.Microbuses
		for _, bcode := range bcodes {
			bus := store.IBuses[bcode]
			for _, destParade := range destParades {
				if bus.GetParaderoIndex(&destParade) != -1 {
					path := models.Path{Origin: parade, Dest: destParade}
					step := models.Ruta{Microbus: bus}
					step.Paraderos = []string{parade.Codigo, destParade.Codigo}
					step.SetDistance(store.IParaderos)
					path.Steps = append(path.Steps, step)
					useful = append(useful, path)
					return useful, nil
				} else {
					rcode := models.RouteCode{BCode: bus.Recorrido, Origin: parade.Codigo, Dest: destParade.Codigo}
					if !noavalaible.Contains(&rcode) {
						noavalaible = append(noavalaible, rcode)
					}
				}
			}
		}
	}
	if len(useful) == 0 {
		for _, rcode := range noavalaible {
			bus := store.IBuses[rcode.BCode]
			origin := store.IParaderos[rcode.Origin]
			destination := store.IParaderos[rcode.Dest]
			routes, unfound := models.GetNextRoutes(&origin, &destination, &bus, &store.IParaderos, &store.IBuses)
			for _, route := range routes {
				path := models.Path{Origin: origin, Dest: destination}
				step := models.Ruta{Microbus: bus}
				step.Paraderos = []string{origin.Codigo, route.Paraderos[0]}
				step.SetDistance(store.IParaderos)
				path.Steps = append(path.Steps, step)
				path.Steps = append(path.Steps, route)
				useful = append(useful, path)
				return useful, nil
			}
			if len(routes) == 0 {
				for _, pair := range unfound {
					transBus := store.IBuses[pair.BCode]
					transPar := store.IParaderos[pair.PCode]
					transferedRoutes, unfoundTrans := models.GetNextRoutes(&transPar, &destination, &transBus, &store.IParaderos, &store.IBuses)
					for _, transRoute := range transferedRoutes {
						path := models.Path{Origin: origin, Dest: destination}
						step := models.Ruta{Microbus: bus}
						step.Paraderos = []string{origin.Codigo, pair.PCode}
						step.SetDistance(store.IParaderos)
						transStep := models.Ruta{Microbus: transBus}
						transStep.Paraderos = []string{pair.PCode, transRoute.Paraderos[0]}
						transStep.SetDistance(store.IParaderos)
						path.Steps = []models.Ruta{step, transStep, transRoute}
						useful = append(useful, path)
						return useful, nil
					}
					if len(transferedRoutes) == 0 {
						for _, cpair := range unfoundTrans {
							sndTransBus := store.IBuses[cpair.BCode]
							sndTransPar := store.IParaderos[cpair.PCode]
							SndTransferedRoutes, _ := models.GetNextRoutes(&sndTransPar, &destination, &sndTransBus, &store.IParaderos, &store.IBuses)
							for _, sndTransRoute := range SndTransferedRoutes {
								path := models.Path{Origin: origin, Dest: destination}
								step := models.Ruta{Microbus: bus}
								step.Paraderos = []string{origin.Codigo, pair.PCode}
								step.SetDistance(store.IParaderos)
								transStep := models.Ruta{Microbus: transBus}
								transStep.Paraderos = []string{pair.PCode, cpair.PCode}
								transStep.SetDistance(store.IParaderos)
								sndTransStep := models.Ruta{Microbus: sndTransBus}
								sndTransStep.Paraderos = []string{cpair.PCode, sndTransRoute.Paraderos[0]}
								sndTransStep.SetDistance(store.IParaderos)
								path.Steps = []models.Ruta{step, transStep, sndTransStep, sndTransRoute}
								useful = append(useful, path)
								return useful, nil
							}
						}
					}
				}
			}
		}
	}
	return useful, nil
}
