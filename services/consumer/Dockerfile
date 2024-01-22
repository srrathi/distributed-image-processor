# Use the official Golang image
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Print the current working directory
RUN echo "Current working directory: $(pwd)"

# Copy only the necessary files for building
COPY go.mod go.sum ./

# Copy necessary project folders
COPY data ./data
COPY database ./database
COPY services ./services
COPY internal ./internal
COPY models ./models
COPY utils ./utils
COPY data.csv .
COPY .env .

# Print the current working directory
RUN echo "Current working directory: $(pwd)"

# Build the binary
RUN cd services/consumer && go build -o consumer && ls -l

# Expose the port the app runs on
EXPOSE 5001

# Run the binary
CMD ["./services/consumer/consumer"]
