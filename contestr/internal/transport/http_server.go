package transport

import (
	"contestr/internal/config"
	"contestr/internal/generated/server"
	httphandlers "contestr/internal/handlers/http"
	"contestr/pkg/logger"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	echo       *echo.Echo
	cfg        *config.Config
	httpServer *http.Server
}

func NewHTTPServer(ctx context.Context, cfg *config.Config) *Server {
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

func (s *Server) RegisterHandlers() {
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running")
	})

	handlers := httphandlers.NewHandlers()
	server.RegisterHandlers(s.echo, handlers)

	routes := s.echo.Routes()
	for _, route := range routes {
		logger.Infof(context.Background(), "route registered: %s %s", route.Method, route.Path)
	}
	logger.Info(context.Background(), "handlers registered")
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		logger.Infof(ctx, "http server started")
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf(ctx, "failed to start http server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info(ctx, "shutting down http server")
	return s.echo.Shutdown(ctx)
}
