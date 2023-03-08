// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package post

import (
	"microblog/internal/microblog/biz"
	"microblog/internal/microblog/store"
)

type PostController struct {
	b biz.IBiz
}

func New(ds store.IStore) *PostController {
	return &PostController{b: biz.NewBiz(ds)}
}
