package router

import (
	"net/http/httptest"
	"testing"

	"gym/global/jwt"
	"gym/models"
	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// 개인정보 목록의 수평 권한 상승(BOLA) 차단 회귀. 2026-08-19 전수 조사에서
// 생성 컨트롤러 36개 전부 소유권 필터가 없던 것을 미들웨어로 막았다.

func ownershipApp() *fiber.App {
	app := fiber.New()
	app.Use(OwnershipGuard)
	// 강제된 소유자 파라미터를 응답에 되비춰 검증할 수 있게 한다.
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user":    c.Query("user"),
			"trainer": c.Query("trainer"),
			"passwd":  c.Query("passwd"),
		})
	})
	return app
}

// tokenFor 는 주어진 등급/역할의 사용자 토큰을 만든다 (tokenUser 가 파싱).
func tokenFor(id int64, level user.Level, role user.Role) string {
	return jwt.MakeToken(models.User{Id: id, Level: level, Role: role})
}

func TestOwnershipGuard_AdminOnly(t *testing.T) {
	app := ownershipApp()
	member := tokenFor(10, user.LevelNormal, user.RoleMember)
	admin := tokenFor(1, user.LevelAdmin, user.RoleGym_admin)

	for _, path := range []string{"/api/token", "/api/loginlog", "/api/setting", "/api/user", "/api/user/count", "/api/stop"} {
		// 회원 → 403
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+member)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("member GET %s = %d, want 403", path, resp.StatusCode)
		}
		// 관리자 → 통과
		req = httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+admin)
		resp, _ = app.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("admin GET %s = %d, want 200", path, resp.StatusCode)
		}
	}
}

func TestOwnershipGuard_UserSelfOnly(t *testing.T) {
	app := ownershipApp()
	member := tokenFor(10, user.LevelNormal, user.RoleMember)

	// 본인 단건 → 통과
	req := httptest.NewRequest("GET", "/api/user/10", nil)
	req.Header.Set("Authorization", "Bearer "+member)
	if resp, _ := app.Test(req); resp.StatusCode != fiber.StatusOK {
		t.Errorf("member GET own /api/user/10 = %d, want 200", resp.StatusCode)
	}
	// 남의 단건 → 통과하되 PII 마스킹 (트레이너 이름 조회 등 정상 흐름).
	// 마스킹 자체는 핸들러 응답이 필요하므로 별도 테스트에서 검증한다.
	req = httptest.NewRequest("GET", "/api/user/11", nil)
	req.Header.Set("Authorization", "Bearer "+member)
	if resp, _ := app.Test(req); resp.StatusCode != fiber.StatusOK {
		t.Errorf("member GET other /api/user/11 = %d, want 200 (masked)", resp.StatusCode)
	}
}

func TestOwnershipGuard_MaskOtherUserPII(t *testing.T) {
	app := fiber.New()
	app.Use(OwnershipGuard)
	// 핸들러가 user 전문을 돌려주는 상황을 모사.
	app.Get("/api/user/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"code": "ok", "item": fiber.Map{
			"id": 11, "name": "김트레이너", "email": "t@x.com",
			"tel": "010", "address": "서울", "birth": "1990",
		}})
	})
	member := tokenFor(10, user.LevelNormal, user.RoleMember)
	req := httptest.NewRequest("GET", "/api/user/11", nil)
	req.Header.Set("Authorization", "Bearer "+member)
	resp, _ := app.Test(req)
	body := make([]byte, 512)
	n, _ := resp.Body.Read(body)
	got := string(body[:n])
	if !contains(got, "김트레이너") {
		t.Errorf("표시용 이름이 사라졌다: %s", got)
	}
	for _, leaked := range []string{"t@x.com", "010", "서울", "1990"} {
		if contains(got, leaked) {
			t.Errorf("PII 가 마스킹되지 않았다 (%s): %s", leaked, got)
		}
	}
}

func TestOwnershipGuard_PasswdOracleStripped(t *testing.T) {
	app := ownershipApp()
	admin := tokenFor(1, user.LevelAdmin, user.RoleGym_admin)
	// 관리자여도 passwd 파라미터는 제거돼야 한다 (해시 오라클 차단).
	req := httptest.NewRequest("GET", "/api/user/1?passwd=$2a$10$abc", nil)
	req.Header.Set("Authorization", "Bearer "+admin)
	resp, _ := app.Test(req)
	body := make([]byte, resp.ContentLength)
	_, _ = resp.Body.Read(body)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("admin GET /api/user/1 = %d, want 200", resp.StatusCode)
	}
	if got := string(body); contains(got, "2a$10") {
		t.Errorf("passwd 파라미터가 남았다: %s", got)
	}
}

func TestOwnershipGuard_OwnerScopeForced(t *testing.T) {
	app := ownershipApp()
	member := tokenFor(10, user.LevelNormal, user.RoleMember)

	// 회원이 남의 user 파라미터를 넣어도 본인 ID 로 덮어써진다.
	req := httptest.NewRequest("GET", "/api/order?user=99", nil)
	req.Header.Set("Authorization", "Bearer "+member)
	resp, _ := app.Test(req)
	body := make([]byte, 256)
	n, _ := resp.Body.Read(body)
	if got := string(body[:n]); !contains(got, `"user":"10"`) {
		t.Errorf("order 목록 user 파라미터가 본인으로 강제되지 않았다: %s", got)
	}
}

func TestOwnershipGuard_TrainerScope(t *testing.T) {
	app := ownershipApp()
	trainer := tokenFor(5, user.LevelManager, user.RoleTrainer)

	// 트레이너의 예약 목록은 trainer 파라미터가 본인으로 강제된다.
	req := httptest.NewRequest("GET", "/api/ptreservation?trainer=99", nil)
	req.Header.Set("Authorization", "Bearer "+trainer)
	resp, _ := app.Test(req)
	body := make([]byte, 256)
	n, _ := resp.Body.Read(body)
	if got := string(body[:n]); !contains(got, `"trainer":"5"`) {
		t.Errorf("ptreservation trainer 파라미터가 본인으로 강제되지 않았다: %s", got)
	}
}

func TestOwnershipGuard_PublicCatalogUntouched(t *testing.T) {
	app := ownershipApp()
	member := tokenFor(10, user.LevelNormal, user.RoleMember)

	for _, path := range []string{"/api/gym", "/api/health", "/api/notice", "/api/paymenttype"} {
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+member)
		if resp, _ := app.Test(req); resp.StatusCode != fiber.StatusOK {
			t.Errorf("member GET %s = %d, want 200 (공개 카탈로그)", path, resp.StatusCode)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
