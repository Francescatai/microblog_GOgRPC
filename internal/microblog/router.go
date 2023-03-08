// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package microblog

import (
	"github.com/gin-gonic/gin"

	"microblog/internal/microblog/controller/v1/post"
	"microblog/internal/microblog/controller/v1/user"
	"microblog/internal/microblog/store"
	"microblog/internal/pkg/core"
	errno "microblog/internal/pkg/err"
	"microblog/internal/pkg/log"
	mw "microblog/internal/pkg/middleware"
	"microblog/pkg/auth"
)

func installRouters(g *gin.Engine) error {

	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	uc := user.New(store.S, authz)
	pc := post.New(store.S)

	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{

		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(mw.Authn(), mw.Authz(authz))
			userv1.GET(":name", uc.Get)
			userv1.PUT(":name", uc.Update)
			userv1.GET("", uc.List)
			userv1.DELETE(":name", uc.Delete)
		}

		postv1 := v1.Group("/posts", mw.Authn())
		{
			postv1.POST("", pc.Create)
			postv1.GET(":postID", pc.Get)
			postv1.PUT(":postID", pc.Update)
			postv1.DELETE("", pc.DeleteCollection)
			postv1.GET("", pc.List)
			postv1.DELETE(":postID", pc.Delete)
		}

	}

	return nil
}
