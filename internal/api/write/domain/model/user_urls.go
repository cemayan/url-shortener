package model

import "gorm.io/gorm"

type UserUrl struct {
	gorm.Model
	UserId string `json:"userId,omitempty"`
}
