// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"microblog/internal/pkg/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.UserModel) error
	Get(ctx context.Context, username string) (*model.UserModel, error)
	Update(ctx context.Context, user *model.UserModel) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UserModel, error)
	Delete(ctx context.Context, username string) error
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

func (u *users) Get(ctx context.Context, username string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *users) Update(ctx context.Context, user *model.UserModel) error {
	return u.db.Save(user).Error
}

// 根據 offset 和 limit 返回 user 列表
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserModel, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserModel{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
