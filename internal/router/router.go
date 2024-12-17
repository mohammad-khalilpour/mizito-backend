package router

import (
	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
)



type Router struct {
	app *fiber.App
	cache bigcache.BigCache
}



func (r *Router) Init() {
	InitAuth(r)
	InitProject(r)
	InitSubtask(r)
	InitTask(r)
	InitUser(r)
}