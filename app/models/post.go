package models

import (
    "time"
)

type Post struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `json:"title"`
    Body        string    `json:"body"`
    PublishDate string `json:"publish_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}