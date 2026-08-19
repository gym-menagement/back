package router

import (
	"encoding/json"
	"regexp"
	"strconv"

	"gym/global/log"
	"gym/models"
	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// OwnershipGuard 는 생성 컨트롤러의 목록/단건 조회에 소유권 경계를 씌운다.
//
// 문제: buildtool 이 생성한 Index/Read 핸들러는 요청자 신원(c.Session)을
// 전혀 쓰지 않는다. 필터는 클라이언트가 보낸 user/gym 쿼리 파라미터로만
// 걸리므로, 로그인한 회원이 `GET /api/order?user=<피해자>` 로 남의 개인정보를
// 조회할 수 있었다 (2026-08-19 전수 조사에서 36개 컨트롤러 전부 확인 —
// 수평 권한 상승 BOLA). 생성 파일은 make server 로 덮어써지므로 컨트롤러를
// 고치는 대신 미들웨어로 잠근다 (paymentguard.go 와 같은 전략).
//
// 세 갈래로 처리한다:
//   1) adminOnlyEntities: 비관리자는 목록/단건/카운트 GET 을 403.
//   2) ownerScopedEntities: 비관리자에게는 목록 쿼리의 소유자 파라미터를
//      본인 ID 로 강제하고(변조 무력화), 단건은 행 소유자를 확인한다.
//      앱은 이미 user=본인 으로 호출하므로 정상 흐름은 안 깨진다.
//   3) /api/user 의 passwd 쿼리 파라미터 제거 — LIKE 검색이 bcrypt 해시를
//      한 글자씩 복원하는 오라클이었다(user.go 컨트롤러 + user 모델 %값%).
//
// SetRouter(=JwtAuthRequired) 보다 먼저 등록되므로 c.Locals("user") 가 아직
// 없다. tokenUser(privilegeguard.go) 로 요청자를 직접 파싱한다 — 토큰이 없으면
// 뒤의 JwtAuthRequired 가 401 을 내므로 여기선 관리자 아님으로 취급해 막는다.

// 회원이 목록으로 열람할 이유가 없는 엔티티. 관리자·플랫폼관리자만.
var adminOnlyEntities = map[string]bool{
	"token":               true, // 남의 JWT 원문 → 계정 탈취
	"loginlog":            true,
	"systemlog":           true,
	"ipblock":             true,
	"role":                true,
	"setting":             true,
	"gymadmin":            true,
	"notificationhistory": true, // 지금은 라우트 미등록이나 재생성 대비 선등록
	"notificationsetting": true,
	"user":                true, // 목록/카운트만. 단건(:id)은 아래에서 본인 허용
	"stop":                true, // 소유자 파라미터가 없어 강제 스코프 불가 → 관리자 전용
}

// 본인 것만 봐야 하는 엔티티 → 소유자 컬럼(쿼리 파라미터명)과 단건 소유 확인용
// (테이블, id 컬럼, 소유자 컬럼) 매핑. 소유자 컬럼이 여러 개(member/trainer)면
// 첫 항목을 목록 강제에 쓰고, 단건 확인은 OR 로 본다.
type ownerSpec struct {
	param   string   // 목록 조회에서 강제할 쿼리 파라미터
	table   string   // 단건 소유 확인용 테이블
	idCol   string   // 단건 PK 컬럼
	ownCols []string // 소유자 컬럼(들) — 하나라도 본인이면 통과
}

var ownerScopedEntities = map[string]ownerSpec{
	"order":          {"user", "order_tb", "o_id", []string{"o_user"}},
	"usehealth":      {"user", "usehealth_tb", "uh_id", []string{"uh_user"}},
	"usehealthusage": {"user", "usehealthusage_tb", "uhu_id", []string{"uhu_user"}},
	"memberbody":     {"user", "memberbody_tb", "mb_id", []string{"mb_user"}},
	"inquiry":        {"user", "inquiry_tb", "iq_id", []string{"iq_user"}},
	"attendance":     {"user", "attendance_tb", "at_id", []string{"at_user"}},
	"alarm":          {"user", "alarm_tb", "al_id", []string{"al_user"}},
	"pushtoken":      {"user", "pushtoken_tb", "pt_id", []string{"pt_user"}},
	"workoutlog":     {"user", "workoutlog_tb", "wl_id", []string{"wl_user"}},
	"membership":     {"user", "membership_tb", "m_id", []string{"m_user"}},
	"qrcode":         {"user", "qrcode_tb", "qr_id", []string{"qr_user"}},
	// 예약 관계: 회원(member)·트레이너(trainer) 양쪽이 당사자.
	"ptreservation":  {"member", "ptreservation_tb", "pr_id", []string{"pr_member", "pr_trainer"}},
	"trainermember":  {"member", "trainermember_tb", "tm_id", []string{"tm_member", "tm_trainer"}},
}

var ownershipPathRe = regexp.MustCompile(`^/api/([a-z]+)(?:/(\d+|count))?$`)

var OwnershipGuard = func(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Next() // 쓰기 경로는 다른 가드(privilege/reservation)와 핸들러가 본다
	}
	m := ownershipPathRe.FindStringSubmatch(c.Path())
	if m == nil {
		return c.Next()
	}
	entity, sub := m[1], m[2]

	me := tokenUser(c)
	admin := isGymAdmin(me) || isPlatformAdmin(me)

	// user 의 passwd 오라클 차단 — 관리자 여부와 무관하게 항상 제거.
	if entity == "user" {
		c.Request().URI().QueryArgs().Del("passwd")
	}

	// 단건 GET /api/user/:id — 본인·관리자는 전문, 그 밖은 PII 마스킹.
	// 회원이 트레이너 이름을 보려고 남의 user 를 조회하는 정상 흐름이 있어
	// (app pt_reservation_screen) 막지 않고 민감 필드만 가린다.
	if entity == "user" && sub != "" && sub != "count" {
		if admin {
			return c.Next()
		}
		if me != nil {
			if id, err := strconv.ParseInt(sub, 10, 64); err == nil && id == me.Id {
				return c.Next()
			}
		}
		return maskUserResponse(c)
	}

	if adminOnlyEntities[entity] {
		if admin {
			return c.Next()
		}
		return forbid(c, entity, sub)
	}

	spec, scoped := ownerScopedEntities[entity]
	if !scoped {
		return c.Next() // 공개 카탈로그 등 — 제한 없음
	}
	if admin {
		return c.Next()
	}
	if me == nil {
		return forbid(c, entity, sub)
	}
	// 트레이너는 예약/담당 관계의 당사자라 목록을 봐야 한다 — 강제 스코프에서 제외
	// 하되, 자기 것만 보도록 trainer 파라미터를 본인으로 강제한다.
	if sub == "" || sub == "count" {
		scopeParam := spec.param
		if me.Role == user.RoleTrainer && len(spec.ownCols) > 1 {
			scopeParam = "trainer"
		}
		c.Request().URI().QueryArgs().Set(scopeParam, strconv.FormatInt(me.Id, 10))
		return c.Next()
	}
	// 단건 — 행 소유자를 조회해 본인인지 확인.
	if !rowOwnedBy(spec, sub, me.Id) {
		return forbid(c, entity, sub)
	}
	return c.Next()
}

// rowOwnedBy 는 단건 행의 소유자 컬럼 중 하나라도 uid 와 같으면 참.
func rowOwnedBy(spec ownerSpec, idStr string, uid int64) bool {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return false
	}
	conn := models.NewConnection()
	defer conn.Close()

	where := spec.ownCols[0] + " = ?"
	args := []interface{}{id, uid}
	for _, col := range spec.ownCols[1:] {
		where += " OR " + col + " = ?"
		args = append(args, uid)
	}
	rows, err := conn.Query(
		"SELECT 1 FROM "+spec.table+" WHERE "+spec.idCol+" = ? AND ("+where+") LIMIT 1", args...)
	if err != nil {
		return false // 판정 불가면 막는다 (fail-closed)
	}
	defer rows.Close()
	return rows.Next()
}

// nonSelfUserFields — 남의 user 를 조회할 때 남겨도 되는 필드(표시용).
// 나머지(email·tel·address·birth·loginid·connectid·logindate 등)는 버린다.
var nonSelfUserFields = map[string]bool{
	"id": true, "name": true, "image": true, "sex": true,
	"level": true, "role": true, "use": true,
}

// maskUserResponse 는 핸들러 응답의 item 에서 민감 필드를 제거한다.
func maskUserResponse(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	var payload map[string]any
	if json.Unmarshal(c.Response().Body(), &payload) != nil {
		return nil
	}
	item, ok := payload["item"].(map[string]any)
	if !ok {
		return nil
	}
	for k := range item {
		if !nonSelfUserFields[k] {
			delete(item, k)
		}
	}
	payload["item"] = item
	if body, err := json.Marshal(payload); err == nil {
		c.Response().SetBody(body)
	}
	return nil
}

func forbid(c *fiber.Ctx, entity, sub string) error {
	var meID int64
	if me := tokenUser(c); me != nil {
		meID = me.Id
	}
	log.Info().Int64("actor", meID).Str("entity", entity).Str("sub", sub).
		Msg("ownership guard blocked")
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
}
