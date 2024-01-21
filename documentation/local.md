# Repository Local Setup
## 3. Repository Local Setup
### **3.1 Cloning the Repository**

Clone this repository on your local machine:

```bash
git clone git@github.com:srrathi/distributed-image-processor.git
```

change directory to project root directory

```bash
cd distributed-image-processor
```

### **3.2 Environment File**
In the root of the repository, create a file named .env and add the following credentials based on the settings from the previous steps:

```env
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=12345678
DB_USER=srrathi
DB_DATABASE=ip_jobs
DB_SSLMODE=disable
RBTMQ_USERNAME=srrathi
RBTMQ_PASSWORD=12345678
RBTMQ_HOST=localhost:5672
RBTMQ_VHOST=jobs
```

### **3.3 Install Dependencies**
In the root of the project folder, where the go.mod file exists, run the following command to download all project dependencies:

```bash
go mod download
```

### **3.4 Database Setup**
Dump the storemaster CSV data into the database. Open the terminal in the root of the project and run the following command:
```bash
go run data/dump.go
```

If successful, you should see "Connected to postgres" and "Data imported successfully" in the terminal.

### **3.5 Running Microservices**
Open four terminal instances in the root of the project folder and run the following commands to start the microservices:

1. Job Status Service:
```bash
go run services/jobStatus/main.go
```

2. Submit Job Service:
```bash
go run services/submitJob/main.go
```

3. Store Visits Service:
```bash
go run services/storeVisits/main.go
```

4. Image Processing Consumer:
```bash
go run services/consumer/main.go
```

This will start the services for the three endpoints and the image processing consumer, which will consume jobs pushed into the RabbitMQ queue.

Now, the microservices are set up, and you can test them out. If you encounter any issues, ensure that you have the correct environment variables, have successfully connected to the database, and have the necessary dependencies installed.