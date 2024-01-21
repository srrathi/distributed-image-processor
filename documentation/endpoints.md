# Endpoints Details
## 4. Endpoints Details
### **4.1 Submit Job**
Create a job to process images collected from stores. The visit_time should be in RFC3339 or ISO format.

- **URL:** http://localhost:5003/api/submit
- **Method:** POST
- **Request Payload:**
```json
{
  "count": 2,
  "visits": [
    {
      "store_id": "S00339218",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/2.jpg",
        "https://www.gstatic.com/webp/gallery/3.jpg"
      ],
      "visit_time": "2024-01-21T16:23:40.898Z"
    },
    {
      "store_id": "S01408764",
      "image_url": ["https://www.gstatic.com/webp/gallery/3.jpg"],
      "visit_time": "2024-01-21T16:23:40.898Z"
    }
  ]
}
```

- **Success Response:**
- **Code:** 201 CREATED
- **Content Example:**

```json
{
  "job_id": 123
}
```

- **Error Responses:**
- **Code:** 400 BAD REQUEST
- **Content Example:**
```json
{
  "error": ""
}
```

### **4.2 Get Job Info**
- **URL:** http://localhost:5001/api/status?jobId=3059701
- **URL Parameters:**
- **jobId:** Job ID received while creating the job
- **Method:** GET
- **Success Response:**
- **Code: 200 OK**
- **Content Example:**

```json
{
  "status": "completed",
  "job_id": ""
}
```

- **Job Status:** failed
```json
{
  "status": "failed",
  "job_id": "",
  "error": [
    {
      "store_id": "S00339218",
      "error": ""
    }
  ]
}
```

- **Error Responses:**
- **Code: 400 BAD REQUEST**
- **Content:**

```json
{}
```

### **4.3 Show Visit Info**
- **URL:** http://localhost:5002/api/visits?area=abc&storeid=S00339218&startdate=stdate&enddate=endate
- **URL Parameters:**
- **area:** Area code from Store Master
- **storeid:** Store ID
- **startdate / enddate:** Date in RFC3339 format to filter data based on the store visit_time
- **Method:** GET
- **Success Response:**
- **Code:** 200 OK
- **Content Example:**

```json
{
  "results": [
    {
      "store_id": "S00339218",
      "area": "",
      "store_name": "",
      "data": [
        {
          "date": "",
          "perimeter": ""
        },
        {
          "date": "",
          "perimeter": ""
        }
      ]
    },
    {
      "store_id": "S01408764",
      "area": "",
      "store_name": "",
      "data": [
        {
          "date": "",
          "perimeter": ""
        },
        {
          "date": "",
          "perimeter": ""
        }
      ]
    }
  ]
}
```

- **Error Responses:**
- **Condition:** If area or storeid does not exist
- **Code:** 400 BAD REQUEST
- **Content:**

```json
{
  "error": ""
}
```
This concludes the detailed information about the three endpoints. You can use these details to interact with the services and test the functionality.