package models

import (
	"eleventh-learn/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	FullName string    `gorm:"not null" json:"full_name" form:"full_name" valid:"required~Your full name is required"`
	Email    string    `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required"`
	Password string    `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required"`
	Products []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
	IsAdmin  bool      `gorm:"default:false" json:"is_admin"`
}

// hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	u.Password = helpers.HashPass(u.Password)

	return
}

type Admin struct {
	GormModel
	User
}

type RegularUser struct {
	GormModel
	User
}
