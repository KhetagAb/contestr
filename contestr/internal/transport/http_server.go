package transport

import (
	"contestr/internal/configs"
	"contestr/internal/generated/server"
	httphandlers "contestr/internal/handlers/http"
	"contestr/pkg/logger"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	echo       *echo.Echo
	cfg        *configs.Config
	httpServer *http.Server
}

func NewHTTPServer(ctx context.Context, cfg *configs.Config) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	httpServer := &http.Server{
		Handler:      e,
		Addr:         ":" + cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	logger.Infof(ctx, "server configuration: address=%s, read_timeout=%v, write_timeout=%v",
		cfg.HTTP.Port, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout)

	return &Server{
		echo:       e,
		cfg:        cfg,
		httpServer: httpServer,
	}
}

func (s *Server) RegisterHandlers(handlers *httphandlers.Handlers) {
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running")
	})

	server.RegisterHandlers(s.echo, handlers)

	routes := s.echo.Routes()
	for _, route := range routes {
		logger.Infof(context.Background(), "route registered: %s %s", route.Method, route.Path)
	}
	logger.Info(context.Background(), "handlers registered")
}

func (s *Server) Start(_ context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info(ctx, "shutting down http server")
	return s.echo.Shutdown(ctx)
}
