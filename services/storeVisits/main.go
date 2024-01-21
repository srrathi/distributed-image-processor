package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/models"
	"github.com/srrathi/distributed-image-processor/utils"
	"gorm.io/gorm"
)

type ErrorInfo struct {
	Error string `json:"error"`
}

type VisitData struct {
	Date      string `json:"date"`
	Perimeter uint   `json:"perimeter"`
}

// ResponseFormat represents the desired response format
type ResponseFormat struct {
	StoreID   string      `json:"store_id"`
	Area      string      `json:"area"`
	StoreName string      `json:"store_name"`
	Data      []VisitData `json:"data"`
}

func main() {
	router := mux.NewRouter()
	router.Use(utils.LoggingMiddleware)

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Could not load database,", err)
	}

	router.HandleFunc("/api/visits", storeVisitsHandler(db)).Methods("GET")
	err = http.ListenAndServe(":5002", router)
	if err != nil {
		log.Println("There's an error with the server,", err)
	}
}

func storeVisitsHandler(db *gorm.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		storeIdStr := req.URL.Query().Get("storeId")
		area := req.URL.Query().Get("area")
		startdateStr := req.URL.Query().Get("startdate")
		enddateStr := req.URL.Query().Get("enddate")

		w.Header().Set("Content-Type", "application/json")

		if storeIdStr == "" && area == "" && startdateStr == "" && enddateStr == "" {
			sendErrorResonse(w, http.StatusBadRequest, "invalid store id")
			return
		}

		var query *gorm.DB
		if storeIdStr != "" {
			query = db.Model(&models.StoreVisits{}).Where("store_id = ?", storeIdStr)
		}

		if area != "" {
			query = db.Model(&models.StoreVisits{}).Where("store_area = ?", area)
		}

		if startdateStr != "" {
			startdate, err := time.Parse(time.RFC3339, startdateStr)
			if err != nil {
				sendErrorResonse(w, http.StatusBadRequest, "invalid format for start date time, acceptable format is RFC3339 time string, "+err.Error())
				return
			}
			query = query.Where("visit_time >= ?", startdate)
		}

		if enddateStr != "" {
			enddate, err := time.Parse(time.RFC3339, enddateStr)
			if err != nil {
				sendErrorResonse(w, http.StatusBadRequest, "invalid format for start date time, acceptable format is RFC3339 time string, "+err.Error())
				return
			}
			query = query.Where("visit_time <= ?", enddate)
		}

		// Get store visits data
		storeVisits, err := database.GetStoreVisits(query)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, "internal server error,"+err.Error(), http.StatusInternalServerError)
			return
		}

		// Group visits data by store IDs
		storeVisitsData := make(map[string][]VisitData)
		for _, visit := range storeVisits {
			visitData := VisitData{
				Date:      visit.VisitTime.Format("2006-01-02"),
				Perimeter: visit.Perimeter,
			}
			storeVisitsData[visit.StoreId] = append(storeVisitsData[visit.StoreId], visitData)
		}

		// Fetch store info for each unique store ID
		storeInfoMap := make(map[string]*models.StoreData)
		for storeID := range storeVisitsData {
			storeInfo, err := database.GetStoreInfoFromStoreId(db, storeID)
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, "internal server error,"+err.Error(), http.StatusInternalServerError)
				return
			}
			storeInfoMap[storeID] = storeInfo
		}

		// Create response objects
		var response []ResponseFormat
		for storeID, visitsData := range storeVisitsData {
			storeInfo := storeInfoMap[storeID]
			if storeInfo != nil {
				resp := ResponseFormat{
					StoreID:   storeID,
					Area:      storeInfo.StoreArea,
					StoreName: storeInfo.StoreName,
					Data:      visitsData,
				}
				response = append(response, resp)
			} else {
				resp := ResponseFormat{
					StoreID:   storeID,
					Area:      "",
					StoreName: "",
					Data:      visitsData,
				}
				response = append(response, resp)
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func sendErrorResonse(w http.ResponseWriter, statusCode int, errMsg string) {
	errorResponse := ErrorInfo{
		Error: errMsg,
	}
	log.Println(errMsg)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)

}
