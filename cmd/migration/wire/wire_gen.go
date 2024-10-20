// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"bubu_admin/internal/repository"
	"bubu_admin/internal/server"
	"bubu_admin/pkg/app"
	"bubu_admin/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	db := repository.NewDB(viperViper, logger)
	migrate := server.NewMigrate(db, logger)
	appApp := newApp(migrate)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewUserRepository, repository.NewMenuRepository)

var serverSet = wire.NewSet(server.NewMigrate)

// build App
func newApp(
	migrate *server.Migrate,
) *app.App {
	return app.NewApp(app.WithServer(migrate), app.WithName("demo-migrate"))
}
