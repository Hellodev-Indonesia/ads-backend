package main

import (
	"log"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/database/seeders"
)

func main() {
	config.LoadEnv()

	config.InitDB()

	seeders.Run(config.DB)

	log.Println("Database successfully seeded!")
}
