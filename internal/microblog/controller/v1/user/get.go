// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package user

import (
	"github.com/gin-gonic/gin"

	"microblog/internal/pkg/core"
	"microblog/internal/pkg/log"
)

func (userContr *UserController) Get(c *gin.Context) {
	log.C(c).Infow("Get user function called")

	user, err := userContr.b.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
