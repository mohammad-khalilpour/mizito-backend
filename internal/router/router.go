package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	return &Router{
		App: app,
		Cfg: cfg,
	}
}

func (r *Router) Init() {
	InitAuth(r)
	InitProject(r)
	InitSubtask(r)
	InitTask(r)
	InitUser(r)
	InitSocket(r)
}

func (r *Router) Run() {
	if err := r.App.Listen(r.Cfg.AppPort); err != nil {
		panic(err)
	}
}
