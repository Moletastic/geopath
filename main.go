package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Moletastic/geopath/handler"
	"github.com/Moletastic/geopath/models"
	"github.com/Moletastic/geopath/router"
	"github.com/Moletastic/geopath/store"
	"github.com/joho/godotenv"
)

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
