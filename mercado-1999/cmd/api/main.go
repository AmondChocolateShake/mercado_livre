package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/config"
	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/database"
	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/item"
	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/order"
)

func main() {
	decimal.MarshalJSONWithoutQuotes = true

	cfg := config.Load()

	if err := database.CreateIfNotExists(cfg.ServerDSN(), cfg.DBName); err != nil {
		log.Fatalf("create database: %v", err)
	}

	db, err := database.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db, "migrations"); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	itemRepo := item.NewRepository(db)
	itemSvc := item.NewService(itemRepo)

	orderRepo := order.NewRepository(db)
	orderSvc := order.NewService(db, orderRepo, itemRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	item.NewHandler(itemSvc).RegisterRoutes(r)
	order.NewHandler(orderSvc).RegisterRoutes(r)

	addr := ":" + cfg.Port
	log.Println("server listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
