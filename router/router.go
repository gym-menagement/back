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

	// 무인증 가입이 관리자 권한을 요구하면 차단한다. JwtAuthRequired 가
	// c.Locals("user") 를 채운 뒤여야 인증된 관리자의 생성과 구분된다.
	apiGroup.Use(BlockPrivilegedSignup)

	// 무인증 아이디 중복확인 응답에서 개인정보를 지운다. 같은 이유로
	// JwtAuthRequired 다음에 온다.
	apiGroup.Use(MaskPublicUserLookup)


	// Setup domain-specific routes
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupMembershipusageRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetupGymtrainerRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
	routers.SetupTrainermemberRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
	routers.SetupUsehealthusageRoutes(apiGroup)
	routers.SetupInquiryRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupNoticeRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetupGymRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupPtreservationRoutes(apiGroup)
	routers.SetupAppversionRoutes(apiGroup)
	routers.SetupQrcodeRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupRockerusageRoutes(apiGroup)
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupUserRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
	routers.SetupMemberqrRoutes(apiGroup)
	routers.SetupWorkoutlogRoutes(apiGroup)
	routers.SetupMemberbodyRoutes(apiGroup)
	routers.SetupPushtokenRoutes(apiGroup)
	routers.SetupAttendanceRoutes(apiGroup)
}