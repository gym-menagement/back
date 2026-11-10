package router

import (
	"strconv"

	"gym/global/log"
	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// UserUpdateGuard 는 사용자 수정(PUT/PATCH /api/user/:id)에 대상 경계를 씌운다.
//
// BlockPrivilegeEscalation 은 level/role 이 올라가는지만 본다. 그 밖의 필드는
// 아무 검사가 없어서, 로그인한 회원이 남의 이름·연락처는 물론 비밀번호까지
// 바꿀 수 있었다(계정 탈취). 생성 핸들러는 make server 로 덮어써지므로 다른
// 가드처럼 미들웨어로 막는다.
//
// 규칙:
//   - 본인: 허용
//   - 플랫폼관리자: 허용
//   - 헬스장관리자: 자기보다 낮은 등급(회원·트레이너·직원)만 허용.
//     다른 헬스장관리자·플랫폼관리자의 계정은 못 고친다(비밀번호 교체로
//     위로 올라가는 길 차단).
//   - 그 밖: 403
//
// 남은 한계 — 헬스장관리자가 "자기 헬스장 소속"인지까지는 보지 않는다.
// 회원과 헬스장의 관계가 회원권·트레이너 배정 등 여러 테이블에 흩어져 있어
// 여기서 판정하지 않았다.
var UserUpdateGuard = func(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPut && c.Method() != fiber.MethodPatch {
		return c.Next()
	}
	m := patchPathRe.FindStringSubmatch(c.Path())
	if m == nil || m[1] != "user" {
		return c.Next()
	}
	targetID, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		return c.Next()
	}

	// SetRouter 보다 먼저 돌아 c.Locals("user") 가 없다 — 토큰을 직접 읽는다.
	me := tokenUser(c)
	if me == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
	}
	if me.Id == targetID || isPlatformAdmin(me) {
		return c.Next()
	}
	if isGymAdmin(me) {
		level, role, ok := userPrivilegeLookup(targetID)
		if !ok {
			return c.Next() // 대상 없음 — 뒤 핸들러가 판정
		}
		if level < int(user.LevelAdmin) && role < int(user.RoleGym_admin) {
			return c.Next()
		}
	}

	log.Info().Int64("actor", me.Id).Int64("target", targetID).Msg("user update blocked")
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"code": "forbidden", "message": "다른 사용자의 정보를 수정할 수 없습니다.",
	})
}

// 테스트에서 DB 없이 대상 등급을 바꿔 끼울 수 있게 변수로 둔다.
var userPrivilegeLookup = currentPrivilege
