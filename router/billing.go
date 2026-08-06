package router

import (
	"os"

	"gym/global/log"
	"gym/models"
	"gym/services/billing"

	"github.com/gofiber/fiber/v2"
)

// SetupBillingRoutes — P1 정기 구독결제 (docs/roadmap-p1-p6.md).
//
// buildtool 이 생성하는 router.go 와 분리된 손코드다. 등록은 SetCustomRouter
// (router/custom.go) 에서 하고, services/http.go 가 SetRouter 다음에 호출하므로
// JwtAuthRequired 를 그대로 통과한다. 웹훅만 isPublicPath 로 열고 서명을 검증한다.
func SetupBillingRoutes(group fiber.Router) {

	// 빌링키 등록 — 앱/웹의 PG 결제창이 발급 성공 후 결과를 보내온다.
	// 본인 것만 만들 수 있다 (user 는 토큰에서, 본문의 user 필드는 무시).
	group.Post("/billing/key", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		var body struct {
			Pg          string `json:"pg"`
			CustomerUid string `json:"customeruid"`
			CardName    string `json:"cardname"`
			CardLast4   string `json:"cardlast4"`
		}
		if err := c.BodyParser(&body); err != nil || body.CustomerUid == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		conn := models.NewConnection()
		defer conn.Close()
		id, err := billing.InsertBillingKey(conn, &billing.BillingKey{
			User: me.Id, Pg: body.Pg, CustomerUid: body.CustomerUid,
			CardName: body.CardName, CardLast4: body.CardLast4,
		})
		if err != nil {
			log.Error().Msg("billing key insert: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok", "id": id})
	})

	// 구독 생성 — 본인 명의 빌링키인지 검증한다.
	group.Post("/subscription", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		var s billing.Subscription
		if err := c.BodyParser(&s); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		if s.Billingkey == 0 || s.Usehealth == 0 || s.Cost <= 0 || s.NextBillingDate == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		s.User = me.Id

		conn := models.NewConnection()
		defer conn.Close()

		key, err := billing.GetBillingKey(conn, s.Billingkey)
		if err != nil || key.User != me.Id || key.Status != billing.BkActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		id, err := billing.InsertSubscription(conn, &s)
		if err != nil {
			log.Error().Msg("subscription insert: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok", "id": id})
	})

	// 내 구독 목록 (회원 앱). 관리자 열람은 아래 /subscription/gym.
	group.Get("/subscription/my", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		conn := models.NewConnection()
		defer conn.Close()
		items, err := billing.FindSubscriptions(conn, `WHERE s_user = ? ORDER BY s_id DESC`, me.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": len(items)})
	})

	// 헬스장별 구독 현황 (관리자 웹). 소유 관계 필수 — 등급만으로는 불가.
	group.Get("/subscription/gym/:id<int>", func(c *fiber.Ctx) error {
		me := currentUser(c)
		conn := models.NewConnection()
		defer conn.Close()
		if !canManageGym(conn, me, gymParamID(c)) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		items, err := billing.FindSubscriptions(conn,
			`WHERE s_gym = ? ORDER BY s_next_billing_date ASC`, gymParamID(c))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": len(items)})
	})

	// 구독 해지 — 본인 또는 관리자.
	group.Delete("/subscription/:id<int>", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		conn := models.NewConnection()
		defer conn.Close()
		id := gymParamID(c)
		s, err := billing.GetSubscription(conn, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
		if s.User != me.Id && !canManageGym(conn, me, s.Gym) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		if err := billing.CancelSubscription(conn, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})

	// 실패 구독 수동 재시도 (관리자). 해당 헬스장 소유 관계 필수.
	group.Post("/subscription/:id<int>/retry", func(c *fiber.Ctx) error {
		me := currentUser(c)
		conn := models.NewConnection()
		defer conn.Close()
		s, err := billing.GetSubscription(conn, gymParamID(c))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
		if !canManageGym(conn, me, s.Gym) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		if err := billing.ChargeSubscription(conn, billing.NewPortOneFromEnv(), s, timeNow()); err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"code": "charge failed", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})

	// PG 웹훅 — 무인증 경로 (isPublicPath 에 등록). 서명 검증이 유일한 방어선.
	// 여기서는 기록만 하고 상태 반영은 하지 않는다 — 상태 변경은 능동 청구
	// (cron/재시도) 결과로만 일어나며, 웹훅은 대사(reconciliation)용 로그다.
	group.Post("/billing/webhook", func(c *fiber.Ctx) error {
		secret := os.Getenv("PORTONE_WEBHOOK_SECRET")
		sig := c.Get("Webhook-Signature")
		if !billing.VerifyWebhookSignature(c.Body(), sig, secret) {
			log.Info().Str("ip", c.IP()).Msg("billing webhook signature mismatch")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		log.Info().Str("body", string(c.Body())).Msg("billing webhook")
		return c.JSON(fiber.Map{"code": "ok"})
	})
}
