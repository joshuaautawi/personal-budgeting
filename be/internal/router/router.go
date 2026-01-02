package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"personal-budgeting/be/internal/handlers"
	"personal-budgeting/be/internal/services"
)

type App = fiber.App

type Deps struct {
	Category    *services.CategoryService
	Budget      *services.BudgetService
	Transaction *services.TxnService
	State       *services.StateService
}

func New(d Deps) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.Health())

	state := handlers.State{Svc: d.State}
	v1.Get("/state", state.Get)
	v1.Put("/state", state.Replace)

	cats := handlers.Categories{Svc: d.Category, BudgetSvc: d.Budget, TxnSvc: d.Transaction}
	v1.Get("/categories", cats.List)
	v1.Post("/categories", cats.Create)
	v1.Patch("/categories/:id", cats.Update)
	v1.Delete("/categories/:id", cats.Delete)

	budgets := handlers.Budgets{Svc: d.Budget, CatSvc: d.Category}
	v1.Get("/budgets", budgets.List)
	v1.Put("/budgets", budgets.Upsert)
	v1.Delete("/budgets/:id", budgets.Delete)

	txns := handlers.Transactions{Svc: d.Transaction, CatSvc: d.Category}
	v1.Get("/transactions", txns.List)
	v1.Post("/transactions", txns.Create)
	v1.Patch("/transactions/:id", txns.Update)
	v1.Delete("/transactions/:id", txns.Delete)

	return app
}


