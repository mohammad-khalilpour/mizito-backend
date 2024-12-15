package router

import "github.com/gofiber/fiber/v2"



type Router struct {
	app *fiber.App
}



func (r *Router) Init() {
	InitAuth(r)
	InitProject(r)
	InitSubtask(r)
	InitTask(r)
	InitUser(r)
}