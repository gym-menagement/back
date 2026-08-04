package router

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// isPublicPath 는 무인증 통과를 결정하는 보안 경계라서 회귀가 나면 바로
// 데이터 노출로 이어진다. 예전 구현이 "/api/user" 접두사를 통째로 열어
// 전체 회원 명단이 무인증 조회됐던 사고를 다시 내지 않기 위한 테스트.
func TestIsPublicPath(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		// 열려 있어야 하는 것 — 로그인 전에 반드시 필요하다.
		{"로그인", "GET", "/api/jwt", true},
		{"회원가입", "POST", "/api/user", true},
		{"아이디 중복확인", "GET", "/api/user/search/loginid", true},

		// 과거 사고 지점: 접두사 매칭이라 전부 뚫려 있었다.
		{"회원 목록 조회", "GET", "/api/user", false},
		{"회원 단건 조회", "GET", "/api/user/1", false},
		{"회원 수정", "PUT", "/api/user/1", false},
		{"회원 삭제", "DELETE", "/api/user/1", false},
		{"회원 일괄등록", "POST", "/api/user/batch", false},
		{"이메일로 검색", "GET", "/api/user/search/email", false},
		{"전화번호로 검색", "GET", "/api/user/search/tel", false},
		{"로그인아이디 단건", "GET", "/api/user/get/loginid/admin", false},

		// 9자 미만 경로 — 예전엔 path[:9] 가 여기서 패닉했다.
		{"짧은 경로", "GET", "/api/gym", false},
		{"더 짧은 경로", "GET", "/api", false},
		{"루트", "GET", "/", false},

		// 메서드가 다르면 막혀야 한다.
		{"jwt 에 POST", "POST", "/api/jwt", false},
		{"가입 경로에 GET", "GET", "/api/user", false},
		{"중복확인 경로에 POST", "POST", "/api/user/search/loginid", false},

		// 접두사·접미사로 우회 시도.
		{"뒤에 슬래시", "POST", "/api/user/", false},
		{"경로 덧붙이기", "POST", "/api/users", false},
		{"중복확인 뒤에 덧붙이기", "GET", "/api/user/search/loginid/all", false},

		// P1~P5 기계 대 기계 경로 — JWT 대신 자체 서명(HMAC)을 검증한다.
		{"PG 웹훅", "POST", "/api/billing/webhook", true},
		{"계약 웹훅", "POST", "/api/contract/webhook", true},
		{"도어 verify", "POST", "/api/door/verify", true},
		// 메서드·경로 변형은 전부 막는다.
		{"PG 웹훅에 GET", "GET", "/api/billing/webhook", false},
		{"도어 verify 에 GET", "GET", "/api/door/verify", false},
		{"도어 QR 발급은 인증 필요", "GET", "/api/door/qr", false},
		{"빌링키 등록은 인증 필요", "POST", "/api/billing/key", false},
		{"구독 생성은 인증 필요", "POST", "/api/subscription", false},
		{"웹훅 경로 덧붙이기", "POST", "/api/billing/webhook/x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPublicPath(tt.method, tt.path); got != tt.want {
				t.Errorf("isPublicPath(%q, %q) = %v, want %v", tt.method, tt.path, got, tt.want)
			}
		})
	}
}

// 가입 요청만으로 최고관리자가 되는 권한 상승을 막는지 확인한다.
// JSON 만 검사하면 form 인코딩으로 그대로 우회되므로 두 경로를 모두 본다.
func TestBlockPrivilegedSignup(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantBlocked bool
	}{
		// 정상 가입 — 웹 가입 화면이 실제로 보내는 값들.
		{"회원 가입", fiber.MIMEApplicationJSON, `{"loginid":"a","role":1,"level":1}`, false},
		{"트레이너 가입", fiber.MIMEApplicationJSON, `{"loginid":"a","role":2,"level":2}`, false},
		{"헬스장관리자 가입", fiber.MIMEApplicationJSON, `{"loginid":"a","role":4,"level":3}`, false},

		// Level 과 Role 은 라벨 순서가 어긋나 같은 숫자가 같은 권한이 아니다.
		// Level 4 는 플랫폼관리자라 막고, Role 4 는 헬스장관리자라 연다.
		{"level 4 는 플랫폼관리자", fiber.MIMEApplicationJSON, `{"loginid":"a","level":4}`, true},
		{"권한 필드 없음", fiber.MIMEApplicationJSON, `{"loginid":"a"}`, false},
		{"빈 본문", fiber.MIMEApplicationJSON, ``, false},

		// 권한 상승 시도.
		{"role 5 요구", fiber.MIMEApplicationJSON, `{"loginid":"a","role":5}`, true},
		{"level 5 요구", fiber.MIMEApplicationJSON, `{"loginid":"a","level":5}`, true},
		{"둘 다 5", fiber.MIMEApplicationJSON, `{"loginid":"a","role":5,"level":5}`, true},
		{"범위 밖 큰 값", fiber.MIMEApplicationJSON, `{"loginid":"a","role":99}`, true},

		// 문자열로 감싸 보내는 우회.
		{"문자열 role", fiber.MIMEApplicationJSON, `{"loginid":"a","role":"5"}`, true},

		// form 인코딩 우회 — JSON 파싱은 실패하지만 핸들러는 값을 읽는다.
		{"form 정상", fiber.MIMEApplicationForm, `loginid=a&role=1`, false},
		{"form 으로 role 5", fiber.MIMEApplicationForm, `loginid=a&role=5`, true},
		{"form 으로 level 5", fiber.MIMEApplicationForm, `loginid=a&level=5`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(BlockPrivilegedSignup)
			app.Post("/api/user", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest(fiber.MethodPost, "/api/user", strings.NewReader(tt.body))
			req.Header.Set(fiber.HeaderContentType, tt.contentType)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("요청 실패: %v", err)
			}

			blocked := resp.StatusCode == fiber.StatusForbidden
			if blocked != tt.wantBlocked {
				t.Errorf("status=%d, blocked=%v, want blocked=%v", resp.StatusCode, blocked, tt.wantBlocked)
			}
		})
	}
}

// 인증을 통과한 관리자가 만드는 계정에는 상한을 적용하지 않는다.
func TestBlockPrivilegedSignupSkipsAuthenticated(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &struct{ Name string }{"admin"}) // JwtAuthRequired 통과 상황 재현
		return c.Next()
	})
	app.Use(BlockPrivilegedSignup)
	app.Post("/api/user", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodPost, "/api/user", strings.NewReader(`{"role":5,"level":5}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("요청 실패: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("인증된 관리자의 생성은 통과해야 한다: status=%d", resp.StatusCode)
	}
}

// 가입 외의 경로는 이 미들웨어가 건드리지 않는다.
func TestBlockPrivilegedSignupIgnoresOtherRoutes(t *testing.T) {
	app := fiber.New()
	app.Use(BlockPrivilegedSignup)
	app.Put("/api/user/1", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/gym", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	for _, r := range []struct{ method, path string }{
		{fiber.MethodPut, "/api/user/1"},
		{fiber.MethodPost, "/api/gym"},
	} {
		req := httptest.NewRequest(r.method, r.path, strings.NewReader(`{"role":5}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("요청 실패: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("%s %s 는 통과해야 한다: status=%d", r.method, r.path, resp.StatusCode)
		}
	}
}

// 무인증 아이디 중복확인이 개인정보를 흘리지 않는지 확인한다.
// 이 경로는 가입 흐름 때문에 열어둬야 해서 응답 쪽에서 막는 수밖에 없다.
func TestMaskPublicUserLookup(t *testing.T) {
	const full = `{"code":"ok","item":{"id":1,"loginid":"admin","email":"admin@gym.com",` +
		`"name":"최고관리자","tel":"010-0001-0001","address":"서울시 강남구",` +
		`"birth":"1980-01-01","level":5,"role":5}}`

	newApp := func(authenticated bool) *fiber.App {
		app := fiber.New()
		if authenticated {
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("user", &struct{ Name string }{"admin"})
				return c.Next()
			})
		}
		app.Use(MaskPublicUserLookup)
		h := func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.SendString(full)
		}
		app.Get("/api/user/search/loginid", h)
		app.Get("/api/user/search/email", h) // 인증 필요 경로 — 미들웨어가 건드리면 안 된다
		return app
	}

	read := func(t *testing.T, app *fiber.App, path string) map[string]any {
		t.Helper()
		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, path, nil))
		if err != nil {
			t.Fatalf("요청 실패: %v", err)
		}
		var out map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("응답 파싱 실패: %v", err)
		}
		return out
	}

	leaky := []string{"email", "name", "tel", "address", "birth", "level", "role"}

	t.Run("무인증이면 개인정보가 제거된다", func(t *testing.T) {
		out := read(t, newApp(false), "/api/user/search/loginid?loginid=admin")
		item, ok := out["item"].(map[string]any)
		if !ok {
			t.Fatalf("item 이 객체가 아니다: %#v", out["item"])
		}
		for _, k := range leaky {
			if _, exists := item[k]; exists {
				t.Errorf("개인정보 %q 가 남아 있다", k)
			}
		}
		// 중복 판단에 필요한 필드는 유지돼야 한다.
		for _, k := range []string{"id", "loginid"} {
			if _, exists := item[k]; !exists {
				t.Errorf("필수 필드 %q 가 사라졌다", k)
			}
		}
		if out["code"] != "ok" {
			t.Errorf("응답 형태가 바뀌었다: code=%v", out["code"])
		}
	})

	t.Run("인증된 요청은 원본 그대로", func(t *testing.T) {
		out := read(t, newApp(true), "/api/user/search/loginid?loginid=admin")
		item := out["item"].(map[string]any)
		if _, exists := item["email"]; !exists {
			t.Error("인증된 요청인데 email 이 제거됐다")
		}
	})

	t.Run("다른 경로는 건드리지 않는다", func(t *testing.T) {
		out := read(t, newApp(false), "/api/user/search/email?email=x")
		item := out["item"].(map[string]any)
		if _, exists := item["email"]; !exists {
			t.Error("대상이 아닌 경로의 응답이 변형됐다")
		}
	})
}
