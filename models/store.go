package models

import "time"

type StoreData struct {
	Id        uint   `gorm:"primary key;autoIncrement" json:"id"`
	StoreId   string `json:"store_id" validate:"required"`
	StoreArea string `json:"store_area" validate:"required"`
	StoreName string `json:"store_name" validate:"required"`
}

type StoreVisits struct {
	Id        uint      `gorm:"primary key;autoIncrement" json:"id"`
	StoreId   string    `json:"store_id" validate:"required"`
	StoreArea string    `json:"store_area" validate:"required"`
	Perimeter uint      `json:"perimeter" validate:"required"`
	VisitTime time.Time `gorm:"type:time" json:"visit_time" validate:"required"`
}

type StoreVisitData struct {
	StoreId   string    `json:"store_id" validate:"required"`
	VisitTime time.Time `json:"visit_time"`
	ImageUrl  []string  `json:"image_url"`
}