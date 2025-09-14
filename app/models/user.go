package models

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/goravel/framework/database/orm"
)

type User struct {
    orm.Model
    Name     string `json:"name" gorm:"not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"-" gorm:"not null"` // Hide password in JSON
    Active   int    `gorm:"column:active" json:"active"`
}

func (User) TableName() string {
    return "users"
}

// HashPassword untuk hash password sebelum disimpan
func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// CheckPassword untuk verify password
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}