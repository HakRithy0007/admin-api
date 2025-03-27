package main

import (
	configs "admin-phone-shop-api/config"
	database "admin-phone-shop-api/config/database"
	redis "admin-phone-shop-api/config/redis"
	"admin-phone-shop-api/handler"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	translate "admin-phone-shop-api/pkg/utils/translate"
	routers "admin-phone-shop-api/routers"
	"fmt"
)

func main() {

	// Initial configurations
	app_configs := configs.NewConfig()

	// Initial databases
	db_pool := database.GetDB()

	// Initialize routers
	app := routers.New(db_pool)

	// Initialize redis clients
	rdb := redis.NewRedisClient()

	// Initialize the translates
	if err := translate.Init(); err != nil {
		custom_log.NewCustomLog("Failed_initialize_i18n", err.Err.Error(), "error")
	}

	handler.NewFrontService(app, db_pool, rdb)

	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))

}
