package processing

import (
	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/models"
	"github.com/srrathi/distributed-image-processor/utils"
	"gorm.io/gorm"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"sync"
)

// ImageData represents the structure of image data fetched from the internet
type ImageData struct {
	URL       string
	Perimeter int
	Error     error
}

func ProcessStoreVisits(jobData models.JobData, db *gorm.DB) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorResults []models.JobErrors
	var successResults []models.StoreVisits

	for _, visit := range jobData.StoreJobs {
		wg.Add(1)

		go func(visit models.StoreVisitData) {
			defer wg.Done()

			// Fetch images concurrently
			imageData := fetchImages(visit.ImageUrl)

			// Calculate total perimeter for the store visit
			perimeterSum := calculatePerimeterSum(imageData)

			// Update the result
			imageErr := getErrorString(imageData)
			if imageErr != "" {
				jobError := models.JobErrors{
					JobId:   uint64(jobData.JobId),
					StoreId: visit.StoreId,
					Error:   imageErr,
				}
				mu.Lock()
				errorResults = append(errorResults, jobError)
				mu.Unlock()
			} else {
				// Fetch store area from database
				storeArea, _ := database.GetStoreAreaFromStoreId(db, visit.StoreId)
				visitData := models.StoreVisits{
					StoreId:   visit.StoreId,
					StoreArea: storeArea,
					Perimeter: uint(perimeterSum),
					VisitTime: visit.VisitTime,
				}
				mu.Lock()
				successResults = append(successResults, visitData)
				mu.Unlock()
			}
		}(visit)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if len(errorResults) > 0 {
		err := database.WriteErrorStoresData(db, &errorResults)
		if err != nil {
			return err
		}

		err = database.UpdateJobStatusDatabase(db, uint64(jobData.JobId), utils.JOB_FAILED)
		if err != nil {
			return err
		}
	} else {
		err := database.UpdateJobStatusDatabase(db, uint64(jobData.JobId), utils.JOB_COMPLETED)
		if err != nil {
			return err
		}
	}

	if len(successResults) > 0 {
		err := database.WriteStoresVisitsData(db, &successResults)
		if err != nil {
			return err
		}
	}
	return nil
}

func fetchImages(imageURLs []string) []ImageData {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var imageDataArray []ImageData

	for _, url := range imageURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			// Fetch image data
			imageData, err := fetchImage(url)

			// Append the result to the imageDataArray
			mu.Lock()
			imageDataArray = append(imageDataArray, ImageData{URL: url, Perimeter: func() int {
				if imageData != nil {
					return *imageData
				}
				return 0
			}(), Error: err})
			mu.Unlock()
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return imageDataArray
}

func fetchImage(url string) (*int, error) {
	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching the image:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the image
	config, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		log.Println("Error determining image format:", err)
		return nil, err
	}

	// Calculate the perimeter (twice the sum of width and height)
	perimeter := 2 * (config.Width + config.Height)
	return &perimeter, nil
}

func calculatePerimeterSum(imageDataArray []ImageData) int {
	var perimeterSum int
	for _, imageData := range imageDataArray {
		if imageData.Perimeter > 0 {
			perimeterSum += imageData.Perimeter
		}
	}
	return perimeterSum
}

func getErrorString(imageDataArray []ImageData) string {
	for _, imageData := range imageDataArray {
		if imageData.Error != nil {
			return imageData.Error.Error()
		}
	}
	return ""
}
