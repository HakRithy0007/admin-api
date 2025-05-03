package main

import (
	config "admin-phone-shop-api/config"
	database "admin-phone-shop-api/config/database"
	redis "admin-phone-shop-api/config/redis"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	translate "admin-phone-shop-api/pkg/utils/translate"
	routers "admin-phone-shop-api/routers"
	handler "admin-phone-shop-api/handler"
	"fmt"
)

func main() {

	// Configuration
	app_configs := config.NewConfig()

	// Database
	db_pool := database.GetDB()

	// Routers
	app := routers.New(db_pool)

	// Redis
	rdb  := redis.NewRedisClient()

	// Translate
	if err := translate.Init(); err != nil {
		custom_log.NewCustomLog("Failed_initialize_i18n", err.Err.Error(), "error")
	}

	handler.NewFrontService(app, db_pool, rdb)

	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))

}
