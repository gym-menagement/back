package router

import (
	"encoding/json"
	"strconv"

	"gym/global/jwt"
	"gym/global/log"
	"gym/models"
	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// BlockPrivilegeEscalation 은 사용자 수정으로 권한이 올라가는 것을 막는다.
//
// BlockPrivilegedSignup(auth.go)은 **무인증 가입**(POST /api/user)만 본다.
// 인증된 수정 경로는 아무 검증이 없어서, 아무나 회원가입 → 로그인 →
// `PUT /api/user/{내 id} {"level":5,"role":5}` 한 번으로 최고관리자가 됐다
// (2026-08-07 프로세스 E2E 5회차에서 실제로 성공). 가입 상한을 막아둔 의미가
// 사라지는 우회로다.
//
// 규칙 — 요청자 등급을 넘는 권한은 누구에게도 부여할 수 없다:
//   - 플랫폼관리자(Level 4+ / Role 5): 제한 없음
//   - 헬스장관리자(Level 3 / Role 4): 대상 권한을 자기 상한(3/4) 이하로만
//   - 그 밖(회원·트레이너·직원): level/role 을 바꾸는 것 자체 금지
//
// 권한을 낮추는 방향은 막지 않는다 — 강등·비활성은 정상 운영 행위다.
//
// PatchMerge 다음에 등록해야 한다. 병합 후 본문에는 기존 값이 채워져 있으므로
// "본문에 필드가 있는지"가 아니라 "DB 의 현재 값보다 커지는지"로 판정한다.
var BlockPrivilegeEscalation = func(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPut && c.Method() != fiber.MethodPatch {
		return c.Next()
	}
	m := patchPathRe.FindStringSubmatch(c.Path())
	if m == nil || m[1] != "user" {
		return c.Next()
	}

	var body map[string]interface{}
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Next()
	}
	_, hasLevel := body["level"]
	_, hasRole := body["role"]
	if !hasLevel && !hasRole {
		return c.Next()
	}

	targetID, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		return c.Next()
	}

	curLevel, curRole, ok := currentPrivilege(targetID)
	if !ok {
		return c.Next() // 대상 없음 — 뒤 핸들러가 판정
	}

	newLevel, newRole := curLevel, curRole
	if v, ok := privilegeInt(body["level"]); hasLevel && ok {
		newLevel = v
	}
	if v, ok := privilegeInt(body["role"]); hasRole && ok {
		newRole = v
	}
	// 낮추거나 그대로면 통과.
	if newLevel <= curLevel && newRole <= curRole {
		return c.Next()
	}

	// 이 미들웨어는 SetRouter 보다 먼저 돌아서 JwtAuthRequired 가 아직
	// c.Locals("user") 를 채우기 전이다 — 토큰을 직접 파싱해 요청자를 얻는다.
	// (currentUser 를 쓰면 항상 nil 이라 관리자의 정상 수정까지 막힌다)
	me := tokenUser(c)
	if isPlatformAdmin(me) {
		return c.Next()
	}
	// 헬스장관리자는 자기 상한까지만 부여할 수 있다 (가입 상한과 동일 기준).
	if isGymAdmin(me) && newLevel <= maxSignupLevel && newRole <= maxSignupRole {
		return c.Next()
	}

	var meID int64
	if me != nil {
		meID = me.Id
	}
	log.Info().Int64("actor", meID).Int64("target", targetID).
		Int("level", newLevel).Int("role", newRole).Msg("privilege escalation blocked")
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"code": "forbidden", "message": "권한을 올릴 수 없습니다.",
	})
}

// tokenUser 는 Authorization 헤더에서 사용자를 파싱한다. 인증 미들웨어보다
// 앞서 도는 미들웨어용 — 토큰이 없거나 깨졌으면 nil 이고, 그 경우 권한 상승
// 요청은 어차피 차단된다(뒤의 JwtAuthRequired 가 401 을 낸다).
func tokenUser(c *fiber.Ctx) *models.User {
	claims, err := jwt.Check(c.Get("Authorization"))
	if err != nil {
		return nil
	}
	u := claims.User
	return &u
}

func currentPrivilege(id int64) (level, role int, ok bool) {
	conn := models.NewConnection()
	defer conn.Close()
	rows, err := conn.Query(`SELECT u_level, u_role FROM user_tb WHERE u_id = ?`, id)
	if err != nil {
		return 0, 0, false
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, 0, false
	}
	if rows.Scan(&level, &role) != nil {
		return 0, 0, false
	}
	return level, role, true
}

// privilegeInt 는 JSON 의 숫자·문자열 표현을 모두 받는다 — 문자열로 감싸
// 검사를 우회하는 시도를 막는다(auth.go 의 privilegeValue 와 같은 이유).
func privilegeInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case string:
		i, err := strconv.Atoi(n)
		return i, err == nil
	}
	return 0, false
}

var _ = user.LevelAdmin // 등급 상수 출처 표시
