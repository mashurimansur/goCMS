package app

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/adapter/http/handler"
	"github.com/mashurimansur/goCMS/internal/adapter/http/router"
	domainperson "github.com/mashurimansur/goCMS/internal/domain/person"
	sqlperson "github.com/mashurimansur/goCMS/internal/repository/person"
	personusecase "github.com/mashurimansur/goCMS/internal/usecase/person"
	"github.com/mashurimansur/goCMS/internal/utils/config"
	"github.com/mashurimansur/goCMS/internal/utils/database"
)

// Application wires all layers (infrastructure, use cases, delivery) so main can remain minimal.
type Application struct {
	engine   *gin.Engine
	httpAddr string
	dbConn   *database.Connection
}

// New creates a fully wired application instance ready to run.
func New(ctx context.Context, cfg config.AppConfig) (*Application, error) {
	dbConn, err := database.NewConnection(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	personRepo, err := buildPersonRepository(dbConn)
	if err != nil {
		return nil, err
	}

	personUseCase := personusecase.New(personRepo)
	personHandler := handler.NewPersonHandler(personUseCase)

	engine := router.NewGinEngine(router.Options{
		Mode:          cfg.GinMode,
		PersonHandler: personHandler,
	})

	app := &Application{
		engine:   engine,
		httpAddr: cfg.HTTPAddr,
		dbConn:   dbConn,
	}

	return app, nil
}

// Run starts the HTTP server using the configured engine and address.
func (a *Application) Run() error {
	if a == nil || a.engine == nil {
		return errors.New("application engine is not configured")
	}

	addr := a.httpAddr
	if addr == "" {
		addr = ":8080"
	}

	return a.engine.Run(addr)
}

// Close releases infrastructure resources such as database connections.
func (a *Application) Close() error {
	if a == nil || a.dbConn == nil {
		return nil
	}

	return a.dbConn.Close()
}

func buildPersonRepository(dbConn *database.Connection) (domainperson.Repository, error) {
	if dbConn == nil || dbConn.DB == nil {
		return nil, errors.New("database connection is required for person repository")
	}

	return sqlperson.New(dbConn.DB)
}
