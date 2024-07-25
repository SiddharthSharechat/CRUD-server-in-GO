package Models

import "time"
import "gorm.io/plugin/soft_delete"

type User struct {
	Id        string
	Name      string
	Email     string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
