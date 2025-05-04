package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	auth "admin-phone-shop-api/internal/auth"
	custom_models "admin-phone-shop-api/pkg/model"
	"admin-phone-shop-api/pkg/utils/response"
	custom_translate "admin-phone-shop-api/pkg/utils/translate"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewJwtMinddleWare(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) {
	errs := godotenv.Load()
	if errs != nil {
		log.Fatalf("Error loading .env file")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")

	// JWT Middleware for HTTP requests
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret_key)},
		ContextKey: "jwt_data",
	}))

	// Middleware to add admin context
	app.Use(func(c *fiber.Ctx) error {
		admin_token := c.Locals("jwt_data").(*jwt.Token)
		uclaim := admin_token.Claims.(jwt.MapClaims)
		return handleAdminContext(c, uclaim, db_pool, redis)
	})
}

func handleAdminContext(c *fiber.Ctx, uclaim jwt.MapClaims, db *sqlx.DB, redis *redis.Client) error {
	login_session, ok := uclaim["login_session"].(string)
	if !ok || login_session == "" {
		smg_error := response.NewResponseError(
			custom_translate.Translate(c, "login_session_missing"),
			-500,
			fmt.Errorf("%s", custom_translate.Translate(c, "login_session_missing")),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	uCtx := custom_models.AdminContext{
		AdminID:     uclaim["admin_id"].(int),
		Admin_Name:     uclaim["admin_name"].(string),
		LoginSession: login_session,
		Exp:          time.Unix(int64(uclaim["exp"].(float64)), 0),
		AdminAgent:    c.Get("Admin-Agent", "unknown"),
		Ip:           c.IP(),
	}
	c.Locals("AdminContext", uCtx)

	sv := auth.NewAuthService(db, redis)
	success, err := sv.CheckSession(login_session, uCtx.AdminID)
	if err != nil || !success {
		smg_error := response.NewResponseError(
			custom_translate.Translate(c, "login_session_invalid"),
			-500,
			fmt.Errorf("%s", custom_translate.Translate(c, "login_session_invalid")),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	return c.Next()
}
