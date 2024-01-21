package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/models"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}

	// Specify the path to your CSV file
	filename := "data.csv"
	csvFilePath := filepath.Join(".", filename)

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all CSV records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV records:", err)
	}

	// Prepare a slice to hold Store instances
	var stores []models.StoreData

	// Iterate over records and populate the slice
	for _, record := range records {
		// Create a Store instance and add it to the slice
		store := models.StoreData{
			StoreArea: record[0],
			StoreName: record[1],
			StoreId:   record[2],
		}
		stores = append(stores, store)
	}

	// Insert data in bulk
	if err := db.Create(&stores).Error; err != nil {
		log.Fatal("Error inserting records into the database:", err)
	}

	log.Println("Data imported successfully.")
}
