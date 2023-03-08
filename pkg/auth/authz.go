// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	aclModel = `[request_definition]
				r = sub, obj, act
				[policy_definition]
				p = sub, obj, act
				[policy_effect]
				e = some(where (p.eft == allow))
				[matchers]
				m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)`
)

type Authz struct {
	*casbin.SyncedEnforcer
}

func NewAuthz(db *gorm.DB) (*Authz, error) {
	// Initialize a Gorm adapter and use it in a Casbin enforcer
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, _ := model.NewModelFromString(aclModel)

	// Initialize the enforcer.
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// Load the policy from DB.
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{enforcer}

	return a, nil
}

func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
