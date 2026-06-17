//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/zhiyusec/kratos-layout/internal/biz"
	"github.com/zhiyusec/kratos-layout/internal/conf"
	"github.com/zhiyusec/kratos-layout/internal/data"
	"github.com/zhiyusec/kratos-layout/internal/server"
	"github.com/zhiyusec/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
