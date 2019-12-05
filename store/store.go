package store

import (
	"errors"
	"sync"

	"github.com/Moletastic/geopath/models"
)

// PathStore almacena los datos de paraderos y buses
type PathStore struct {
	Paraderos  models.Paraderos
	IParaderos models.IndParaderos
	Buses      models.MicroBuses
	IBuses     models.IndMicroBuses
}

// NewPathStore crea un nuevo PathStore
func NewPathStore(paraderos models.Paraderos, buses models.MicroBuses) *PathStore {
	return &PathStore{
		Paraderos:  paraderos,
		Buses:      buses,
		IParaderos: paraderos.ToIndParaderos(),
		IBuses:     buses.ToIndMicroBuses(),
	}
}

// GetPathToDest retorna el Path con menor cantidad de trasbordos, o un error
func (store *PathStore) GetPathToDest(origin, dest *models.Coordenada) ([]models.Path, error) {
	store.Paraderos.SortByCoordDistance(origin)
	originParades := store.Paraderos.GetNextParaderos(origin)
	store.Paraderos.SortByCoordDistance(dest)
	destParades := store.Paraderos.GetNextParaderos(dest)
	if len(originParades) == 0 || len(destParades) == 0 {
		return nil, errors.New("No existen paradas cercanas")
	}
	useful := make(models.Paths, 0)
	noavalaible := make(models.RouteCodes, 0)
	// Buscando Path sin trasbordo
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
		var jobs sync.WaitGroup
		jobs.Add(len(noavalaible))
		sndnoavalaible := make([]models.RouteCodes, 0)
		thdnoavalaible := make([]models.RouteCodes, 0)
		// Buscando Path con 1 trasbordo
		for _, rcode := range noavalaible {
			go func(rcode models.RouteCode) {
				bus := store.IBuses[rcode.BCode]
				origin := store.IParaderos[rcode.Origin]
				destination := store.IParaderos[rcode.Dest]
				routes, unfound := models.GetNextRoutes(&origin, &destination, &bus, &store.IParaderos, &store.IBuses)
				if len(routes) != 0 {
					for _, route := range routes {
						if route.Microbus.Recorrido != bus.Recorrido {
							path := models.Path{Origin: origin, Dest: destination}
							step := models.Ruta{Microbus: bus}
							step.Paraderos = []string{origin.Codigo, route.Paraderos[0]}
							step.SetDistance(store.IParaderos)
							path.Steps = append(path.Steps, step)
							path.Steps = append(path.Steps, route)
							useful = append(useful, path)
						}
					}
				} else {
					unavalaible := make(models.RouteCodes, 0)
					for _, codepair := range unfound {
						if codepair.BCode != bus.Recorrido {
							transBus := store.IBuses[codepair.BCode]
							transParade := store.IParaderos[codepair.PCode]
							originRC := models.RouteCode{BCode: bus.Recorrido, Origin: rcode.Origin, Dest: transParade.Codigo}
							transRC := models.RouteCode{BCode: transBus.Recorrido, Origin: transParade.Codigo, Dest: rcode.Dest}
							unavalaible = models.RouteCodes{originRC, transRC}
						}
					}
					if len(unavalaible) != 0 {
						sndnoavalaible = append(sndnoavalaible, unavalaible)
					}
				}
				jobs.Done()
			}(rcode)
		}
		jobs.Wait()
		if len(useful) != 0 {
			return []models.Path{*useful.GetBest(&store.IParaderos)}, nil
		}
		noavalaible = nil
		jobs.Add(len(sndnoavalaible))
		// Buscando Path con 2 trasbordos
		for _, routecodes := range sndnoavalaible {
			go func(routecodes models.RouteCodes) {
				index := len(routecodes) - 1
				bus := store.IBuses[routecodes[index].BCode]
				origin := store.IParaderos[routecodes[index].Origin]
				for _, dest := range destParades {
					routes, unfound := models.GetNextRoutes(&origin, &dest, &bus, &store.IParaderos, &store.IBuses)
					if len(routes) != 0 {
						for _, route := range routes {
							path := models.Path{Origin: origin, Dest: dest}
							step := models.Ruta{Microbus: store.IBuses[routecodes[0].BCode]}
							step.Paraderos = []string{routecodes[0].Origin, routecodes[0].Dest}
							step.SetDistance(store.IParaderos)
							transtep := models.Ruta{Microbus: store.IBuses[routecodes[1].BCode]}
							transtep.Paraderos = []string{routecodes[1].Origin, route.Paraderos[0]}
							transtep.SetDistance(store.IParaderos)
							path.Steps = []models.Ruta{step, transtep, route}
							useful = append(useful, path)
						}
					} else {
						unavalaible := make(models.RouteCodes, 0)
						for _, codepair := range unfound {
							unavalaible = routecodes
							route := models.RouteCode{BCode: codepair.BCode, Origin: codepair.PCode, Dest: dest.Codigo}
							unavalaible[1].Dest = route.Origin
							unavalaible = append(unavalaible, route)
						}
						if len(unavalaible) != 0 {
							thdnoavalaible = append(thdnoavalaible, unavalaible)
						}
					}
				}
				jobs.Done()
			}(routecodes)
		}
		jobs.Wait()
		if len(useful) != 0 {
			return []models.Path{*useful.GetBest(&store.IParaderos)}, nil
		}
		// Buscando Path con 3 trasbordos
		jobs.Add(len(thdnoavalaible))
		for _, routecodes := range thdnoavalaible {
			go func(routecodes models.RouteCodes) {
				index := len(routecodes) - 1
				bus := store.IBuses[routecodes[index].BCode]
				origin := store.IParaderos[routecodes[index].Origin]
				for _, dest := range destParades {
					routes, _ := models.GetNextRoutes(&origin, &dest, &bus, &store.IParaderos, &store.IBuses)
					if len(routes) != 0 {
						for _, route := range routes {
							path := models.Path{Origin: origin, Dest: dest}
							step := models.Ruta{Microbus: store.IBuses[routecodes[0].BCode]}
							step.Paraderos = []string{routecodes[0].Origin, routecodes[0].Dest}
							step.SetDistance(store.IParaderos)
							transtep := models.Ruta{Microbus: store.IBuses[routecodes[1].BCode]}
							transtep.Paraderos = []string{routecodes[1].Origin, routecodes[2].Origin}
							transtep.SetDistance(store.IParaderos)
							sndtranstep := models.Ruta{Microbus: store.IBuses[routecodes[2].BCode]}
							sndtranstep.Paraderos = []string{routecodes[2].Origin, route.Paraderos[0]}
							sndtranstep.SetDistance(store.IParaderos)
							path.Steps = []models.Ruta{step, transtep, sndtranstep, route}
							useful = append(useful, path)
						}
					}
				}
				jobs.Done()
			}(routecodes)
		}
		jobs.Wait()
		if len(useful) != 0 {
			return []models.Path{*useful.GetBest(&store.IParaderos)}, nil
		}
	}
	return useful, nil
}

// GetParadeByID retorna el paradero con el id entregado
func (store *PathStore) GetParadeByID(id string) (*models.Paradero, error) {
	paradero := store.IParaderos[id]
	if paradero.Codigo != "" {
		return &paradero, nil
	}
	return nil, errors.New("Paradero no encontrado")
}

// GetMicroBusByID retorna el microbus con el id entregado
func (store *PathStore) GetMicroBusByID(id string) (*models.MicroBus, error) {
	bus := store.IBuses[id]
	if bus.Recorrido != "" {
		return &bus, nil
	}
	return nil, errors.New("MicroBus no encontrado")
}
