package transport

import (
	"contestr/internal/auth"
	"contestr/internal/configs"
	"contestr/internal/generated/server"
	"contestr/internal/handlers"
	authmiddleware "contestr/internal/middleware"
	"contestr/pkg/logger"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type HTTPServer struct {
	echo       *echo.Echo
	cfg        *configs.Config
	httpServer *http.Server
}

func NewHTTPServer(ctx context.Context, handlers *handlers.Handlers, cfg *configs.Config, authService *auth.Service) *HTTPServer {
	logger.Infof(ctx, "server configuration: address=%s, read_timeout=%v, write_timeout=%v",
		cfg.HTTP.Port, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout)

	e := newEcho()
	httpServer := &http.Server{
		Handler:      e,
		Addr:         ":" + cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	httpServerWrapper := &HTTPServer{
		echo:       e,
		cfg:        cfg,
		httpServer: httpServer,
	}

	httpServerWrapper.RegisterHandlers(handlers, authService)

	return httpServerWrapper
}

func newEcho() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.Logger())
	return e
}

func (s *HTTPServer) RegisterHandlers(handlers *handlers.Handlers, authService *auth.Service) {
	s.echo.Use(authmiddleware.AdminJWT(authService))
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HTTPServer is running")
	})

	server.RegisterHandlers(s.echo, handlers)
	registerPublicRegattaRoutes(s.echo, handlers)

	routes := s.echo.Routes()
	for _, route := range routes {
		logger.Infof(context.Background(), "route registered: %s %s", route.Method, route.Path)
	}
	logger.Info(context.Background(), "handlers registered")
}

func (s *HTTPServer) Start(_ context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	logger.Info(ctx, "shutting down http server")
	return s.echo.Shutdown(ctx)
}
