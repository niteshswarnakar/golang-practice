package echoserver

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	casbin_middleware "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Runner() {
	e := echo.New()

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", "localhost", "test_user", "test_password", "5432", "test_db")
	dbAdapter, err := gormadapter.NewAdapter("postgres", dsn, true)
	if err != nil {
		e.Logger.Fatal("Failed to create Gorm adapter:", err)
	}

	enforcer, err := casbin.NewEnforcer("echoserver/model.conf", dbAdapter)
	if err != nil {
		e.Logger.Fatal(err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		e.Logger.Fatal("Failed to load policy:", err)
	}

	// HOW TO GET LIST OF ROLES
	roles, err := enforcer.GetAllRoles()
	if err != nil {
		e.Logger.Fatal("Failed to get roles:", err)
	}
	fmt.Println("All roles:", roles)

	// GET ALL THE USERS (Method 1: Get all subjects from grouping policies)
	allUsers, err := enforcer.GetAllSubjects()
	if err != nil {
		fmt.Printf("Error getting subjects: %v\n", err)
	} else {
		fmt.Println("All users (subjects):", allUsers)
	}

	// GET ALL USERS (Method 2: Extract users from grouping policies manually)
	groupings, _ := enforcer.GetGroupingPolicy()
	users := []string{}
	for _, grouping := range groupings {
		if len(grouping) >= 1 {
			users = append(users, grouping[0]) // First element is the user
		}
	}

	fmt.Println("\nExtracted users from grouping policies:", users)

	// Create 'viewer' role with read permissions
	enforcer.AddPolicy("viewer", "/api/v1/kathmandu", "GET")
	enforcer.AddPolicy("viewer", "/api/v1/pokhara", "GET")
	fmt.Println("Created 'viewer' role with read permissions")

	// Create 'editor' role with read/write permissions
	enforcer.AddPolicy("editor", "/api/v1/kathmandu", "GET")
	enforcer.AddPolicy("editor", "/api/v1/kathmandu", "POST")
	enforcer.AddPolicy("editor", "/api/v1/pokhara", "GET")
	enforcer.AddPolicy("editor", "/api/v1/pokhara", "POST")
	fmt.Println("Created 'editor' role with read/write permissions")

	// Create 'manager' role with all permissions
	enforcer.AddPolicy("manager", "/api/v1/kathmandu", "*")
	enforcer.AddPolicy("manager", "/api/v1/pokhara", "*")
	enforcer.AddPolicy("manager", "/api/v1/public", "*")
	fmt.Println("Created 'manager' role with all permissions")

	// STEP 2: Verify roles were created
	fmt.Println("\n=== STEP 2: Verifying Roles ===")
	allRoles, _ := enforcer.GetAllRoles()
	fmt.Println("All available roles:", allRoles)

	// STEP 3: Create users by assigning them to roles (users are created implicitly)
	fmt.Println("\n=== STEP 3: Creating Users and Assigning Roles ===")

	// Create user 'john' and assign to 'viewer' role
	enforcer.AddGroupingPolicy("john", "viewer")
	fmt.Println("Created user 'john' and assigned to 'viewer' role")

	// Create user 'sarah' and assign to 'editor' role
	enforcer.AddGroupingPolicy("sarah", "editor")
	fmt.Println("Created user 'sarah' and assigned to 'editor' role")

	// Create user 'mike' and assign to 'manager' role
	enforcer.AddGroupingPolicy("mike", "manager")
	fmt.Println("Created user 'mike' and assigned to 'manager' role")

	// STEP 4: Later, reassign users to different roles
	fmt.Println("\n=== STEP 4: Reassigning User Roles ===")

	// Promote john from viewer to editor
	enforcer.RemoveGroupingPolicy("john", "viewer")
	enforcer.AddGroupingPolicy("john", "editor")
	fmt.Println("Promoted 'john' from 'viewer' to 'editor'")

	// Give sarah multiple roles (editor + manager)
	enforcer.AddGroupingPolicy("sarah", "manager")
	fmt.Println("Added 'manager' role to 'sarah' (now has both editor and manager)")

	// STEP 5: Verify final user-role assignments
	fmt.Println("\n=== STEP 5: Final User-Role Assignments ===")
	finalUsers, _ := enforcer.GetAllSubjects()
	fmt.Println("All users:", finalUsers)

	for _, user := range finalUsers {
		userRoles, _ := enforcer.GetRolesForUser(user)
		fmt.Printf("User '%s' has roles: %v\n", user, userRoles)
	}

	// STEP 6: Test permissions
	fmt.Println("\n=== STEP 6: Testing Permissions ===")

	// Test john's access (should have editor permissions)
	johnCanRead, _ := enforcer.Enforce("john", "/api/v1/kathmandu", "GET")
	johnCanWrite, _ := enforcer.Enforce("john", "/api/v1/kathmandu", "POST")
	fmt.Printf("John can READ: %v, can WRITE: %v\n", johnCanRead, johnCanWrite)

	// Test mike's access (should have manager permissions)
	mikeCanDelete, _ := enforcer.Enforce("mike", "/api/v1/kathmandu", "DELETE")
	fmt.Printf("Mike can DELETE: %v\n", mikeCanDelete)

	// Save all changes to database
	enforcer.SavePolicy()
	fmt.Println("\n=== All changes saved to database ===")

	fmt.Println("\n=== Starting Echo Server ===")

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
