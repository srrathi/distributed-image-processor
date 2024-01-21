package models

type JobStatus struct {
	JobId     uint64 `gorm:"primary key;autoIncrement" json:"job_id"`
	JobStatus string `json:"job_status" validate:"required"`
}

type JobErrors struct {
	Id      uint   `gorm:"primary key;autoIncrement" json:"id"`
	JobId   uint64 `json:"job_id" validate:"required"`
	StoreId string `json:"store_id" validate:"required"`
	Error   string `json:"error" validate:"required"`
}

type JobData struct {
	JobId     int             `json:"jobId" validate:"required"`
	StoreJobs []StoreVisitData `json:"store_jobs" valiadte:"required"`
}
