package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	"github.com/Moletastic/geopath/utils"
)

// Paradero almacena la información de un paradero,
// además del código de aquellos buses disponibles en este
type Paradero struct {
	Codigo     string   `json:"codigo"`
	Direccion  string   `json:"direccion"`
	Latitud    float64  `json:"latitud"`
	Longitud   float64  `json:"longitud"`
	Microbuses []string `json:"microbuses"`
}

// ExtractCoord retorna una Coordenada apartir de
// la latitud y longitud de ubicación del paradero
func (p *Paradero) ExtractCoord() Coordenada {
	return Coordenada{Latitud: p.Latitud, Longitud: p.Longitud}
}

// Paraderos representa un array de Paradero
type Paraderos []Paradero

// GetNextParaderos retorna los paraderos más cercanos a una coordenada
func (par *Paraderos) GetNextParaderos(location *Coordenada) []Paradero {
	ps := make([]Paradero, 0)
	for _, paradero := range *par {
		l := (&paradero).ExtractCoord()
		if location.EqualsOrClose(l, 0.3) {
			ps = append(ps, paradero)
		}
	}
	return ps
}

// GetNearest obtiene el paradero más cercano a una ubicación
func (par *Paraderos) GetNearest(location Coordenada) Paradero {
	var p Paradero
	d := 999.0
	for _, paradero := range *par {
		l := Coordenada{Latitud: paradero.Latitud, Longitud: paradero.Longitud}
		if GetDistance(l, location) < d {
			d = GetDistance(l, location)
			p = paradero
		}
	}
	return p
}

// SortByCoordDistance ordena los paraderos almacenados
// de acuerdo a la cercanía de estos respecto a una ubicación
func (par *Paraderos) SortByCoordDistance(location *Coordenada) {
	sort.Slice((*par)[:], func(i, j int) bool {
		l1 := (*par)[i].ExtractCoord()
		l2 := (*par)[j].ExtractCoord()
		return GetDistance(l1, *location) < GetDistance(l2, *location)
	})
}

// ToIndParaderos retorna un mapa indexado de Paraderos por su código
func (par *Paraderos) ToIndParaderos() IndParaderos {
	indexed := make(map[string]Paradero, 0)
	for _, paradero := range *par {
		indexed[paradero.Codigo] = paradero
	}
	return indexed
}

// IndParaderos es un mapa de estructura Paradero
// generalmente indexado por el campo Codigo
type IndParaderos map[string]Paradero

// MicroBus almacena la información de un microbus,
// además del código de aquellos paraderos por los que pasa
type MicroBus struct {
	Recorrido string   `json:"recorrido"`
	Tipo      string   `json:"tipo"`
	Paraderos []string `json:"paraderos"`
}

// GetParaderoIndex retorna el índice del paradero de acuerdo
// al recorrido de un microbus. Si este microbus no pasa por dicho
// paradero retorna -1
func (m *MicroBus) GetParaderoIndex(paradero *Paradero) int {
	for index, p := range m.Paraderos {
		if p == paradero.Codigo {
			return index
		}
	}
	return -1
}

// GetNextParaderosFrom retorna los paraderos del recorrido del microbus
// desde un paradero entregado, sin contarlo
func (m MicroBus) GetNextParaderosFrom(paradero *Paradero) []string {
	currentIndex := m.GetParaderoIndex(paradero)
	if currentIndex != -1 {
		codes := m.Paraderos[currentIndex+1:]
		return codes
	}
	return nil
}

// GetParaderosBetween retorna los códigos de paraderos
// del recorrido del MicroBus entre un paradero de origen y destino
func (m MicroBus) GetParaderosBetween(origin, dest *Paradero) []string {
	originIndex := m.GetParaderoIndex(origin)
	destIndex := m.GetParaderoIndex(dest)
	if originIndex != -1 && destIndex != -1 {
		codes := m.Paraderos[originIndex:destIndex]
		return codes
	}
	return nil
}

// IsNextParadero retorna sí un paradero destino, se encuentra
// en el recorrido del MicroBus luego de un paradero de origen
func (m MicroBus) IsNextParadero(origin, dest *Paradero) bool {
	nextParaderos := m.GetNextParaderosFrom(origin)
	if utils.StrListContains(nextParaderos, &dest.Codigo) {
		return true
	}
	return false
}

// MicroBuses representa un array de MicroBus
type MicroBuses []MicroBus

// Contains retorna sí un MicroBus se encuentra almacenado
func (ms *MicroBuses) Contains(m MicroBus) bool {
	for _, mb := range *ms {
		if mb.Recorrido == m.Recorrido {
			return true
		}
	}
	return false
}

// ToIndMicroBuses retorna un mapa de microbuses indexado por Código
func (ms *MicroBuses) ToIndMicroBuses() IndMicroBuses {
	indexed := make(IndMicroBuses, 0)
	for _, micro := range *ms {
		indexed[micro.Recorrido] = micro
	}
	return indexed
}

// IndMicroBuses es un mapa de estructura MicroBus
// generalmente indexado por el campo Recorrido
type IndMicroBuses map[string]MicroBus

// Ruta almacena el intervalo del recorrido de un microbus
type Ruta struct {
	Paraderos []string `json:"paraderos"`
	Microbus  MicroBus `json:"microbus"`
	Distancia float64  `json:"distancia"`
}

// SetDistance fija la distancia de ruta de acuerdo a
// la distancia entre los paraderos de origen y destino
func (r *Ruta) SetDistance(paraderos IndParaderos) {
	origin := paraderos[r.Paraderos[0]]
	dest := paraderos[r.Paraderos[1]]
	r.Distancia = GetDistance(origin.ExtractCoord(), dest.ExtractCoord())
}

// CodePair almacena el código de un microbus y un paradero
type CodePair struct {
	BCode string
	PCode string
}

// RouteCode almacena el código de un microbus, y de paraderos de origen y destino
type RouteCode struct {
	BCode  string
	Origin string
	Dest   string
}

// RouteCodes representa un array de RouteCode
type RouteCodes []RouteCode

// Contains retorna sí un RouteCode se encuentra en el array
func (rcodes *RouteCodes) Contains(rcode *RouteCode) bool {
	for _, rc := range *rcodes {
		if rc.Origin == rcode.Origin && rc.Dest == rcode.Dest && rc.BCode == rcode.BCode {
			return true
		}
	}
	return false
}

// Path almacena los pasos a seguir para viajar
// desde un paradero origen a uno destino
type Path struct {
	Origin Paradero `json:"origin"`
	Dest   Paradero `json:"dest"`
	Steps  []Ruta   `json:"steps"`
}

// HasMicroBus revisa si existe un MicroBus en el recorrido
// construido
func (p *Path) HasMicroBus(m MicroBus) bool {
	for _, step := range p.Steps {
		if step.Microbus.Recorrido == m.Recorrido {
			return true
		}
	}
	return false
}

// GetDistance retorna la distancia total de todos los paraderos
// en los que se realiza un trasbordo o fin de ruta
func (p *Path) GetDistance(paraderos IndParaderos) float64 {
	totalDistance := 0.0
	for _, step := range p.Steps {
		if step.Distancia == 0 {
			step.SetDistance(paraderos)
		}
		distance := step.Distancia
		totalDistance += distance
	}
	return totalDistance
}

// GetNextRoutes retorna las posibles rutas de a un destino, basado
// en el recorrido de todos los otros buses, que pueden ser tomados
// en los próximos paraderos. Además retorna los pares de código,
// de los buses que no pasan por el destino, y el paradero en el que se tomó
func GetNextRoutes(origin, dest *Paradero, bus *MicroBus, parades *IndParaderos, buses *IndMicroBuses) ([]Ruta, []CodePair) {
	nextParades := bus.GetNextParaderosFrom(origin)
	unfound := make([]CodePair, 0)
	routes := make([]Ruta, 0)
	for _, pcode := range nextParades {
		parada := (*parades)[pcode]
		for _, bcode := range parada.Microbuses {
			microbus := (*buses)[bcode]
			if microbus.GetParaderoIndex(dest) != -1 {
				r := Ruta{}
				r.Microbus = microbus
				r.Paraderos = append(r.Paraderos, parada.Codigo)
				r.Paraderos = append(r.Paraderos, dest.Codigo)
				r.SetDistance(*parades)
				routes = append(routes, r)
			} else {
				c := CodePair{BCode: bcode, PCode: pcode}
				unfound = append(unfound, c)
			}
		}
	}
	return routes, unfound
}

// Paths representa un array de Path
type Paths []Path

// GetBest retorna el Path con menor distancia total entre
// todos sus path
func (ps Paths) GetBest(parades *IndParaderos) *Path {
	var best Path
	minD := 999.0
	for _, path := range ps {
		distance := path.GetDistance(*parades)
		if distance < minD {
			minD = distance
			best = path
		}
	}
	return &best
}

// PathResponse contiene los datos de una respuesta
// a peticiones de un Path
type PathResponse struct {
	Data   []Path `json:"data"`
	Status int    `json:"status"`
}

// GetBuses retorna una estructura MicroBuses con los datos
// contenidos en el archivo con el nombre fijado, o un error
func GetBuses(filename string) (MicroBuses, error) {
	buses := make(MicroBuses, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &buses)
	if err != nil {
		return nil, err
	}
	return buses, nil
}

// GetParaderos retorna una estructura Paraderos con los datos
// contenidos en el archivo con el nombre fijado, o un error
func GetParaderos(filename string) (Paraderos, error) {
	paraderos := make(Paraderos, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	p, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(p, &paraderos)
	if err != nil {
		return nil, err
	}
	return paraderos, nil
}
