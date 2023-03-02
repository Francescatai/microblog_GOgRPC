// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package model

import (
	"time"

	"gorm.io/gorm"

	"microblog/pkg/auth"
)


type UserModel struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	Password  string    `gorm:"column:password;not null"`
	Nickname  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (u *UserModel) TableName() string {
	return "user"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
    // Encrypt the user password
    u.Password, err = auth.Encrypt(u.Password)
    if err != nil {
        return err
    }

    return nil
}