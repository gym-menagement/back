package router

import (
	"encoding/json"
	"net/http/httptest"
	"regexp"

	"gym/global/log"
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// PatchMerge 는 PATCH/PUT 부분 본문을 기존 리소스와 병합하는 미들웨어다.
//
// 배경: buildtool 생성 Update 핸들러는 본문을 제로값 구조체에 바인딩해 전
// 컬럼을 덮어쓴다. 구 Spring 백엔드는 부분 병합(로드 후 지정 필드만 세팅)
// 이었고 프론트 20여 곳이 그 계약대로 부분 본문을 보낸다 — Go 컷오버 후
// `User.update(id, {use})` 한 번에 이름·비밀번호가 비워지는 회귀가 됐다
// (2026-08-07 프로세스 E2E 루프에서 발견).
//
// 동작: PATCH/PUT /api/{entity}/{id} 요청이면 같은 앱에 내부 GET 을 날려
// 현재 리소스(item)를 얻고, 그 위에 요청 본문을 덮어쓴 병합 본문으로
// 교체한 뒤 다음 핸들러로 넘긴다. 본문에 명시된 필드만 바뀌는 원래 계약이
// 복원된다 (명시적 0/빈문자열 설정도 그대로 반영됨).
//
// - 생성 파일(router.go, routers/*)을 건드리지 않으므로 buildtool 재생성과
//   무관하게 유지된다. services/http.go 가 SetRouter 이전에 등록한다.
// - 내부 GET 은 요청의 Authorization 을 그대로 전달한다 — 조회 권한이 없는
//   호출자는 병합 없이 원본 본문으로 진행돼 기존 인증 흐름이 판정한다.
// - GET 실패·item 없음·JSON 아님 등 모든 예외는 병합 생략(원본 통과)으로
//   처리한다. 이 미들웨어는 계약 복원 장치지 검증 장치가 아니다.
var patchPathRe = regexp.MustCompile(`^/api/([a-z]+)/(\d+)$`)

func PatchMerge(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPatch && c.Method() != fiber.MethodPut {
		return c.Next()
	}
	m := patchPathRe.FindStringSubmatch(c.Path())
	if m == nil {
		return c.Next()
	}

	var body map[string]interface{}
	if err := json.Unmarshal(c.Body(), &body); err != nil || len(body) == 0 {
		return c.Next()
	}

	// 같은 앱으로 내부 GET — 라우팅·인증·컨트롤러를 그대로 태운다.
	req := httptest.NewRequest(fiber.MethodGet, c.Path(), nil)
	if auth := c.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := c.App().Test(req, 5000)
	if err != nil || resp.StatusCode != fiber.StatusOK {
		return mergeUnavailable(c, m[1])
	}
	defer resp.Body.Close()

	var envelope struct {
		Item map[string]interface{} `json:"item"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil || envelope.Item == nil {
		return mergeUnavailable(c, m[1])
	}

	// 현재 값 위에 요청 본문을 덮는다. extra 같은 클라이언트 무관 필드는
	// Update 구조체가 무시하므로 남아 있어도 무해하다.
	for k, v := range body {
		envelope.Item[k] = v
	}

	// user 의 passwd 는 특별 취급이 필요하다:
	//  - GET 이 passwd 를 마스킹하므로 병합 본문에 없음 → Update 가 빈 값을 씀
	//  - 해시를 주입해 주면 Update 가 그걸 다시 bcrypt 해싱해 로그인 불능
	// 그래서 본문에 passwd 가 없으면(=변경 의도 없음) 기존 해시를 확보해 두고
	// Update 실행 뒤 원래 해시를 되돌려 쓴다. 빈 문자열도 "변경 없음"으로
	// 본다 — 빈 비밀번호는 원래 유효한 상태가 아니다.
	restorePasswd := ""
	if m[1] == "user" {
		if pw, _ := body["passwd"].(string); pw == "" {
			delete(envelope.Item, "passwd")
			restorePasswd = currentPasswdHash(m[2])
		}
	}
	merged, err := json.Marshal(envelope.Item)
	if err != nil {
		return c.Next()
	}
	c.Request().SetBody(merged)
	c.Request().Header.SetContentLength(len(merged))

	log.Debug().Str("path", c.Path()).Int("fields", len(body)).Msg("patch merged")
	err = c.Next()

	// Update 는 전 컬럼을 덮어쓰므로 위에서 뺀 passwd 는 빈 값이 됐다 —
	// 확보해 둔 해시를 되돌린다 (본문에 passwd 를 명시한 요청은 해당 없음).
	if restorePasswd != "" {
		conn := models.NewConnection()
		defer conn.Close()
		if _, e := conn.Exec(`UPDATE user_tb SET u_passwd = ? WHERE u_id = ?`, restorePasswd, m[2]); e != nil {
			log.Error().Str("id", m[2]).Msg("passwd restore failed: " + e.Error())
		}
	}
	return err
}

// mergeUnavailable 은 현재 값을 못 읽어 병합할 수 없을 때의 처리다.
//
// 그냥 통과시키면(fail-open) 요청 본문에 없는 컬럼이 전부 제로값으로 덮여
// 행이 조용히 파괴된다 — 실제로 usehealth 에서 회원권 기간·상태·소유자가
// 한 번에 날아갔다(2026-08-07 프로세스 E2E 2회차). 생성 모델의 단건 조회는
// 연관 테이블 INNER JOIN 이라 FK 가 0 인 행이 조회되지 않는데, 그런 행일수록
// 부분 수정으로 고치려 드는 대상이라 위험이 겹친다.
//
// 그래서 막는다(fail-closed). 클라이언트는 전체 본문 PUT 으로 명시적으로
// 덮어쓰거나, 조회되지 않는 원인(누락 FK)을 먼저 고쳐야 한다.
func mergeUnavailable(c *fiber.Ctx, entity string) error {
	log.Info().Str("path", c.Path()).Str("entity", entity).Msg("patch merge unavailable — rejected")
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"code":    "merge unavailable",
		"message": "현재 값을 읽을 수 없어 부분 수정을 처리할 수 없습니다. 전체 필드를 포함해 다시 요청하세요.",
	})
}

// currentPasswdHash 는 병합용으로 사용자의 현재 비밀번호 해시를 읽는다.
// 실패 시 빈 문자열 — 호출부가 필드 제거로 폴백한다.
func currentPasswdHash(id string) string {
	conn := models.NewConnection()
	defer conn.Close()
	rows, err := conn.Query(`SELECT u_passwd FROM user_tb WHERE u_id = ?`, id)
	if err != nil {
		return ""
	}
	defer rows.Close()
	if !rows.Next() {
		return ""
	}
	var hash string
	if rows.Scan(&hash) != nil {
		return ""
	}
	return hash
}
