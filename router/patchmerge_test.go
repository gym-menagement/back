package router

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// PatchMerge 는 부분 수정이 행을 파괴하지 않게 막는 장치라, 경로 판정과
// fail-closed 동작에 회귀가 나면 곧바로 데이터 손실로 이어진다.

func TestPatchMergePathMatching(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		match bool
	}{
		{"엔티티 단건", "/api/user/1", true},
		{"여러자리 id", "/api/usehealth/1234", true},
		{"목록", "/api/user", false},
		{"batch", "/api/user/batch", false},
		{"검색", "/api/user/search/loginid", false},
		{"손코드 중첩경로", "/api/subscription/1/retry", false},
		{"api 밖", "/user/1", false},
		{"id 가 숫자 아님", "/api/user/abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := patchPathRe.MatchString(tt.path) ; got != tt.match {
				t.Errorf("patchPathRe.MatchString(%q) = %v, want %v", tt.path, got, tt.match)
			}
		})
	}
}

// 현재 값을 못 읽으면 통과(fail-open)가 아니라 409 로 막아야 한다.
// 통과시키면 본문에 없는 컬럼이 제로값으로 덮여 행이 파괴된다.
func TestPatchMergeFailsClosed(t *testing.T) {
	app := fiber.New()
	app.Use(PatchMerge)

	updated := false
	// 단건 GET 이 없는(=현재 값을 읽을 수 없는) 엔티티를 흉내낸다.
	app.Patch("/api/thing/:id", func(c *fiber.Ctx) error {
		updated = true
		return c.JSON(fiber.Map{"code": "ok"})
	})

	req := httptest.NewRequest(fiber.MethodPatch, "/api/thing/1",
		strings.NewReader(`{"name":"only-this-field"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusConflict {
		t.Errorf("status = %d, want %d (fail-closed)", resp.StatusCode, fiber.StatusConflict)
	}
	if updated {
		t.Error("병합 불가 상태에서 업데이트 핸들러가 실행되면 안 된다")
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["code"] != "merge unavailable" {
		t.Errorf("code = %v, want \"merge unavailable\"", body["code"])
	}
}

// 본문이 JSON 이 아니거나 비어 있으면 병합 대상이 아니므로 그대로 통과한다.
func TestPatchMergeSkipsNonJSON(t *testing.T) {
	app := fiber.New()
	app.Use(PatchMerge)
	reached := false
	app.Patch("/api/thing/:id", func(c *fiber.Ctx) error {
		reached = true
		return c.JSON(fiber.Map{"code": "ok"})
	})

	req := httptest.NewRequest(fiber.MethodPatch, "/api/thing/1", strings.NewReader(``))
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if !reached {
		t.Error("빈 본문은 병합 대상이 아니므로 통과해야 한다")
	}
}
