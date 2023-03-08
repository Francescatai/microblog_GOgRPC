// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package post

import (
	"github.com/gin-gonic/gin"

	"microblog/internal/pkg/core"
	"microblog/internal/pkg/known"
	"microblog/internal/pkg/log"
)

func (postContr *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called")

	post, err := postContr.b.Posts().Get(c, c.GetString(known.XUsernameKey), c.Param("postID"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, post)
}
