package router

import (
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// P4 CEO 리포트 (docs/roadmap-p1-p6.md).
// 신규 테이블 없이 원본(payment/order/usehealth/attendance) 집계로 시작한다.
// 느려지면 그때 일별 스냅샷을 도입한다 — 선최적화 금지.
func SetupReportRoutes(group fiber.Router) {

	// GET /api/report/summary?gym=&from=&to=  (from/to: YYYY-MM-DD)
	// gym_admin 은 자기 헬스장만, platform_admin 은 전체(gym=0 허용).
	group.Get("/report/summary", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if !isGymAdmin(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		gymID := int64(c.QueryInt("gym"))
		from := c.Query("from")
		to := c.Query("to")
		if from == "" || to == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		if gymID == 0 && !isPlatformAdmin(me) {
			// 전체 집계는 플랫폼관리자만. 헬스장관리자는 gym 파라미터 필수.
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		gymWhere := ""
		params := []interface{}{from, to}
		if gymID > 0 {
			gymWhere = " AND p_gym = ?"
			params = append(params, gymID)
		}

		// 1) 월별 매출 (성공 결제만)
		monthly := []fiber.Map{}
		rows, err := conn.Query(
			`SELECT DATE_FORMAT(p_date, '%Y-%m') ym, SUM(p_cost), COUNT(*)
			 FROM payment_tb
			 WHERE p_status = 1 AND DATE(p_date) BETWEEN ? AND ?`+gymWhere+`
			 GROUP BY ym ORDER BY ym`, params...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		for rows.Next() {
			var ym string
			var sum, cnt int64
			if err := rows.Scan(&ym, &sum, &cnt); err == nil {
				monthly = append(monthly, fiber.Map{"month": ym, "revenue": sum, "count": cnt})
			}
		}
		rows.Close()

		// 2) 결제수단별 매출
		byMethod := []fiber.Map{}
		rows, err = conn.Query(
			`SELECT p_method, SUM(p_cost), COUNT(*)
			 FROM payment_tb
			 WHERE p_status = 1 AND DATE(p_date) BETWEEN ? AND ?`+gymWhere+`
			 GROUP BY p_method`, params...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		for rows.Next() {
			var method int
			var sum, cnt int64
			if err := rows.Scan(&method, &sum, &cnt); err == nil {
				byMethod = append(byMethod, fiber.Map{"method": method, "revenue": sum, "count": cnt})
			}
		}
		rows.Close()

		// 3) 신규 회원권 발급 수 / 4) 재등록률 — 기간 내 만료 회원권 중
		// 만료 후 30일 안에 같은 회원이 새 회원권을 산 비율.
		uhWhere := ""
		uhParams := []interface{}{from, to}
		if gymID > 0 {
			uhWhere = " AND uh_gym = ?"
			uhParams = append(uhParams, gymID)
		}

		var newCount int64
		rows, err = conn.Query(
			`SELECT COUNT(*) FROM usehealth_tb WHERE DATE(uh_date) BETWEEN ? AND ?`+uhWhere, uhParams...)
		if err == nil && rows.Next() {
			_ = rows.Scan(&newCount)
		}
		if rows != nil {
			rows.Close()
		}

		var expired, renewed int64
		rows, err = conn.Query(
			`SELECT COUNT(*),
			        CAST(IFNULL(SUM(EXISTS(
			          SELECT 1 FROM usehealth_tb n
			          WHERE n.uh_user = e.uh_user AND n.uh_gym = e.uh_gym
			            AND n.uh_id <> e.uh_id
			            AND n.uh_date > e.uh_endday
			            AND n.uh_date <= DATE_ADD(e.uh_endday, INTERVAL 30 DAY))), 0) AS SIGNED)
			 FROM usehealth_tb e
			 WHERE DATE(e.uh_endday) BETWEEN ? AND ?`+uhWhere, uhParams...)
		if err == nil && rows.Next() {
			_ = rows.Scan(&expired, &renewed)
		}
		if rows != nil {
			rows.Close()
		}

		// 5) 회원권 판매 순위 (주문 기준 상위 10)
		topHealth := []fiber.Map{}
		oWhere := ""
		oParams := []interface{}{from, to}
		if gymID > 0 {
			oWhere = " AND o_gym = ?"
			oParams = append(oParams, gymID)
		}
		rows, err = conn.Query(
			`SELECT o.o_health, IFNULL(h.h_name, ''), COUNT(*)
			 FROM order_tb o LEFT JOIN health_tb h ON h.h_id = o.o_health
			 WHERE DATE(o.o_date) BETWEEN ? AND ?`+oWhere+`
			 GROUP BY o.o_health, h.h_name ORDER BY COUNT(*) DESC LIMIT 10`, oParams...)
		if err == nil {
			for rows.Next() {
				var hid int64
				var name string
				var cnt int64
				if err := rows.Scan(&hid, &name, &cnt); err == nil {
					topHealth = append(topHealth, fiber.Map{"health": hid, "name": name, "count": cnt})
				}
			}
			rows.Close()
		}

		// 6) 출석 히트맵 (요일 0=일 ~ 6=토 × 시간대)
		heatmap := []fiber.Map{}
		aWhere := ""
		aParams := []interface{}{from, to}
		if gymID > 0 {
			aWhere = " AND at_gym = ?"
			aParams = append(aParams, gymID)
		}
		rows, err = conn.Query(
			`SELECT DAYOFWEEK(at_checkintime)-1, HOUR(at_checkintime), COUNT(*)
			 FROM attendance_tb
			 WHERE at_type = 1 AND DATE(at_checkintime) BETWEEN ? AND ?`+aWhere+`
			 GROUP BY 1, 2`, aParams...)
		if err == nil {
			for rows.Next() {
				var dow, hour int
				var cnt int64
				if err := rows.Scan(&dow, &hour, &cnt); err == nil {
					heatmap = append(heatmap, fiber.Map{"dow": dow, "hour": hour, "count": cnt})
				}
			}
			rows.Close()
		}

		renewRate := 0.0
		if expired > 0 {
			renewRate = float64(renewed) / float64(expired)
		}

		return c.JSON(fiber.Map{
			"code": "ok",
			"item": fiber.Map{
				"monthly":     monthly,
				"bymethod":    byMethod,
				"newmembers":  newCount,
				"expired":     expired,
				"renewed":     renewed,
				"renewrate":   renewRate,
				"tophealth":   topHealth,
				"attendance":  heatmap,
			},
		})
	})
}
