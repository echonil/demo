package core

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	sqlxadapter "github.com/memwey/casbin-sqlx-adapter"
)

const enforceModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

func NewEnforcer(cfg DBConfig) (*casbin.Enforcer, error) {
	opts := &sqlxadapter.AdapterOptions{
		DriverName:     cfg.Driver,
		DataSourceName: cfg.DSN,
		TableName:      "auth_rules",
	}
	m, err := model.NewModelFromString(enforceModel)
	if err != nil {
		return nil, err
	}

	a := sqlxadapter.NewAdapterFromOptions(opts)
	e, err := casbin.NewEnforcer()
	if err != nil {
		return nil, err
	}
	if err = e.InitWithModelAndAdapter(m, a); err != nil {
		return nil, err
	}

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	return e, nil
}
