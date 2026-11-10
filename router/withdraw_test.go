package router

import (
	"net/http/httptest"
	"strings"
	"testing"

	"gym/models/user"

	"github.com/gofiber/fiber/v2"
)

// 사용자 삭제 권한 회귀. 생성된 DELETE /api/user/:id 는 요청자를 보지 않아
// 로그인한 아무 회원이 남의 계정을 지울 수 있었다. 거부 경로는 DB 를 건드리기
// 전에 끊기므로 DB 없이 검증한다 (허용 경로의 실제 삭제는 DB 가 필요해 제외).

func userDeleteApp(reached *bool) *fiber.App {
	app := fiber.New()
	app.Use(UserDeleteGuard)
	app.All("/*", func(c *fiber.Ctx) error {
		*reached = true
		return c.JSON(fiber.Map{"code": "ok"})
	})
	return app
}

func TestUserDeleteGuard_Blocks(t *testing.T) {
	member := tokenFor(10, user.LevelNormal, user.RoleMember)
	gymAdmin := tokenFor(2, user.LevelAdmin, user.RoleGym_admin)

	cases := []struct {
		name, token, path, body string
		want                    int
	}{
		{"무인증", "", "/api/user/10", "", fiber.StatusUnauthorized},
		{"회원이 남의 계정", member, "/api/user/11", "", fiber.StatusForbidden},
		{"헬스장관리자가 남의 계정", gymAdmin, "/api/user/11", "", fiber.StatusForbidden},
		{"회원 일괄삭제", member, "/api/user/batch", "[11,12]", fiber.StatusForbidden},
		{"헬스장관리자 일괄삭제", gymAdmin, "/api/user/batch", `[{"id":11}]`, fiber.StatusForbidden},
	}
	for _, tc := range cases {
		reached := false
		app := userDeleteApp(&reached)
		req := httptest.NewRequest("DELETE", tc.path, strings.NewReader(tc.body))
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		resp, _ := app.Test(req)
		if resp.StatusCode != tc.want {
			t.Errorf("%s: DELETE %s = %d, want %d", tc.name, tc.path, resp.StatusCode, tc.want)
		}
		if reached {
			t.Errorf("%s: 생성 핸들러까지 도달함", tc.name)
		}
	}
}

func TestUserDeleteGuard_PassThrough(t *testing.T) {
	member := tokenFor(10, user.LevelNormal, user.RoleMember)
	// 삭제가 아닌 메서드와 user 가 아닌 엔티티는 건드리지 않는다.
	for _, tc := range [][2]string{{"GET", "/api/user/11"}, {"PUT", "/api/user/11"}, {"DELETE", "/api/inquiry/3"}} {
		reached := false
		app := userDeleteApp(&reached)
		req := httptest.NewRequest(tc[0], tc[1], nil)
		req.Header.Set("Authorization", "Bearer "+member)
		if resp, _ := app.Test(req); resp.StatusCode != fiber.StatusOK || !reached {
			t.Errorf("%s %s 가 통과하지 않음 (%d)", tc[0], tc[1], resp.StatusCode)
		}
	}
}

func TestBatchUserIDs(t *testing.T) {
	if ids, ok := batchUserIDs([]byte(`[1,2]`)); !ok || len(ids) != 2 || ids[1] != 2 {
		t.Errorf("id 배열 파싱 실패: %v %v", ids, ok)
	}
	if ids, ok := batchUserIDs([]byte(`[{"id":3}]`)); !ok || len(ids) != 1 || ids[0] != 3 {
		t.Errorf("객체 배열 파싱 실패: %v %v", ids, ok)
	}
	for _, bad := range []string{``, `{}`, `[0]`, `[{"id":0}]`, `["x"]`} {
		if _, ok := batchUserIDs([]byte(bad)); ok {
			t.Errorf("잘못된 본문 %q 이 통과함", bad)
		}
	}
}
