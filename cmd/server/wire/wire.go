//go:build wireinject
// +build wireinject

package wire

import (
	"bubu_admin/internal/handler"
	"bubu_admin/internal/repository"
	"bubu_admin/internal/server"
	"bubu_admin/internal/service"
	"bubu_admin/pkg/app"
	"bubu_admin/pkg/jwt"
	"bubu_admin/pkg/casbin"
	"bubu_admin/pkg/log"
	"bubu_admin/pkg/server/http"
	"bubu_admin/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewMenuRepository,
	repository.NewRoleRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewMenuService,
	service.NewRoleService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewMenuHandler,
	handler.NewRoleHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		casbin.NewCasbin,
		newApp,
	))
}
