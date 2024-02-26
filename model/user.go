package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id         int32  `json:"id" gorm:"column:id" gorm:"primaryKey"`
	Email      string `json:"email" gorm:"column:email" gorm:"uniqueIndex"`
	Password   string `json:"password" gorm:"column:password"`
	VerifyCode string `json:"verify_code" gorm:"column:verify_code"`

	//CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	//UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	//DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "main_user"
}
