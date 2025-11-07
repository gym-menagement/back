package router

import (
  "strconv"
  "strings"

	"gym/router/routers"

	"github.com/gofiber/fiber/v2"
)

func getArrayCommal(name string) []int64 {
	values := strings.Split(name, ",")

	var items []int64
	for _, item := range values {
        n, _ := strconv.ParseInt(item, 10, 64)
		items = append(items, n)
	}

	return items
}

func getArrayCommai(name string) []int {
	values := strings.Split(name, ",")

	var items []int
	for _, item := range values {
        n, _ := strconv.Atoi(item)
		items = append(items, n)
	}

	return items
}

func SetRouter(r *fiber.App) {

    r.Get("/api/jwt", func(c *fiber.Ctx) error {
		loginid := c.Query("loginid")
        passwd := c.Query("passwd")
        return c.JSON(JwtAuth(c, loginid, passwd))
	})

	apiGroup := r.Group("/api")

	apiGroup.Use(JwtAuthRequired)


	// Setup domain-specific routes
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupUserRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupGymRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
}