package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/srrathi/distributed-image-processor/database"
	"github.com/srrathi/distributed-image-processor/models"
	"github.com/srrathi/distributed-image-processor/services/consumer/processing"
	"github.com/srrathi/distributed-image-processor/utils"
	"golang.org/x/sync/errgroup"
)

func main() {
	mqClient, err := utils.ConnectToRBMQ()
	if err != nil {
		panic(err)
	}

	messageBus, err := mqClient.Consume(utils.RBTMQ_QUEUE_NAME, utils.RBTMQ_CONSUMER, false)
	if err != nil {
		panic(err)
	}

	// To connect to database
	db, err := database.NewConnection()
	if err != nil {
		panic(err)
	}

	// blocking is used to block forever
	var blocking chan struct{}

	// Set a timeout for 15 secs
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create an Errgroup to manage concurrecy
	g, _ := errgroup.WithContext(ctx)
	// Set amount of concurrent tasks
	g.SetLimit(utils.RBTMQ_CONCURRENT_TASK_LIMIT)
	go func() {
		for message := range messageBus {
			// Spawn a worker
			msg := message
			g.Go(func() error {
				// Unmarshal the JSON data into the struct
				var jobData models.JobData
				err := json.Unmarshal(msg.Body, &jobData)
				if err != nil {
					log.Println("Error:", err)
					return err
				}

				// Update job status to running
				err = database.UpdateJobStatusDatabase(db, uint64(jobData.JobId), utils.JOB_RUNNING)
				if err != nil {
					log.Println("Error:", err)
					return err
				}
				// Process images
				err = processing.ProcessStoreVisits(jobData, db)
				if err != nil {
					log.Println("Error:", err)
					return err
				}

				// Multiple means that we acknowledge a batch of messages, leave false for now
				if err := msg.Ack(false); err != nil {
					log.Printf("Acknowledged message failed: Retry ? Handle manually %s\n", msg.MessageId)
					return err
				}

				log.Println("Acknowledged message: ", msg.MessageId)
				return nil
			})
		}
	}()

	log.Println("Consuming, to close the program press CTRL+C")
	// This will block forever
	<-blocking

}
