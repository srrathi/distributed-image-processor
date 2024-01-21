package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/srrathi/distributed-image-processor/internal"
)

var (
	JOB_FAILED    = "failed"
	JOB_COMPLETED = "completed"
	JOB_RUNNING   = "running"
	JOB_CREATED   = "created"
)

var (
	RBTMQ_QUEUE_NAME            = "jobs_schedule"
	RBTMQ_BINDING               = "jobs.create.*"
	RBTMQ_EXCHANGE              = "jobs_events"
	RBTMQ_IP_JOB_ROUTING_KEY    = "jobs.create.ip"
	RBTMQ_CONSUMER              = "image-processor"
	RBTMQ_CONCURRENT_TASK_LIMIT = 10
)

type Config struct {
	Username    string
	Password    string
	Host        string
	VirtualHost string
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getRBTMQConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	config := &Config{
		Username:    os.Getenv("RBTMQ_USERNAME"),
		Password:    os.Getenv("RBTMQ_PASSWORD"),
		Host:        os.Getenv("RBTMQ_HOST"),
		VirtualHost: os.Getenv("RBTMQ_VHOST"),
	}
	return config, nil
}

func ConnectToRBMQ() (*internal.RabbitClient, error) {
	config, err := getRBTMQConfig()
	if err != nil {
		return nil, err
	}
	conn, err := internal.ConnectRabbitMQ(config.Username, config.Password, config.Host, config.VirtualHost)
	if err != nil {
		return nil, err
	}
	mqClient, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		return nil, err
	}
	return &mqClient, nil
}
