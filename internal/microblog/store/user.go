// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package store

import (
    "context"

    "gorm.io/gorm"

    "microblog/internal/pkg/model"
)


type UserStore interface {
    Create(ctx context.Context, user *model.UserModel) error
}


type users struct {
    db *gorm.DB
}


var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
    return &users{db}
}


func (u *users) Create(ctx context.Context, user *model.UserModel) error {
    return u.db.Create(&user).Error
}