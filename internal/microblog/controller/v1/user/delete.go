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

func (userContr *UserController) Delete(c *gin.Context) {
	log.C(c).Infow("Delete user function called")

	username := c.Param("name")

	if err := userContr.b.Users().Delete(c, username); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	if _, err := userContr.a.RemoveNamedPolicy("p", username, "", ""); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
