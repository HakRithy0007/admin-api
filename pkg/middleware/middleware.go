package middlewares

import (
	auth "admin-phone-shop-api/internal/auth"
	model "admin-phone-shop-api/pkg/model"
	response "admin-phone-shop-api/pkg/utils/response"
	translate "admin-phone-shop-api/pkg/utils/translate"

	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	app.Use(func(c *fiber.Ctx) error {
		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			webSocketProtocol := c.Get("Sec-webSocket-Protocol")
			if webSocketProtocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing WebSocket protocol for authentication",
				})
			}

			parts := strings.Split(webSocketProtocol, ",")
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid webSocket protocal authenticaion format",
				})
			}

			tokenString := strings.TrimSpace(parts[1])

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret_key), nil
			})
			if err != nil || !token.Valid {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired JWT token",
				})
			}

			c.Locals("jwt_data", token)
			c.Set("Sec-WebSocket-Protocol", "Bearer")
			return c.Next()
		}

		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret_key)},
			ContextKey: "jwt_data",
		})(c)
	})

	app.Use(func(c *fiber.Ctx) error {
		admin_token := c.Locals("jwt_data").(*jwt.Token)
		uclaim := admin_token.Claims.(jwt.MapClaims)

		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			return handleAdminContext(c, uclaim, db_pool, redis)
		}
		return handleAdminContext(c, uclaim, db_pool, redis)
	})
}

func handleAdminContext(c *fiber.Ctx, uclaim jwt.MapClaims, db *sqlx.DB, redis *redis.Client) error {

	login_session, ok := uclaim["login_session"].(string)
	if !ok || login_session == "" {
		smg_error := response.NewResponseError(
			translate.Translate("loginSessionMissing", nil, c),
			-500,
			fmt.Errorf("%s", translate.Translate("loginSessionMissing", nil, c)),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	uCtx := model.AdminContext{
		AdminID:      uclaim["admin_id"].(float64),
		Admin_Name:   uclaim["admin_name"].(string),
		RoleID:       int(uclaim["role_id"].(float64)),
		LoginSession: uclaim["login_session"].(string),
		Exp:          time.Unix(int64(uclaim["exp"].(float64)), 0),
		AdminAgent:   c.Get("Admin-Agent", "unknown"),
		Ip:           c.IP(),
	}
	c.Locals("AdminContext", uCtx)

	sv := auth.NewAuthService(db, redis)
	success, err := sv.CheckSession(login_session, uCtx.AdminID)
	if err != nil || !success {
		smg_error := response.NewResponseError(
			translate.Translate("loginSessionInvalid", nil, c),
			-500,
			fmt.Errorf("%s", translate.Translate("loginSessionInvalid", nil, c)),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	return c.Next()
}
