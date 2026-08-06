package router

import (
	"strconv"
	"time"

	"gym/models"
	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// SetCustomRouter 는 buildtool 이 생성하지 않는 손코드 라우트(P1~P5)를 등록한다.
//
// services/http.go 가 router.SetRouter(app) 다음에 호출한다 — SetRouter 가
// 등록한 JwtAuthRequired 미들웨어가 /api 접두사 전체에 걸려 있으므로 여기
// 라우트들도 그대로 인증을 통과해야 도달한다. 무인증이 필요한 경로(웹훅,
// 도어 verify)는 isPublicPath(auth.go) 완전 일치 목록에 추가하고 각 핸들러가
// 자체 서명 검증을 한다.
//
// buildtool 재생성(make server) 시에도 이 파일과 각 Setup*Routes 파일은
// 남는다 — router.go 만 덮어써진다.
func SetCustomRouter(app *fiber.App) {
	group := app.Group("/api")

	SetupBillingRoutes(group)  // P1 정기 구독결제
	SetupDoorRoutes(group)     // P2 출입 도어 제어
	SetupScheduleRoutes(group) // P3 트레이너 근무 스케줄
	SetupReportRoutes(group)   // P4 CEO 리포트
	SetupContractRoutes(group) // P5 전자계약
}

// ---- 공용 헬퍼 ----

// currentUser 는 JwtAuthRequired 가 채운 인증 사용자를 돌려준다. 무인증
// 공개 경로로 들어온 요청이면 nil.
func currentUser(c *fiber.Ctx) *models.User {
	u, _ := c.Locals("user").(*models.User)
	return u
}

// isGymAdmin — 헬스장관리자 이상. Level 과 Role 은 같은 숫자가 다른 권한을
// 뜻하므로(auth.go 상단 주석) 반드시 둘 다 본다: Level 3+/Role 4+ 가 관리자.
func isGymAdmin(u *models.User) bool {
	if u == nil {
		return false
	}
	return u.Level >= user.LevelAdmin || u.Role >= user.RoleGym_admin
}

// isPlatformAdmin — 플랫폼관리자 이상 (Level 4+/Role 5).
func isPlatformAdmin(u *models.User) bool {
	if u == nil {
		return false
	}
	return u.Level >= user.LevelSuperadmin || u.Role >= user.RolePlatform_admin
}

// canManageGym 은 해당 헬스장의 관리 데이터(매출 리포트·구독 현황)에 접근할
// 자격을 판정한다. 등급(level 3)만 보면 안 된다 — 가입 정책이 헬스장관리자
// 자가가입을 허용하므로, 등급만 확인하면 아무나 가입해서 남의 헬스장 매출을
// 볼 수 있다 (2026-08-06 테스트에서 실제 확인된 구멍).
//
// 허용 조건: 플랫폼관리자 | gym_tb.g_user 소유자 | gymadmin_tb 등록 관리자.
func canManageGym(conn *models.Connection, u *models.User, gymID int64) bool {
	if u == nil || gymID <= 0 {
		return false
	}
	if isPlatformAdmin(u) {
		return true
	}
	if !isGymAdmin(u) {
		return false
	}
	rows, err := conn.Query(
		`SELECT 1 FROM gym_tb WHERE g_id = ? AND g_user = ?
		 UNION SELECT 1 FROM gymadmin_tb WHERE ga_gym = ? AND ga_user = ? LIMIT 1`,
		gymID, u.Id, gymID, u.Id)
	if err != nil {
		return false // 판정 불가면 막는다 (fail-closed)
	}
	defer rows.Close()
	return rows.Next()
}

func gymParamID(c *fiber.Ctx) int64 {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	return id
}

// timeNow 는 테스트에서 갈아끼우기 위한 간접 호출.
var timeNow = time.Now
