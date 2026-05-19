package seeders

import (
	"log"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("Seeding data...")
	
	// Panggil seeder spesifik di sini
	SeedCore(db)
	
	log.Println("Seeding completed successfully!")
}
