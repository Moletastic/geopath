package main

import (
	"fmt"
	"log"
	"os"
	_"github.com/Moletastic/geopath/docs"
	"github.com/Moletastic/geopath/handler"
	"github.com/Moletastic/geopath/models"
	"github.com/Moletastic/geopath/router"
	"github.com/Moletastic/geopath/store"
	"github.com/joho/godotenv"
)

// @title Proyecto Computación Distribuida
// @version 1.0
// @description Esta API entregará el recorrido en micro (con menos transbordos posibles) para dos coordenadas de inicio y destino
// @host ec2-3-136-84-231.us-east-2.compute.amazonaws.com
// @BasePath /api
// @tag.name Geopaths

func main() {
	r := router.New()

	v4 := r.Group("/api")

	err := godotenv.Load()
	if err != nil {
		r.Logger.Fatal(err)
	}

	buses, err := models.GetBuses("data/microbuses.json")
	if err != nil {
		log.Fatal(err)
	}

	paraderos, err := models.GetParaderos("data/paradas.json")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewPathStore(paraderos, buses)
	h := handler.NewHandler(store)
	h.Register(v4)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == ":" {
		port = ":8080"
	}
	r.Logger.Fatal(r.Start(port))
}
