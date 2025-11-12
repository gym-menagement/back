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
	routers.SetupUserRoutes(apiGroup)
	routers.SetupMemberqrRoutes(apiGroup)
	routers.SetupRockerusageRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupMemberbodyRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupPtreservationRoutes(apiGroup)
	routers.SetupMembershipusageRoutes(apiGroup)
	routers.SetupPushtokenRoutes(apiGroup)
	routers.SetupWorkoutlogRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupTrainermemberRoutes(apiGroup)
	routers.SetupAttendanceRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupAppversionRoutes(apiGroup)
	routers.SetupGymRoutes(apiGroup)
	routers.SetupNoticeRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupInquiryRoutes(apiGroup)
}