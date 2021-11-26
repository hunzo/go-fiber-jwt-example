package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

const JWT_SECRET = "thisissecret"

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", Public)
	app.Get("/token/:payload", func(c *fiber.Ctx) error {
		p := c.Params("payload")

		token, err := GenerateToken(p)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"success": true,
			"payload": p,
			"token":   token,
		})
	})

	private := app.Group("/private", AuthedRequired())
	private.Get("/", Private)

	app.Listen(":8080")
}

func GenerateToken(payload string) (string, error) {
	claims := jwt.MapClaims{
		"payload":    payload,
		"role":       "admin",
		"token_type": "access_token",
		"exp":        time.Now().Add(time.Minute * 10).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return token, nil
}

func AuthedRequired() func(c *fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey:   []byte(JWT_SECRET),
		ErrorHandler: jwtError,
	}
	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}

func Public(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"info": "Public page",
	})
}

func Private(c *fiber.Ctx) error {
	// Default context c.Locals("user") jwtware
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.JSON(fiber.Map{
		"info":       "Private",
		"role":       claims["role"].(string),
		"token_type": claims["token_type"].(string),
		"payload":    claims["payload"].(string),
	})
}
