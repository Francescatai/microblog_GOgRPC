// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"microblog/internal/pkg/known"
)

// 用来在每個 HTTP 請求的 context, response 中加入 `X-Request-ID` key-value
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查request header中是否有 `X-Request-ID`，沒有則新增
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 將 RequestID 保存在 gin.Context 中
		c.Set(known.XRequestIDKey, requestID)

		// 將 RequestID 保存在 HTTP response header，Header key: `X-Request-ID`
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}

