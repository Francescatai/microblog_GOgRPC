// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	
	S *datastore
)


type IStore interface {
	Users() UserStore
}


type datastore struct {
	db *gorm.DB
}


var _ IStore = (*datastore)(nil)


func NewStore(db *gorm.DB) *datastore {
	// 確保 S 只被初始化一次
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}


func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}
