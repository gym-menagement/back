package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// 결제 데이터 접근 제한의 회귀. 생성 목록 라우트에 사용자 필터가 없어
// 로그인한 회원 누구나 전체 결제 기록을 볼 수 있던 구멍(2026-08-19)을
// 다시 열지 않도록 고정한다.

func TestPaymentEntityPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/api/payment", true},
		{"/api/payment/1", true},
		{"/api/payment/count", true},
		{"/api/paymentform", true},
		{"/api/paymentform/1", true},
		{"/api/paymentform/my", true}, // 매칭은 되지만 미들웨어에서 예외 처리
		{"/api/paymenttype", false},   // 상품 유형 — 접두사만 겹침
		{"/api/paymenttype/1", false},
		{"/api/user", false},
		{"/api/order", false},
	}
	for _, tt := range tests {
		if got := paymentEntityPath(tt.path); got != tt.want {
			t.Errorf("paymentEntityPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

// 토큰 없는(=관리자 아님) 요청 기준의 미들웨어 동작. 관리자 통과 분기는
// isGymAdmin/isPlatformAdmin 을 그대로 쓰므로 별도 검증하지 않는다.
func TestRestrictPaymentAccess(t *testing.T) {
	app := fiber.New()
	app.Use(RestrictPaymentAccess)
	app.All("/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{"GET", "/api/paymentform", fiber.StatusForbidden},
		{"GET", "/api/paymentform/1", fiber.StatusForbidden},
		{"GET", "/api/payment", fiber.StatusForbidden},
		{"GET", "/api/payment/count", fiber.StatusForbidden},
		{"PUT", "/api/payment/1", fiber.StatusForbidden},
		{"PATCH", "/api/payment/1", fiber.StatusForbidden},
		{"DELETE", "/api/payment/1", fiber.StatusForbidden},
		{"POST", "/api/payment", fiber.StatusOK},        // 구매 흐름 (앱)
		{"POST", "/api/paymentform", fiber.StatusOK},    // 구매 흐름 (앱)
		{"GET", "/api/paymentform/my", fiber.StatusOK},  // 본인 이력 전용 라우트
		{"GET", "/api/paymenttype", fiber.StatusOK},     // 제한 대상 아님
		{"GET", "/api/order", fiber.StatusOK},           // 무관 엔티티
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("%s %s: %v", tt.method, tt.path, err)
		}
		if resp.StatusCode != tt.want {
			t.Errorf("%s %s = %d, want %d", tt.method, tt.path, resp.StatusCode, tt.want)
		}
	}
}
