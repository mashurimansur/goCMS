package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/adapter/http/handler"
)

// Options configure the HTTP router and its dependencies.
type Options struct {
	Mode          string
	PersonHandler *handler.PersonHandler
}

// NewGinEngine wires middleware stack and registers feature routes.
func NewGinEngine(opts Options) *gin.Engine {
	if opts.Mode != "" {
		gin.SetMode(opts.Mode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	admin := engine.Group("/api/v1/admin")
	if opts.PersonHandler != nil {
		opts.PersonHandler.Register(admin)
	}

	// Health probe for readiness checks.
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return engine
}
