package Models

import "time"
import "gorm.io/plugin/soft_delete"

type User struct {
	Id        string `gorm:"type:varchar(100);primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Location  string `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
