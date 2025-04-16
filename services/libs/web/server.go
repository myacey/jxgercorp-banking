package web

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Listen            string        `yaml:"listen"`
	DrainInterval     time.Duration `yaml:"drainInterval"`
	Profile           bool          `yaml:"profile"`
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	Env               string        `yaml:"env"`
}

// Server if an interface for web http server.
type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Router() gin.IRouter
	Ready() bool
}

var _ Server = (*BaseServer)(nil)

// BaseServer is a default implementation of Server interface.
type BaseServer struct {
	engine     *gin.Engine
	httpServer *http.Server

	isNotReady int32
}

// NewServer returns new *BaseServer.
func NewServer(cfg ServerConfig, handler *gin.Engine) *BaseServer {
	s := &BaseServer{
		engine: handler,
	}

	// server tweaks
	s.httpServer = &http.Server{
		Addr:              cfg.Listen,
		Handler:           s.engine,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return s
}

func (s *BaseServer) Router() gin.IRouter {
	return s.engine
}

func (s *BaseServer) Ready() bool {
	return atomic.LoadInt32(&s.isNotReady) == 0
}

func (s *BaseServer) Run(ctx context.Context) error {
	go func() {
		for {
			<-ctx.Done()
			return
		}
	}()

	log.Printf("start server on %v", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *BaseServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
