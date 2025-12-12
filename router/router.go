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
	routers.SetupGymtrainerRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupAttendanceRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupGymRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupMemberqrRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupNoticeRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupInquiryRoutes(apiGroup)
	routers.SetupQrcodeRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetupPtreservationRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupUsehealthusageRoutes(apiGroup)
	routers.SetupWorkoutlogRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupTrainermemberRoutes(apiGroup)
	routers.SetupMembershipusageRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupRockerusageRoutes(apiGroup)
	routers.SetupMemberbodyRoutes(apiGroup)
	routers.SetupAppversionRoutes(apiGroup)
	routers.SetupUserRoutes(apiGroup)
	routers.SetupPushtokenRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
}