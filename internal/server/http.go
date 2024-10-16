package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "bubu_admin/api/v1"
	"bubu_admin/docs"
	"bubu_admin/internal/handler"
	"bubu_admin/internal/middleware"
	"bubu_admin/pkg/jwt"
	"bubu_admin/pkg/log"
	"bubu_admin/pkg/server/http"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	menuHandler *handler.MenuHandler,
	roleHandler *handler.RoleHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	api := s.Group("/api")

	v1 := api.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
		}
		/* // Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		} */

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			// strictAuthRouter.PUT("user/user", userHandler.UpdateProfile)

			// ========用户相关========
			// 为用户添加角色
			strictAuthRouter.POST("/user/addRoleToUser", userHandler.AddRoleToUser)
			// 获取用户信息
			strictAuthRouter.GET("/user/getProfile", userHandler.GetProfile)
			// 获取用户列表
			strictAuthRouter.POST("/user/listUser", userHandler.ListUser)
			// 获取用户信息
			strictAuthRouter.GET("/user/getUserByID/:id", userHandler.GetUserByID)
			// 删除用户
			strictAuthRouter.DELETE("/user/deleteUser/:id", userHandler.DeleteUser)
			// 更新用户
			strictAuthRouter.PUT("/user/updateUser/:id", userHandler.UpdateUser)

			// ========菜单相关========
			// 添加菜单
			strictAuthRouter.POST("/menu/createMenu", menuHandler.CreateMenu)

			// 获取菜单列表
			strictAuthRouter.POST("/menu/listMenu", menuHandler.ListMenu)
			// 更新菜单
			strictAuthRouter.PUT("/menu/updateMenu/:id", menuHandler.UpdateMenu)
			// 获取单个菜单
			strictAuthRouter.GET("/menu/getMenuByID/:id", menuHandler.GetMenuByID)
			// 删除菜单
			strictAuthRouter.DELETE("/menu/deleteMenu/:id", menuHandler.DeleteMenu)

			// ========角色相关========
			// 添加角色
			strictAuthRouter.POST("/role/createRole", roleHandler.CreateRole)
			// 修改角色
			strictAuthRouter.PUT("/role/updateRole/:id", roleHandler.UpdateRole)
			// 获取角色列表 
			strictAuthRouter.POST("/role/listRole", roleHandler.ListRole)
			// 删除角色
			strictAuthRouter.DELETE("/role/deleteRole/:id", roleHandler.DeleteRole)
			// 获取角色信息
			strictAuthRouter.GET("/role/getRoleByID/:id", roleHandler.GetRoleByID)
		}
	}

	return s
}
