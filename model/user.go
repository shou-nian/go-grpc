package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id         int32  `json:"id" gorm:"column:id" gorm:"primaryKey"`
	Email      string `json:"email" gorm:"column:email" gorm:"unique"`
	Password   string `json:"password" gorm:"column:password"`
	VerifyCode string `json:"verify_code" gorm:"column:verify_code"`

	Token Token `json:"token,omitempty" gorm:"foreignKey:UserId;references:Id"`
}

type Token struct {
	gorm.Model
	Id     int32     `json:"id" gorm:"column:id" gorm:"primaryKey"`
	UserId int32     `json:"user_id" gorm:"column:user_id" gorm:"unique"`
	Token  string    `json:"token" gorm:"column:token"`
	InDate time.Time `json:"in_date" gorm:"column:in_date"`
}

func (User) TableName() string { return "main_user" }

func (Token) TableName() string { return "main_token" }
