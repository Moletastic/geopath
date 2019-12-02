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

// ExtractCoord will be commented
func (p *Paradero) ExtractCoord() Coordenada {
	return Coordenada{Latitud: p.Latitud, Longitud: p.Longitud}
}

// Paraderos will be commented
type Paraderos []Paradero

// GetNextParaderos will be commented
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

// GetNearest will be commented
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

// SortByCoordDistance will be commented
func (par *Paraderos) SortByCoordDistance(location *Coordenada) {
	sort.Slice((*par)[:], func(i, j int) bool {
		l1 := (*par)[i].ExtractCoord()
		l2 := (*par)[j].ExtractCoord()
		return GetDistance(l1, *location) < GetDistance(l2, *location)
	})
}

// ToIndParaderos will be commented
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

// GetParaderosBetween will be commented
func (m MicroBus) GetParaderosBetween(origin, dest *Paradero) []string {
	originIndex := m.GetParaderoIndex(origin)
	destIndex := m.GetParaderoIndex(dest)
	if originIndex != -1 && destIndex != -1 {
		codes := m.Paraderos[originIndex:destIndex]
		return codes
	}
	return nil
}

// IsNextParadero will be commented
func (m MicroBus) IsNextParadero(origin, dest *Paradero) bool {
	nextParaderos := m.GetNextParaderosFrom(origin)
	if utils.StrListContains(nextParaderos, &dest.Codigo) {
		return true
	}
	return false
}

// MicroBuses will be commented
type MicroBuses []MicroBus

// Contains will be commented
func (ms *MicroBuses) Contains(m MicroBus) bool {
	for _, mb := range *ms {
		if mb.Recorrido == m.Recorrido {
			return true
		}
	}
	return false
}

// ToIndMicroBuses will be commented
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

type RouteCode struct {
	BCode  string
	Origin string
	Dest   string
}

type RouteCodes []RouteCode

func (rcodes *RouteCodes) Contains(rcode *RouteCode) bool {
	for _, rc := range *rcodes {
		if rc.Origin == rcode.Origin && rc.Dest == rcode.Dest && rc.BCode == rcode.BCode {
			return true
		}
	}
	return false
}

/* // CodePairs will be commented
type CodePairs []CodePair

// Contains will be commented
func (pairs *CodePairs) Contains(pair *CodePair) bool {
	for _, p := range *pairs {
		if p.PCode == pair.PCode && p.BCode == pair.BCode {
			return true
		}
	}
	return false
} */

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

// GetDistance will be commented
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

// GetNextRoutes will be commented
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

// PathResponse will be commented
type PathResponse struct {
	Data   []Path `json:"data"`
	Status int    `json:"status"`
}

// GetBuses will be commented
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

// GetParaderos will be commented
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
