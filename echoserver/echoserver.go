package echoserver

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	casbin_middleware "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Runner() {
	e := echo.New()
	enforcer, err := casbin.NewEnforcer("echoserver/model.conf", "echoserver/policy.csv")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Add middleware to extract user from Basic Auth
	e.Use(extractUserMiddleware)

	// Use Casbin middleware with custom config
	e.Use(casbin_middleware.MiddlewareWithConfig(casbin_middleware.Config{
		Enforcer: enforcer,
		UserGetter: func(c echo.Context) (string, error) {
			// Get user from context (set by extractUserMiddleware)
			if user := c.Get("user"); user != nil {
				return user.(string), nil
			}
			return "anonymous", nil // Default user
		},
	}))

	e.GET("/api/v1/kathmandu", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello, this is a protected route! /api/v1/kathmandu : Read access granted",
		})
	})

	e.POST("/api/v1/kathmandu", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello, this is a protected route! /api/v1/kathmandu : Write access granted",
		})
	})
	e.GET("/api/v1/pokhara", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello, this is a protected route! /api/v1/pokhara : Read access granted",
		})
	})

	e.POST("/api/v1/pokhara", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello, this is a protected route! /api/v1/pokhara : Write access granted",
		})
	})

	e.GET("/api/v1/public", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello, this is a public route! /api/v1/public : No access control",
		})
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8000"))
}

// Middleware to extract user from Basic Auth
func extractUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		fmt.Println("BASIC AUTH : ", auth)
		if auth == "" {
			c.Set("user", "anonymous")
			return next(c)
		}

		if strings.HasPrefix(auth, "Basic ") {
			userCreds := auth[6:]

			pair := strings.SplitN(userCreds, ":", 2)
			if len(pair) == 2 {
				fmt.Printf("Username: %s, Password: %s\n", pair[0], pair[1])
				c.Set("user", pair[0]) // Set username
			} else {
				fmt.Println("Invalid credentials format")
				c.Set("user", "anonymous")
			}
		} else {
			c.Set("user", "anonymous")
		}

		return next(c)
	}
}
