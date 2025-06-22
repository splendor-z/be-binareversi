package model

import "time"

type Player struct {
	ID         string    `json:"id" gorm:"not null;column:id;primaryKey"`
	Name       string    `json:"name" gorm:"not null;column:name"`
	LastUsedAt time.Time `json:"lastUsedAt" gorm:"column:last_used_at;"`
}
