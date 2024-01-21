package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/utils"
	"gorm.io/gorm"
)

type APIResponse struct {
	Status string      `json:"status"`
	JobID  uint64      `json:"job_id"`
	Error  []ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

func main() {
	router := mux.NewRouter()
	router.Use(utils.LoggingMiddleware)

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Could not load database,", err)
	}

	router.HandleFunc("/api/status", jobStatusHandler(db)).Methods("GET")
	err = http.ListenAndServe(":5001", router)
	if err != nil {
		log.Println("There's an error with the server,", err)
	}
}

func jobStatusHandler(db *gorm.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		jobIdStr := req.URL.Query().Get("jobId")
		jobIdStr = strings.Trim(jobIdStr, " ")
		w.Header().Set("Content-Type", "application/json")

		if jobIdStr == "" {
			log.Println("invalid job id")
			w.WriteHeader(http.StatusBadRequest) // Return 400 Bad Request.
			w.Write([]byte(`{}`))
			return
		}

		// Parse the string into uint
		jobIdInt, err := strconv.ParseUint(jobIdStr, 10, 64)
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(http.StatusBadRequest) // Return 400 Bad Request.
			return
		}

		// fetching job status from database
		jobStatusData, err := database.GetJobStatusData(db, jobIdInt)
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(http.StatusBadRequest) // Return 400 Bad Request.
			return
		}

		// creating response object
		response := APIResponse{
			Status: jobStatusData.JobStatus,
			JobID:  jobStatusData.JobId,
		}
		if jobStatusData.JobStatus == utils.JOB_FAILED {
			storeErrors, err := database.GetJobErrors(db, jobIdInt)
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Add errors to the response
			for _, storeError := range storeErrors {
				errorInfo := ErrorInfo{
					StoreID: storeError.StoreId, // Replace with the actual store ID
					Error:   storeError.Error,
				}
				response.Error = append(response.Error, errorInfo)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
