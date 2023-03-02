// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package user

import (
	"microblog/internal/microblog/biz"
	"microblog/internal/microblog/store"
)


type UserController struct {
    b biz.IBiz
}

// New 创建一个 user controller.
func New(ds store.IStore) *UserController {
    return &UserController{b: biz.NewBiz(ds)}
}