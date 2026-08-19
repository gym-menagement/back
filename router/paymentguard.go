package router

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RestrictPaymentAccess 는 결제 데이터(payment/paymentform)를 관리자 전용으로
// 잠근다.
//
// 생성 컨트롤러의 목록 조회에는 사용자 필터가 없어서, 로그인한 회원 누구나
// `GET /api/paymentform?page=1&pagesize=10` 로 **전체 회원의 결제 기록**을
// 볼 수 있었다 (2026-08-19 확인 — 회원 웹 결제 이력 화면이 실제로 이렇게
// 호출하고 있었다). 생성 파일은 buildtool 재생성 때 덮어써지므로 컨트롤러
// 수정 대신 미들웨어로 막는다.
//
// 규칙 — 헬스장관리자 미만(회원·트레이너·직원)은:
//   - POST(구매 흐름의 결제/결제양식 생성)만 허용
//   - GET/PUT/PATCH/DELETE 는 403 (fail-closed)
//   - 본인 결제 이력은 전용 라우트 GET /api/paymentform/my 로 조회
//
// paymenttype 은 상품 유형 목록이라 제한하지 않는다 (경로 접두사가 겹치므로
// 세그먼트 단위로 판정). BlockPrivilegeEscalation 처럼 SetRouter 이전에
// 등록되어 Locals("user") 가 아직 없다 — tokenUser 로 요청자를 판별한다.
var RestrictPaymentAccess = func(c *fiber.Ctx) error {
	path := c.Path()
	if !paymentEntityPath(path) {
		return c.Next()
	}
	if c.Method() == fiber.MethodPost {
		return c.Next()
	}
	// 회원 본인 이력 조회는 전용 핸들러가 자체적으로 p_user = 나 로 제한한다.
	if path == "/api/paymentform/my" {
		return c.Next()
	}
	me := tokenUser(c)
	if isGymAdmin(me) || isPlatformAdmin(me) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
}

// paymentEntityPath 는 /api/payment, /api/paymentform (및 하위 경로)만 참.
// /api/paymenttype 처럼 접두사만 겹치는 엔티티는 제외한다.
func paymentEntityPath(path string) bool {
	for _, base := range []string{"/api/payment", "/api/paymentform"} {
		if path == base || strings.HasPrefix(path, base+"/") {
			return true
		}
	}
	return false
}
