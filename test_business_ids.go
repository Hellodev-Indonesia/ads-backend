package main

import (
	"fmt"
	"log"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/business"
)

func main() {
	config.LoadEnv()
	config.InitDB()

	repo := business.NewRepository(config.DB)
	ids, err := repo.GetUniqueBusinessIDsFromAdAccounts()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Found %d ids: %v\n", len(ids), ids)
}
