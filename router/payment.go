package router

import (
	"strconv"
	"strings"

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

	// PATCH /api/payment/:id<int>/refund — 환불 처리 (success → refunded).
	//
	// 돈의 이동은 시스템 밖에서 일어나고(현금 반환, PG 관리자 콘솔) 여기는
	// 장부 상태만 바꾼다. 기본으로 같은 주문(p_order)의 회원권(usehealth)을
	// 종료해 QR 입장과 재등록률 분모에서 빠지게 한다 — 부분 환불처럼
	// 회원권을 유지해야 하면 body 에 keepmembership=true 를 보낸다.
	group.Patch("/payment/:id<int>/refund", func(c *fiber.Ctx) error {
		me := currentUser(c)

		var body struct {
			KeepMembership bool `json:"keepmembership"`
		}
		_ = c.BodyParser(&body) // body 없이 호출해도 기본값으로 동작

		conn := models.NewConnection()
		defer conn.Close()

		id := gymParamID(c)
		rows, err := conn.Query(
			`SELECT p_gym, p_order, p_status FROM payment_tb WHERE p_id = ?`, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		if !rows.Next() {
			rows.Close()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
		var gymID, orderID int64
		var status int
		_ = rows.Scan(&gymID, &orderID, &status)
		rows.Close()

		if !canManageGym(conn, me, gymID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		// 성공 결제만 환불 대상 — 취소·환불·실패 건을 다시 환불하면 장부가 꼬인다.
		if status != billing.PaySuccess {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "invalid status"})
		}

		conn.Begin()
		res, err := conn.Exec(
			`UPDATE payment_tb SET p_status = ? WHERE p_id = ? AND p_status = ?`,
			billing.PayRefunded, id, billing.PaySuccess)
		if err != nil {
			conn.Rollback()
			log.Error().Msg("payment refund update: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		if n, _ := res.RowsAffected(); n == 0 { // 조회 후 다른 요청이 먼저 바꾼 경우
			conn.Rollback()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"code": "conflict"})
		}
		if !body.KeepMembership && orderID > 0 {
			if _, err := conn.Exec(
				`UPDATE usehealth_tb SET uh_status = ?
				 WHERE uh_order = ? AND uh_status IN (?, ?)`,
				billing.UhStatusTerminated, orderID,
				billing.UhStatusUse, billing.UhStatusPaused); err != nil {
				conn.Rollback()
				log.Error().Msg("payment refund terminate usehealth: " + err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
			}
		}
		if err := conn.Commit(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		log.Info().Int64("payment", id).Bool("keepmembership", body.KeepMembership).
			Msg("payment refunded")
		return c.JSON(fiber.Map{"code": "ok"})
	})

	// GET /api/payment/statusmeta?gym=&ids=1,2,3 — 생성 payment 모델에 없는
	// p_status/p_method 를 관리자 목록 화면에 보여주기 위한 임시 라우트.
	// buildtool 재생성으로 두 컬럼이 모델에 들어오면 manual 라우트와 함께 삭제.
	group.Get("/payment/statusmeta", func(c *fiber.Ctx) error {
		me := currentUser(c)
		gymID := int64(c.QueryInt("gym"))

		ids := []interface{}{}
		for _, s := range strings.Split(c.Query("ids"), ",") {
			if v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil && v > 0 {
				ids = append(ids, v)
			}
		}
		if gymID <= 0 || len(ids) == 0 || len(ids) > 200 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		if !canManageGym(conn, me, gymID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		// gym 조건을 같이 걸어 다른 헬스장 결제 상태는 조회되지 않게 한다.
		args := append([]interface{}{gymID}, ids...)
		rows, err := conn.Query(
			`SELECT p_id, p_status, p_method FROM payment_tb
			 WHERE p_gym = ? AND p_id IN (?`+strings.Repeat(",?", len(ids)-1)+`)`, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		defer rows.Close()

		items := []fiber.Map{}
		for rows.Next() {
			var pid int64
			var st, method int
			if err := rows.Scan(&pid, &st, &method); err == nil {
				items = append(items, fiber.Map{"id": pid, "status": st, "method": method})
			}
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": len(items)})
	})

	// GET /api/paymentform/my?page=&pagesize= — 본인 결제 이력.
	//
	// 생성 목록 라우트(GET /paymentform)는 사용자 필터가 없어 RestrictPaymentAccess
	// 가 관리자 미만을 차단한다 — 회원은 이 라우트로만 본다. 생성 매니저를
	// 재사용해 응답 항목 형태(extra 의 gym/payment/paymenttype 포함)를 생성
	// 라우트와 동일하게 유지하고, 생성 모델에 없는 결제 상태·수단만 얹는다.
	group.Get("/paymentform/my", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		page := c.QueryInt("page", 1)
		pagesize := c.QueryInt("pagesize", 10)
		if page < 1 || pagesize < 1 || pagesize > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		manager := models.NewPaymentformManager(conn)
		rows := manager.Find([]interface{}{
			models.Where{Column: "p_user", Value: me.Id, Compare: "="},
			models.Paging(page, pagesize),
			models.Ordering("pf_id desc"),
		})

		// 생성 Count 쿼리는 조인 대상이 빠져 있어 못 쓴다 — 직접 센다.
		var total int64
		cr, err := conn.Query(
			`SELECT COUNT(*) FROM paymentform_tb, payment_tb
			 WHERE pf_payment = p_id AND p_user = ?`, me.Id)
		if err == nil && cr.Next() {
			_ = cr.Scan(&total)
		}
		if cr != nil {
			cr.Close()
		}

		type myPaymentform struct {
			models.Paymentform
			Status int `json:"status"`
			Method int `json:"method"`
		}
		items := make([]myPaymentform, 0, len(rows))
		if len(rows) > 0 {
			ids := make([]interface{}, 0, len(rows))
			for _, r := range rows {
				ids = append(ids, r.Payment)
			}
			meta := map[int64][2]int{}
			mr, err := conn.Query(
				`SELECT p_id, p_status, p_method FROM payment_tb
				 WHERE p_id IN (?`+strings.Repeat(",?", len(ids)-1)+`)`, ids...)
			if err == nil {
				for mr.Next() {
					var pid int64
					var st, method int
					if mr.Scan(&pid, &st, &method) == nil {
						meta[pid] = [2]int{st, method}
					}
				}
				mr.Close()
			}
			for _, r := range rows {
				m := meta[r.Payment]
				items = append(items, myPaymentform{Paymentform: r, Status: m[0], Method: m[1]})
			}
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": total})
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
