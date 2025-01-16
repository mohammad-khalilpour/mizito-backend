package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"mizito/internal/database"
	"mizito/internal/env"
	"mizito/internal/middleware"
)

type Router struct {
	App *fiber.App
	Cfg *env.Config
}

func InitApp(cfg *env.Config) *Router {
	app := fiber.New()

	app.Use(recover.New())

	app.Use("/ws/:id", middleware.UpgradeMiddleware)
	app.Use(middleware.AuthMiddleware)

	return &Router{
		App: app,
		Cfg: cfg,
	}
}

func (r *Router) Init(env *env.Config) {

	redis := database.NewRedisHandler(env)
	mongo := database.NewMongoHandler(env)
	postgreSql := database.NewDatabaseHandler(env)

	InitAuth(r, env.AuthorizationSecret, redis, postgreSql)
	InitProject(r, postgreSql)
	InitSubtask(r, postgreSql)
	InitTask(r, postgreSql)
	InitUser(r)
	InitSocket(r, redis, mongo, postgreSql, env)
}

func (r *Router) Run() {
	if err := r.App.Listen(r.Cfg.AppPort); err != nil {
		panic(err)
	}
}
