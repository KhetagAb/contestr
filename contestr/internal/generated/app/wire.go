// wire.go
//go:build wireinject
// +build wireinject

package app

import (
	"contestr/cmd/app/wiresets"
	"contestr/internal/configs"
	"contestr/internal/services/contest_sync"
	"contestr/internal/transport"
	"context"

	"github.com/google/wire"
)

type App struct {
	Ctx         context.Context
	Cfg         *configs.Config
	HttpServer  *transport.HTTPServer
	TgBot       *transport.TgBot
	ContestSync *contest_sync.ContestSyncService
}

func InitializeService() (*App, error) {
	wire.Build(
		wiresets.All,

		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
