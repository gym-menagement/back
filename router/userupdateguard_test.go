package router

import (
	"net/http/httptest"
	"strings"
	"testing"

	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// 사용자 수정 대상 경계 회귀. 회원이 남의 프로필·비밀번호를 바꿀 수 있었다.

func TestUserUpdateGuard(t *testing.T) {
	// 대상 등급: 11 = 회원, 20 = 헬스장관리자, 30 = 플랫폼관리자, 99 = 없음
	orig := userPrivilegeLookup
	userPrivilegeLookup = func(id int64) (int, int, bool) {
		switch id {
		case 11:
			return int(user.LevelNormal), int(user.RoleMember), true
		case 20:
			return int(user.LevelAdmin), int(user.RoleGym_admin), true
		case 30:
			return int(user.LevelSuperadmin), int(user.RolePlatform_admin), true
		}
		return 0, 0, false
	}
	defer func() { userPrivilegeLookup = orig }()

	member := tokenFor(10, user.LevelNormal, user.RoleMember)
	gymAdmin := tokenFor(2, user.LevelAdmin, user.RoleGym_admin)
	platform := tokenFor(1, user.LevelSuperadmin, user.RolePlatform_admin)

	cases := []struct {
		name, token, method, path string
		want                      int
	}{
		{"무인증", "", "PUT", "/api/user/10", fiber.StatusUnauthorized},
		{"회원 본인", member, "PUT", "/api/user/10", fiber.StatusOK},
		{"회원 본인 PATCH", member, "PATCH", "/api/user/10", fiber.StatusOK},
		{"회원이 남의 계정", member, "PUT", "/api/user/11", fiber.StatusForbidden},
		{"회원이 남의 계정 PATCH", member, "PATCH", "/api/user/11", fiber.StatusForbidden},
		{"회원이 관리자 계정", member, "PUT", "/api/user/30", fiber.StatusForbidden},
		{"헬스장관리자가 회원", gymAdmin, "PUT", "/api/user/11", fiber.StatusOK},
		{"헬스장관리자가 다른 헬스장관리자", gymAdmin, "PUT", "/api/user/20", fiber.StatusForbidden},
		{"헬스장관리자가 플랫폼관리자", gymAdmin, "PATCH", "/api/user/30", fiber.StatusForbidden},
		{"헬스장관리자, 대상 없음", gymAdmin, "PUT", "/api/user/99", fiber.StatusOK},
		{"플랫폼관리자가 아무나", platform, "PUT", "/api/user/20", fiber.StatusOK},
		{"다른 엔티티는 통과", member, "PUT", "/api/inquiry/3", fiber.StatusOK},
		{"조회는 통과", member, "GET", "/api/user/11", fiber.StatusOK},
	}
	for _, tc := range cases {
		reached := false
		app := fiber.New()
		app.Use(UserUpdateGuard)
		app.All("/*", func(c *fiber.Ctx) error {
			reached = true
			return c.JSON(fiber.Map{"code": "ok"})
		})
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{"name":"x"}`))
		req.Header.Set("Content-Type", "application/json")
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		resp, _ := app.Test(req)
		if resp.StatusCode != tc.want {
			t.Errorf("%s: %s %s = %d, want %d", tc.name, tc.method, tc.path, resp.StatusCode, tc.want)
		}
		if (tc.want == fiber.StatusOK) != reached {
			t.Errorf("%s: 핸들러 도달 여부가 기대와 다름 (reached=%v)", tc.name, reached)
		}
	}
}
