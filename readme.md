# Image Processing Distributed Service

Welcome to the Image Processing Distributed Service! This service allows you to process thousands of images collected from stores in a distributed and scalable manner. Follow the steps below to set up and run the service locally.

## Table of Contents
### [1. Setting up RabbitMQ](/documentation/rabbitmq.md)
### [2. Setting up PostgreSQL](/documentation/postgres.md)
### [3. Repository Local Setup](/documentation/local.md)
### [4. Endpoints Details](/documentation/endpoints.md)

## Prerequisites
### [1. Docker](https://www.docker.com/)
### [2. Golang](https://go.dev/)

## Getting Started
### 1. Clone the repository:

```bash
git clone git@github.com:srrathi/distributed-image-processor.git
```

change directory to project root directory

```bash
cd distributed-image-processor
```

### 2. Follow the steps outlined in the documentation:

#### [1. Setting up RabbitMQ](/documentation/rabbitmq.md)
#### [2. Setting up PostgreSQL](/documentation/postgres.md)
#### [3. Repository Local Setup](/documentation/local.md)
#### [4. Endpoints Details](/documentation/endpoints.md)

## Running the Microservices

1. After setting up everything, open four terminal instances in the root of the project folder.

2. Run the following commands to start the 4 microservices in 4 terminals:

```bash
go run services/jobStatus/main.go
go run services/submitJob/main.go
go run services/storeVisits/main.go
go run services/consumer/main.go
```

3. The services for the three endpoints and the image processing consumer will start working in a first-in-first-out manner.

## Testing Endpoints
Use the provided [Postman Collection](https://documenter.getpostman.com/view/14089377/2s9YymFPnz) for testing the endpoints.

## Issues and Contributions
If you encounter any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

*Happy processing!*