// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package model

import (
	"time"

	"gorm.io/gorm"

	"microblog/pkg/util/id"
)

type PostModel struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	PostID    string    `gorm:"column:postID;not null"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (p *PostModel) TableName() string {
	return "post"
}

// BeforeCreate 在创建数据库记录之前生成 postID.
func (p *PostModel) BeforeCreate(tx *gorm.DB) error {
	p.PostID = "post-" + id.GenShortID()

	return nil
}
