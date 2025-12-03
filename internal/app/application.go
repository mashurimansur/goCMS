package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/adapter/http/handler"
	"github.com/mashurimansur/goCMS/internal/adapter/http/router"
	domainperson "github.com/mashurimansur/goCMS/internal/domain/person"
	sqlperson "github.com/mashurimansur/goCMS/internal/repository/person"
	sqluser "github.com/mashurimansur/goCMS/internal/repository/user"
	personusecase "github.com/mashurimansur/goCMS/internal/usecase/person"
	userusecase "github.com/mashurimansur/goCMS/internal/usecase/user"
	"github.com/mashurimansur/goCMS/internal/utils/config"
	"github.com/mashurimansur/goCMS/internal/utils/database"
	"github.com/mashurimansur/goCMS/internal/utils/token"
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

	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	tokenDuration, err := time.ParseDuration(cfg.TokenDuration)
	if err != nil {
		return nil, fmt.Errorf("cannot parse token duration: %w", err)
	}

	userRepo := sqluser.NewUserRepository(dbConn.DB)
	userUseCase := userusecase.NewUserUseCase(userRepo, tokenMaker, tokenDuration)
	userHandler := handler.NewUserHandler(userUseCase)

	engine := router.NewGinEngine(router.Options{
		Mode:          cfg.GinMode,
		PersonHandler: personHandler,
		UserHandler:   userHandler,
		TokenMaker:    tokenMaker,
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
