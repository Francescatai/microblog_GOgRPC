// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package v1

type CreatePostRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}

//  `POST /v1/posts`
type CreatePostResponse struct {
	PostID string `json:"postID"`
}

//  `GET /v1/posts/{postID}`
type GetPostResponse PostInfo

//  `PUT /v1/posts`
type UpdatePostRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}

type PostInfo struct {
	Username  string `json:"username,omitempty"`
	PostID    string `json:"postID,omitempty"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

//  `GET /v1/posts`
type ListPostRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

//  `GET /v1/posts`
type ListPostResponse struct {
	TotalCount int64       `json:"totalCount"`
	Posts      []*PostInfo `json:"posts"`
}
