// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errno "microblog/internal/pkg/err"
)

type ErrResponse struct {
	Code string `json:"code"`

	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		httpCode, code, message := errno.Decode(err)
		c.JSON(httpCode, ErrResponse{
			Code:    code,
			Message: message,
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
