// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package microblog

import (
	"github.com/gin-gonic/gin"

	"microblog/internal/microblog/controller/v1/user"
	"microblog/internal/microblog/store"
	"microblog/internal/pkg/core"
	errno "microblog/internal/pkg/err"
	"microblog/internal/pkg/log"
)


func installRouters(g *gin.Engine) error {

	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)

	v1 := g.Group("/v1")
	{

		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
		}
	}

	return nil
}
