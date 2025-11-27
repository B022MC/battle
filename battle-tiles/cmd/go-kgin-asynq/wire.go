//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"battle-tiles/internal/biz"
	"battle-tiles/internal/conf"
	"battle-tiles/internal/dal/repo"
	"battle-tiles/internal/infra"
	"battle-tiles/internal/server"
	"battle-tiles/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Global, *conf.Server, *conf.Data, *conf.Log, log.Logger) (*kratos.App, func(), error) {

	panic(wire.Build(server.ProviderSet, infra.ProviderSet, repo.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
