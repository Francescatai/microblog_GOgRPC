// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"microblog/internal/microblog/store"
	errno "microblog/internal/pkg/err"
	"microblog/internal/pkg/log"
	"microblog/internal/pkg/model"
	v1 "microblog/pkg/api/microblog/v1"
	"microblog/pkg/auth"
	"microblog/pkg/token"
)

type UserBiz interface {
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
	Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserModel
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)

	resp.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")

	return &resp, nil
}

func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user := item
		users = append(users, &v1.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}

func (b *userBiz) Update(ctx context.Context, username string, user *v1.UpdateUserRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if user.Email != nil {
		userM.Email = *user.Email
	}

	if user.Nickname != nil {
		userM.Nickname = *user.Nickname
	}

	if user.Phone != nil {
		userM.Phone = *user.Phone
	}

	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

func (b *userBiz) Delete(ctx context.Context, username string) error {
	if err := b.ds.Users().Delete(ctx, username); err != nil {
		return err
	}

	return nil
}
