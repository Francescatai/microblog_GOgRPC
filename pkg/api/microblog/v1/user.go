// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package v1

type LoginRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" valid:"required,stringlength(6|18)"`
	NewPassword string `json:"newPassword" valid:"required,stringlength(6|18)"`
}

type CreateUserRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" valid:"required,stringlength(1|255)"`
	Email    string `json:"email" valid:"required,email"`
	Phone    string `json:"phone" valid:"required,stringlength(10|10)"`
}

type GetUserResponse UserInfo

type UserInfo struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PostCount int64  `json:"postCount"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

//  `GET /v1/users`
type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

//  `GET /v1/users`
type ListUserResponse struct {
	TotalCount int64       `json:"totalCount"`
	Users      []*UserInfo `json:"users"`
}

// `PUT /v1/users/{name}`
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" valid:"stringlength(1|255)"`
	Email    *string `json:"email" valid:"email"`
	Phone    *string `json:"phone" valid:"stringlength(11|11)"`
}
