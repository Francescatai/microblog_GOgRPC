// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package err

var (
	
	OK = &Errno{HTTP: 200, Code: "", Message: ""}
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}
	ErrPageNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found."}

	ErrBind = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: "Error occurred while binding the request body to the struct."}
	ErrInvalidParameter = &Errno{HTTP: 400, Code: "InvalidParameter", Message: "Parameter verification failed."}

	ErrSignToken = &Errno{HTTP: 401, Code: "AuthFailure.SignTokenError", Message: "Error occurred while signing the JSON web token."}
	ErrTokenInvalid = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "Token was invalid."}

	ErrUnauthorized = &Errno{HTTP: 401, Code: "AuthFailure.Unauthorized", Message: "Unauthorized."}
)