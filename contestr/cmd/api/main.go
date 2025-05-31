package api

import (
	"contestr/internal/configs"
	httphandlers "contestr/internal/handlers/http"
	"contestr/internal/services/codeforces"
	"contestr/internal/transport"
	"context"
)

func CreateHttpServer(ctx context.Context, cfg *configs.Config) *transport.Server {
	cfService := codeforces.NewService(cfg)
	handlers := httphandlers.NewHandlers(cfService)

	httpServer := transport.NewHTTPServer(ctx, cfg)
	httpServer.RegisterHandlers(handlers)
	return httpServer
}
