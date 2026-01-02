package app

import (
	"context"
	"log"

	"personal-budgeting/be/internal/clock"
	"personal-budgeting/be/internal/db"
	"personal-budgeting/be/internal/dbmodel"
	"personal-budgeting/be/internal/id"
	"personal-budgeting/be/internal/repositories"
	"personal-budgeting/be/internal/router"
	"personal-budgeting/be/internal/services"
)

// New wires the full application (router → handlers → services → repositories).
func New() *router.App {
	clk := clock.Real{}
	ids := id.RandomHex{}

	pgCfg := db.LoadPostgresConfigFromEnv()
	if pgCfg.Enabled() {
		gdb, sqlDB, err := db.OpenPostgresGorm(context.Background(), pgCfg)
		if err != nil {
			log.Fatalf("postgres connect error: %v", err)
		} else {
			// optional: close DB on process exit (Fiber doesn't manage it)
			_ = sqlDB

			// Ensure tables exist if you don't run migrations yet.
			if err := dbmodel.AutoMigrate(gdb); err != nil {
				log.Fatalf("postgres automigrate error: %v", err)
			} else {
				catRepo := repositories.NewGormCategoryRepo(gdb)
				budgetRepo := repositories.NewGormBudgetRepo(gdb)
				txnRepo := repositories.NewGormTxnRepo(gdb)

				categorySvc := services.NewCategoryService(clk, ids, catRepo)
				budgetSvc := services.NewBudgetService(clk, ids, budgetRepo)
				txnSvc := services.NewTxnService(clk, ids, txnRepo)
				stateSvc := services.NewStateService(catRepo, budgetRepo, txnRepo)

				return router.New(router.Deps{
					Category:    categorySvc,
					Budget:      budgetSvc,
					Transaction: txnSvc,
					State:       stateSvc,
				})
			}
		}
	}

	log.Fatal("DATABASE_URL (or PG* env vars) is required; memory repo removed")
	return nil
}
