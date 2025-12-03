package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/adapter/http/handler"
	"github.com/mashurimansur/goCMS/internal/adapter/http/middleware"
	"github.com/mashurimansur/goCMS/internal/utils/token"
)

// Options configure the HTTP router and its dependencies.
type Options struct {
	Mode          string
	PersonHandler *handler.PersonHandler
	UserHandler   *handler.UserHandler
	TokenMaker    token.Maker
}

// NewGinEngine wires middleware stack and registers feature routes.
func NewGinEngine(opts Options) *gin.Engine {
	if opts.Mode != "" {
		gin.SetMode(opts.Mode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	// Public routes
	if opts.UserHandler != nil {
		// Register public auth routes and protected user routes
		authMiddleware := middleware.AuthMiddleware(opts.TokenMaker)
		opts.UserHandler.Register(engine.Group("/api/v1"), authMiddleware)
	}

	admin := engine.Group("/api/v1/admin")
	if opts.TokenMaker != nil {
		authMiddleware := middleware.AuthMiddleware(opts.TokenMaker)
		admin.Use(authMiddleware)
	}
	if opts.PersonHandler != nil {
		opts.PersonHandler.Register(admin)
	}

	// Health probe for readiness checks.
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return engine
}
