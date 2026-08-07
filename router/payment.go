package router

import (
	"gym/global/log"
	"gym/models"
	"gym/services/billing"

	"github.com/gofiber/fiber/v2"
)

// 수동 결제 등록 (오프라인 결제 흡수, docs/roadmap-p1-p6.md P6).
//
// buildtool 생성 payment 모델에는 P1 에서 추가한 p_status/p_method 컬럼이
// 없어서, 생성 라우트로는 결제수단을 지정할 수 없다(전부 카드·성공으로 기록).
// model.json 에 필드를 추가하고 재생성하는 게 정공법이지만 buildtool 이
// 로컬에 없으므로, 그때까지 손코드 라우트로 채운다.
//
// 이 API 가 있으면 P6(POS 연동)을 미뤄도 현금·계좌이체 매출이 리포트에
// 정확히 잡힌다 — P6 을 보류한 전제가 실제로 성립하게 된다.
func SetupPaymentRoutes(group fiber.Router) {

	// POST /api/payment/manual — 헬스장 관리자가 오프라인 결제를 기록한다.
	group.Post("/payment/manual", func(c *fiber.Ctx) error {
		me := currentUser(c)

		var body struct {
			Gym    int64  `json:"gym"`
			User   int64  `json:"user"`
			Order  int64  `json:"order"`
			Cost   int    `json:"cost"`
			Method int    `json:"method"` // 1:카드 2:계좌이체 3:현금
			Date   string `json:"date"`   // 비우면 현재 시각
			Note   string `json:"note"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		if body.Cost <= 0 || body.User == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		// 정기결제(4)는 PG 청구로만 기록돼야 한다 — 수동으로 만들면 구독
		// 매출과 실제 청구가 어긋난다.
		if body.Method < billing.MethodCard || body.Method > billing.MethodCash {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "invalid method"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		if !canManageGym(conn, me, body.Gym) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		dateExpr, args := "NOW()", []interface{}{
			body.Gym, body.Order, body.User, body.Cost,
			billing.PaySuccess, body.Method, body.Note,
		}
		if body.Date != "" {
			dateExpr = "?"
			args = append(args, body.Date)
		}

		res, err := conn.Exec(
			`INSERT INTO payment_tb
			 (p_gym, p_order, p_user, p_cost, p_status, p_method, p_fail_reason, p_date)
			 VALUES (?, ?, ?, ?, ?, ?, ?, `+dateExpr+`)`, args...)
		if err != nil {
			log.Error().Msg("manual payment insert: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		id, _ := res.LastInsertId()
		log.Info().Int64("payment", id).Int("cost", body.Cost).Int("method", body.Method).
			Msg("manual payment recorded")
		return c.JSON(fiber.Map{"code": "ok", "id": id})
	})

	// PATCH /api/payment/:id<int>/cancel — 오기입 취소 (삭제 대신 상태 변경).
	// 매출 기록은 지우지 않는 게 원칙이라 취소 상태로 남긴다.
	group.Patch("/payment/:id<int>/cancel", func(c *fiber.Ctx) error {
		me := currentUser(c)
		conn := models.NewConnection()
		defer conn.Close()

		id := gymParamID(c)
		rows, err := conn.Query(`SELECT p_gym FROM payment_tb WHERE p_id = ?`, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		if !rows.Next() {
			rows.Close()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
		var gymID int64
		_ = rows.Scan(&gymID)
		rows.Close()

		if !canManageGym(conn, me, gymID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		if _, err := conn.Exec(
			`UPDATE payment_tb SET p_status = ? WHERE p_id = ?`, billing.PayCancelled, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})
}
