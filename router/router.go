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
	routers.SetupAttendanceRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupWorkoutlogRoutes(apiGroup)
	routers.SetupGymRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupMemberqrRoutes(apiGroup)
	routers.SetupPtreservationRoutes(apiGroup)
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupRockerusageRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupPushtokenRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupInquiryRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetupAppversionRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupNoticeRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
	routers.SetupTrainermemberRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetupMemberbodyRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupMembershipusageRoutes(apiGroup)
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupUserRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
}