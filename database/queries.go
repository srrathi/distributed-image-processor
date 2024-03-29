package database

import (
	"fmt"
	"log"

	"github.com/srrathi/distributed-image-processor/models"
	"gorm.io/gorm"
)

func UpdateJobStatusDatabase(db *gorm.DB, jobId uint64, jobStatus string) error {
	// Create a new JobStatus instance with the target JobID and other data
	newJobStatus := models.JobStatus{
		JobId:     jobId,
		JobStatus: jobStatus,
	}

	// Use FirstOrCreate to perform the upsert operation based on the JobID condition
	result := db.Model(&models.JobStatus{}).Where(models.JobStatus{JobId: jobId}).Assign(models.JobStatus{JobStatus: jobStatus}).FirstOrCreate(&newJobStatus)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetStoreAreaFromStoreId(db *gorm.DB, storeId string) (string, error) {
	var record models.StoreData
	err := db.Model(&models.StoreData{}).Where("store_id = ?", storeId).First(&record).Error

	if err == gorm.ErrRecordNotFound {
		return "", nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if record.StoreArea == "" {
		return "", nil
	}
	return record.StoreArea, nil
}

func WriteErrorStoresData(db *gorm.DB, data *[]models.JobErrors) error {
	result := db.Model(&models.JobErrors{}).Create(data)

	if result.Error != nil {
		log.Println("Error performing bulk write:", result.Error)
		return result.Error
	}
	return nil
}

func WriteStoresVisitsData(db *gorm.DB, data *[]models.StoreVisits) error {
	result := db.Model(&models.StoreVisits{}).Create(data)

	if result.Error != nil {
		log.Println("Error performing bulk write:", result.Error)
		return result.Error
	}
	return nil
}

func GetJobStatusData(db *gorm.DB, jobId uint64) (*models.JobStatus, error) {
	// Fetch job status by job ID
	var jobStatusData models.JobStatus
	result := db.Model(&models.JobStatus{}).First(&jobStatusData, "job_id = ?", jobId)

	if result.Error == gorm.ErrRecordNotFound {
		// Handle case where job ID does not exist
		return nil, fmt.Errorf("job with ID %d not found in the database", jobId)
	} else if result.Error != nil {
		// Handle other errors
		return nil, result.Error
	}
	return &jobStatusData, nil
}

func GetJobErrors(db *gorm.DB, jobId uint64) ([]models.JobErrors, error) {
	var jobErrors []models.JobErrors
	result := db.Model(&models.JobErrors{}).Find(&jobErrors, "job_id=?", jobId)

	if result.Error != nil {
		// Handle errors during the query
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		// Handle case where no errors were found for the given job ID
		return nil, fmt.Errorf("no errors found for Job ID %d", jobId)
	}
	return jobErrors, nil
}

func GetStoreVisits(query *gorm.DB) ([]models.StoreVisits, error) {
	var storeVisits []models.StoreVisits
	result := query.Find(&storeVisits)
	if result.Error != nil {
		return nil, result.Error
	}
	return storeVisits, nil
}

func GetStoreInfoFromStoreId(db *gorm.DB, storeId string) (*models.StoreData, error) {
	var storeInfo models.StoreData
	err := db.Model(&models.StoreData{}).Where("store_id = ?", storeId).First(&storeInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &storeInfo, nil
}
