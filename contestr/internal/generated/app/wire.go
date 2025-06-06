// wire.go
//go:build wireinject
// +build wireinject

package app

import (
	"contestr/cmd/app/wiresets"
	"contestr/internal/configs"
	"contestr/internal/transport"
	"context"
	"github.com/google/wire"
)

type App struct {
	Ctx        context.Context
	Cfg        *configs.Config
	HttpServer *transport.HTTPServer
	TgBot      *transport.TgBot
}

func InitializeService() (*App, error) {
	wire.Build(
		wiresets.CommonSet,
		wiresets.HTTPSet,
		wiresets.TgBotSet,

		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
