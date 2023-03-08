// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package user

import (
	"microblog/internal/microblog/biz"
	"microblog/internal/microblog/store"
	"microblog/pkg/auth"
	pb "microblog/pkg/proto/microblog/v1"
)

type UserController struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedMicroblogServer
}

// New 創建一個 user controller
func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
