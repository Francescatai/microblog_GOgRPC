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

type PostStore interface {
	Create(ctx context.Context, post *model.PostModel) error
	Get(ctx context.Context, username, postID string) (*model.PostModel, error)
	Update(ctx context.Context, post *model.PostModel) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.PostModel, error)
	Delete(ctx context.Context, username string, postIDs []string) error
}

type posts struct {
	db *gorm.DB
}

var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

func (u *posts) Create(ctx context.Context, post *model.PostModel) error {
	return u.db.Create(&post).Error
}

func (u *posts) Get(ctx context.Context, username, postID string) (*model.PostModel, error) {
	var post model.PostModel
	if err := u.db.Where("username = ? and postID = ?", username, postID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (u *posts) Update(ctx context.Context, post *model.PostModel) error {
	return u.db.Save(post).Error
}

func (u *posts) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.PostModel, err error) {
	err = u.db.Where("username = ?", username).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *posts) Delete(ctx context.Context, username string, postIDs []string) error {
	err := u.db.Where("username = ? and postID in (?)", username, postIDs).Delete(&model.PostModel{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
