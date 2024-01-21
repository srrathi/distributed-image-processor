package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/models"
	"github.com/srrathi/distributed-image-processor/utils"
	"gorm.io/gorm"
)

type RequestBody struct {
	Count  int                     `json:"count" validate:"required"`
	Visits []models.StoreVisitData `json:"visits" validate:"required"`
}

type IError struct {
	Field string
	Tag   string
	Value string
}

type ErrorInfo struct {
	Error string `json:"error"`
}

type SuccessInfo struct {
	JobId int `json:"job_id"`
}

var Validator = validator.New()

func main() {
	router := mux.NewRouter()
	router.Use(utils.LoggingMiddleware)

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Could not load database,", err)
	}

	router.HandleFunc("/api/submit", submitJobHandler(db)).Methods("POST")
	err = http.ListenAndServe(":5003", router)
	if err != nil {
		log.Println("There's an error with the server,", err)
	}
}

func submitJobHandler(db *gorm.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := new(RequestBody)
		err := json.NewDecoder(req.Body).Decode(data)
		if err != nil {
			log.Println("There was an error decoding the request body into the struct")
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		// Validate request body
		err = Validator.Struct(data)
		if err != nil {
			handleValidationError(err, w)
			return
		}

		if data.Count != len(data.Visits) {
			handleError(w, http.StatusBadRequest, errors.New("count and number of objects in visits array should be equal"))
			return
		}

		// generate a job ID
		jobId := generateUniqueIntegerID(7)
		jobData := models.JobData{
			JobId:     jobId,
			StoreJobs: data.Visits,
		}

		// send data to exchanger
		err = sendDataToRBMQExchanger(jobData)
		if err != nil {
			log.Println(err.Error())
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		// update status in database for jobs
		err = database.UpdateJobStatusDatabase(db, uint64(jobId), utils.JOB_CREATED)
		if err != nil {
			log.Println(err.Error())
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		// return created job response
		successJobResponse := SuccessInfo{
			JobId: jobId,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(successJobResponse)
	}
}

func handleValidationError(validationError error, w http.ResponseWriter) {
	var errors []*IError
	for _, err := range validationError.(validator.ValidationErrors) {
		var el IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors = append(errors, &el)
	}
	jsonStr, err := json.Marshal(errors)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errorResponse := ErrorInfo{
		Error: string(jsonStr),
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

func generateUniqueIntegerID(digits int) int {
	if digits < 0 {
		digits = 7
	}
	// Calculate the range based on the number of digits
	minValue := int(math.Pow10(digits - 1))
	maxValue := int(math.Pow10(digits)) - 1

	// Generate a random integer within the specified range
	return rand.Intn(maxValue-minValue+1) + minValue
}

func sendDataToRBMQExchanger(data models.JobData) error {
	client, err := utils.ConnectToRBMQ()
	if err != nil {
		return err
	}

	err = client.CreateQueue(utils.RBTMQ_QUEUE_NAME, true, false)
	if err != nil {
		return err
	}

	// Create binding between the jobs_events exchange and the customers-created queue
	err = client.CreateBinding(utils.RBTMQ_QUEUE_NAME, utils.RBTMQ_BINDING, utils.RBTMQ_EXCHANGE)
	if err != nil {
		return err
	}

	// Create context to manage timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dataStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = client.Send(ctx, utils.RBTMQ_EXCHANGE, utils.RBTMQ_IP_JOB_ROUTING_KEY, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent, // This tells rabbitMQ that this message should be Saved if no resources accepts it before a restart (durable)
		Body:         dataStr,
	})
	if err != nil {
		return err
	}

	return nil
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorInfo{
		Error: err.Error(),
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
