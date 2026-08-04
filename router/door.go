package router

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gym/global/jwt"
	"gym/global/log"
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// P2 출입 도어 제어 (docs/roadmap-p1-p6.md).
//
// 흐름: 회원 앱이 TTL 서명 QR 을 표시 → 리더기가 스캔해 POST /api/door/verify
// → 서버가 QR·회원권 검증 → attendance 기록(+횟수 차감) → open 여부 응답
// → 리더기가 릴레이를 연다. 서버 다운 시 도어는 열리지 않는다(fail-closed).

// qrSecret 은 QR TTL 서명 키. JWT 시크릿에서 파생하되 용도 접두사를 붙여
// JWT 서명과 키를 분리한다 (한쪽 유출이 다른 쪽 위조로 번지지 않게).
func qrSecret() []byte {
	sum := sha256.Sum256([]byte("gym-door-qr:" + jwt.Secret()))
	return sum[:]
}

// MakeDoorQR 은 "uh:<usehealthID>:<만료unix>:<서명>" 형태의 QR 페이로드를 만든다.
// 앱이 30초마다 재발급받아 표시한다 — 스크린샷 재사용을 TTL 로 막는다.
func MakeDoorQR(usehealthID int64, now time.Time, ttl time.Duration) string {
	exp := now.Add(ttl).Unix()
	base := fmt.Sprintf("uh:%d:%d", usehealthID, exp)
	mac := hmac.New(sha256.New, qrSecret())
	mac.Write([]byte(base))
	return base + ":" + hex.EncodeToString(mac.Sum(nil)[:16])
}

// parseDoorQR 은 TTL QR 을 검증한다. 형식·서명·만료 중 하나라도 어긋나면 실패.
func parseDoorQR(payload string, now time.Time) (usehealthID int64, ok bool) {
	parts := strings.Split(payload, ":")
	if len(parts) != 4 || parts[0] != "uh" {
		return 0, false
	}
	id, err1 := strconv.ParseInt(parts[1], 10, 64)
	exp, err2 := strconv.ParseInt(parts[2], 10, 64)
	if err1 != nil || err2 != nil || id <= 0 {
		return 0, false
	}
	base := fmt.Sprintf("uh:%d:%d", id, exp)
	mac := hmac.New(sha256.New, qrSecret())
	mac.Write([]byte(base))
	want := hex.EncodeToString(mac.Sum(nil)[:16])
	if !hmac.Equal([]byte(want), []byte(parts[3])) {
		return 0, false
	}
	if now.Unix() > exp {
		return 0, false
	}
	return id, true
}

// doorSign 은 리더기 요청 서명 — HMAC-SHA256(도어시크릿, 요청 본문).
func doorSign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func SetupDoorRoutes(group fiber.Router) {

	// 회원 앱: TTL QR 발급. 활성 회원권(usehealth)이 있어야 한다.
	group.Get("/door/qr", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		conn := models.NewConnection()
		defer conn.Close()

		// 가장 최근의 사용중 회원권 1건. uh_status 1 = USE(사용중).
		rows, err := conn.Query(
			`SELECT uh_id FROM usehealth_tb
			 WHERE uh_user = ? AND uh_status = 1 AND uh_endday >= NOW()
			 ORDER BY uh_endday DESC LIMIT 1`, me.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		defer rows.Close()
		if !rows.Next() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "no active membership"})
		}
		var uhID int64
		if err := rows.Scan(&uhID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok", "qr": MakeDoorQR(uhID, timeNow(), 30*time.Second), "ttl": 30})
	})

	// 리더기: QR 검증 + 개방 판정. 무인증 경로(isPublicPath) — 도어 시크릿
	// HMAC 이 인증을 대신한다. 실패 응답도 200 에 open:false 로 통일해
	// 리더기 구현을 단순하게 유지한다 (403 은 서명 불일치일 때만).
	group.Post("/door/verify", func(c *fiber.Ctx) error {
		doorID, _ := strconv.ParseInt(c.Get("X-Door-Id"), 10, 64)
		sig := c.Get("X-Door-Sign")
		if doorID <= 0 || sig == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		var gymID int64
		var secret string
		var status int
		rows, err := conn.Query(`SELECT d_gym, d_secret, d_status FROM door_tb WHERE d_id = ?`, doorID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		if !rows.Next() {
			rows.Close()
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		if err := rows.Scan(&gymID, &secret, &status); err != nil {
			rows.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		rows.Close()

		if status != 1 || !hmac.Equal([]byte(doorSign(secret, c.Body())), []byte(sig)) {
			log.Info().Int64("door", doorID).Str("ip", c.IP()).Msg("door verify signature mismatch")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}

		var body struct {
			Qr string `json:"qr"`
		}
		if err := c.BodyParser(&body); err != nil || body.Qr == "" {
			return c.JSON(fiber.Map{"code": "ok", "open": false, "reason": "bad qr"})
		}

		open, reason, userID := verifyAndCheckIn(conn, gymID, body.Qr)
		if !open {
			log.Info().Int64("door", doorID).Str("reason", reason).Msg("door entry denied")
		}
		return c.JSON(fiber.Map{"code": "ok", "open": open, "reason": reason, "user": userID})
	})
}

// verifyAndCheckIn 은 QR → 회원권 검증 → 출석 기록(횟수 차감)까지 수행한다.
// 거부 사유를 문자열로 돌려줘 리더기/관리자 로그에서 원인을 추적할 수 있게 한다.
func verifyAndCheckIn(conn *models.Connection, gymID int64, qr string) (open bool, reason string, userID int64) {
	uhID, ok := parseDoorQR(qr, timeNow())
	var legacy bool
	if !ok {
		// 구버전 앱의 고정 QR(uh_qrcode) 폴백 — 앱 전환 완료 후 제거 예정.
		legacy = true
	}

	var q string
	var arg interface{}
	if legacy {
		q, arg = `uh_qrcode = ?`, qr
	} else {
		q, arg = `uh_id = ?`, uhID
	}

	rows, err := conn.Query(
		`SELECT uh_id, uh_user, uh_gym, uh_status, uh_endday >= NOW(), uh_totalcount, uh_remainingcount
		 FROM usehealth_tb WHERE `+q+` LIMIT 1`, arg)
	if err != nil {
		return false, "db error", 0
	}
	defer rows.Close()
	if !rows.Next() {
		return false, "membership not found", 0
	}
	var id, uid, gid int64
	var status, valid, total, remaining int
	if err := rows.Scan(&id, &uid, &gid, &status, &valid, &total, &remaining); err != nil {
		return false, "db error", 0
	}

	switch {
	case gid != gymID:
		return false, "wrong gym", uid
	case status != 1:
		return false, "membership inactive", uid
	case valid != 1:
		return false, "membership expired", uid
	case total > 0 && remaining <= 0:
		return false, "no remaining count", uid
	}

	// 횟수제면 차감. remaining 재확인 조건을 걸어 동시 스캔 이중 차감을 막는다.
	if total > 0 {
		res, err := conn.Exec(
			`UPDATE usehealth_tb
			 SET uh_usedcount = uh_usedcount + 1, uh_remainingcount = uh_remainingcount - 1,
			     uh_lastuseddate = NOW()
			 WHERE uh_id = ? AND uh_remainingcount > 0`, id)
		if err != nil {
			return false, "db error", uid
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return false, "no remaining count", uid
		}
	} else {
		_, _ = conn.Exec(`UPDATE usehealth_tb SET uh_lastuseddate = NOW() WHERE uh_id = ?`, id)
	}

	_, err = conn.Exec(
		`INSERT INTO attendance_tb (at_user, at_usehealth, at_gym, at_type, at_method, at_status, at_note)
		 VALUES (?, ?, ?, 1, 1, 1, 'door')`, uid, id, gymID)
	if err != nil {
		log.Error().Int64("usehealth", id).Msg("attendance insert failed: " + err.Error())
	}
	return true, "", uid
}

