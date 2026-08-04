package router

import (
	"os"

	"gym/global/log"
	"gym/models"
	"gym/services/billing"

	"github.com/gofiber/fiber/v2"
)

// P5 전자계약 (docs/roadmap-p1-p6.md).
//
// 서명 요청 생성·문서 관리·법적 증빙은 외부 서비스(모두싸인류)가 담당하고,
// 우리는 "어느 주문에 어떤 문서가 어떤 상태인지" 만 기록한다. 외부 API 호출
// (문서 생성)은 관리자 웹이 직접 하고 결과를 여기 등록하는 흐름 — 서버가
// 외부 계약을 대행하는 구조는 서비스 계약 확정 후에 붙인다 (roadmap P5 조건).
func SetupContractRoutes(group fiber.Router) {

	// 계약 기록 생성 (관리자 또는 본인 주문).
	group.Post("/contract", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		var body struct {
			User       int64  `json:"user"`
			Gym        int64  `json:"gym"`
			Order      int64  `json:"order"`
			Provider   string `json:"provider"`
			DocumentId string `json:"documentid"`
			SignUrl    string `json:"signurl"`
		}
		if err := c.BodyParser(&body); err != nil || body.DocumentId == "" || body.Provider == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		if body.User != me.Id && !isGymAdmin(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		conn := models.NewConnection()
		defer conn.Close()
		res, err := conn.Exec(
			`INSERT INTO contract_tb (ct_user, ct_gym, ct_order, ct_provider, ct_document_id, ct_sign_url, ct_status)
			 VALUES (?, ?, ?, ?, ?, ?, 1)`,
			body.User, body.Gym, body.Order, body.Provider, body.DocumentId, body.SignUrl)
		if err != nil {
			log.Error().Msg("contract insert: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		id, _ := res.LastInsertId()
		return c.JSON(fiber.Map{"code": "ok", "id": id})
	})

	// 내 계약 목록.
	group.Get("/contract/my", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		conn := models.NewConnection()
		defer conn.Close()
		rows, err := conn.Query(
			`SELECT ct_id, ct_gym, ct_order, ct_provider, ct_document_id, ct_sign_url, ct_status,
			        IFNULL(ct_signed_date, ''), ct_date
			 FROM contract_tb WHERE ct_user = ? ORDER BY ct_id DESC`, me.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		defer rows.Close()
		items := []fiber.Map{}
		for rows.Next() {
			var id, gym, order int64
			var provider, doc, url, signed, date string
			var status int
			if err := rows.Scan(&id, &gym, &order, &provider, &doc, &url, &status, &signed, &date); err == nil {
				items = append(items, fiber.Map{
					"id": id, "gym": gym, "order": order, "provider": provider,
					"documentid": doc, "signurl": url, "status": status,
					"signeddate": signed, "date": date,
				})
			}
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": len(items)})
	})

	// 외부 서명 서비스 웹훅 — 무인증(isPublicPath) + HMAC 서명 검증.
	// 완료/거절/만료 상태만 반영한다.
	group.Post("/contract/webhook", func(c *fiber.Ctx) error {
		secret := os.Getenv("CONTRACT_WEBHOOK_SECRET")
		sig := c.Get("Webhook-Signature")
		if !billing.VerifyWebhookSignature(c.Body(), sig, secret) {
			log.Info().Str("ip", c.IP()).Msg("contract webhook signature mismatch")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		var body struct {
			Provider   string `json:"provider"`
			DocumentId string `json:"documentid"`
			Event      string `json:"event"` // signed / declined / expired
		}
		if err := c.BodyParser(&body); err != nil || body.DocumentId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		status := 0
		signedExpr := ""
		switch body.Event {
		case "signed":
			status, signedExpr = 2, ", ct_signed_date = NOW()"
		case "declined":
			status = 3
		case "expired":
			status = 4
		default:
			return c.JSON(fiber.Map{"code": "ok"}) // 모르는 이벤트는 무시
		}

		conn := models.NewConnection()
		defer conn.Close()
		_, err := conn.Exec(
			`UPDATE contract_tb SET ct_status = ?`+signedExpr+` WHERE ct_provider = ? AND ct_document_id = ?`,
			status, body.Provider, body.DocumentId)
		if err != nil {
			log.Error().Msg("contract webhook update: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})
}
