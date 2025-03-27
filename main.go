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

	// Initial configuration
	app_configs := configs.NewConfig()

	// Initial database
	db_pool := database.GetDB()

	// Initialize router
	app := routers.New(db_pool)

	// Initialize redis client
	rdb := redis.NewRedisClient()

	// Initialize the translate
	if err := translate.Init(); err != nil {
		custom_log.NewCustomLog("Failed_initialize_i18n", err.Err.Error(), "error")
	}

	handler.NewFrontService(app, db_pool, rdb)

	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
	
}
