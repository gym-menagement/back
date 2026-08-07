package router

import (
	"database/sql"
	"strconv"
	"strings"

	"gym/global/log"
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// DualUserFix 는 user_tb 를 두 번 참조하는 엔티티의 조회를 대체한다.
//
// buildtool 생성 모델은 자기참조 FK 가 둘인 테이블에서 `from x_tb, user_tb,
// user_tb, gym_tb` 처럼 **별칭 없이** user_tb 를 두 번 나열하고
// `x_trainer = u_id and x_member = u_id` 조건을 건다. MySQL 은 이걸
// "Not unique table/alias: 'user_tb'" 로 거부하고, 컨트롤러는 에러를 삼켜
// 200 + 빈 배열을 돌려준다 — 즉 **조회가 조용히 전부 실패**한다.
// (2026-08-07 프로세스 E2E 3회차에서 PT 예약 목록이 항상 비어 발견)
//
// 영향: ptreservation, trainermember, inquiry, memberbody — 관리자 웹의
// PT 예약 관리·트레이너 배정·문의 관리 화면과 앱의 예약/체성분 조회가
// 전부 빈 화면이었다.
//
// 생성 파일은 make server 에서 덮어써지므로 고칠 수 없고, 생성기(buildtool)는
// 로컬에 없다. 그래서 라우트보다 먼저 도는 미들웨어에서 목록·단건 GET 을
// 가로채 별칭을 붙인 SQL 로 직접 응답한다. 응답 형태는 생성 컨트롤러와
// 동일하게 유지한다({code, items, total} / {code, item}) — 프론트 인터셉터
// 계약이 그대로 성립해야 한다.
//
// buildtool 이 고쳐지면 이 파일을 통째로 삭제하고 http.go 의 등록만 빼면 된다.

type dualUserEntity struct {
	table   string // ptreservation_tb
	prefix  string // pr
	columns []string
}

var dualUserEntities = map[string]dualUserEntity{
	"ptreservation": {"ptreservation_tb", "pr", []string{
		"id", "trainer", "member", "gym", "reservationdate", "starttime", "endtime",
		"duration", "status", "note", "cancelreason", "createddate", "updateddate", "date"}},
	"trainermember": {"trainermember_tb", "tm", []string{
		"id", "trainer", "member", "gym", "startdate", "enddate", "status", "note", "date"}},
	"inquiry": {"inquiry_tb", "iq", []string{
		"id", "user", "gym", "type", "title", "content", "status", "answer",
		"answeredby", "answereddate", "createddate", "date"}},
	"memberbody": {"memberbody_tb", "mb", []string{
		"id", "gym", "user", "height", "weight", "bodyfat", "musclemass", "bmi",
		"skeletalmuscle", "bodywater", "chest", "waist", "hip", "arm", "thigh",
		"note", "measureddate", "measuredby", "date"}},
}

func DualUserFix(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Next()
	}
	parts := strings.Split(strings.TrimPrefix(c.Path(), "/api/"), "/")
	ent, ok := dualUserEntities[parts[0]]
	if !ok || len(parts) > 2 {
		return c.Next()
	}

	conn := models.NewConnection()
	defer conn.Close()

	// 실제 컬럼만 고른다 — 엔티티별 스키마 차이(오타·미존재 컬럼)로
	// 전체 조회가 깨지지 않게 한다.
	cols := existingColumns(conn, ent)
	if len(cols) == 0 {
		return c.Next() // 판정 불가 — 기존 동작 유지
	}

	sel := make([]string, len(cols))
	for i, col := range cols {
		sel[i] = ent.prefix + "_" + col
	}
	base := "SELECT " + strings.Join(sel, ", ") + " FROM " + ent.table

	// 단건: /api/{entity}/{id}
	if len(parts) == 2 {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return c.Next() // /count, /search/... 등은 원래 라우트로
		}
		rows, err := conn.Query(base+" WHERE "+ent.prefix+"_id = ?", id)
		if err != nil {
			log.Error().Msg("dualuserfix single: " + err.Error())
			return c.Next()
		}
		defer rows.Close()
		if !rows.Next() {
			return c.JSON(fiber.Map{"code": "ok"})
		}
		item, err := scanRow(rows, cols)
		if err != nil {
			return c.Next()
		}
		return c.JSON(fiber.Map{"code": "ok", "item": item})
	}

	// 목록: /api/{entity}?필드=값&page=&pagesize=
	where, args := []string{}, []interface{}{}
	for _, col := range cols {
		if col == "id" {
			continue
		}
		if v := c.Query(col); v != "" {
			where = append(where, ent.prefix+"_"+col+" = ?")
			args = append(args, v)
		}
	}
	q := base
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}

	total := countRows(conn, ent, where, args)

	q += " ORDER BY " + ent.prefix + "_id DESC"
	pagesize := c.QueryInt("pagesize")
	if pagesize <= 0 {
		pagesize = 20
	}
	q += " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa(c.QueryInt("page")*pagesize)

	rows, err := conn.Query(q, args...)
	if err != nil {
		log.Error().Msg("dualuserfix list: " + err.Error())
		return c.Next()
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	for rows.Next() {
		item, err := scanRow(rows, cols)
		if err != nil {
			return c.Next()
		}
		items = append(items, item)
	}
	return c.JSON(fiber.Map{"code": "ok", "items": items, "total": total})
}

func existingColumns(conn *models.Connection, ent dualUserEntity) []string {
	rows, err := conn.Query(
		`SELECT column_name FROM information_schema.columns
		 WHERE table_schema = DATABASE() AND table_name = ?`, ent.table)
	if err != nil {
		return nil
	}
	defer rows.Close()
	have := map[string]bool{}
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			have[strings.TrimPrefix(name, ent.prefix+"_")] = true
		}
	}
	var cols []string
	for _, col := range ent.columns {
		if have[col] {
			cols = append(cols, col)
		}
	}
	return cols
}

func countRows(conn *models.Connection, ent dualUserEntity, where []string, args []interface{}) int64 {
	q := "SELECT COUNT(*) FROM " + ent.table
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	rows, err := conn.Query(q, args...)
	if err != nil {
		return 0
	}
	defer rows.Close()
	var n int64
	if rows.Next() {
		_ = rows.Scan(&n)
	}
	return n
}

// scanRow 는 컬럼값을 JSON 친화 타입으로 변환한다. 숫자 컬럼이 문자열로
// 나가면 프론트 비교(=== 1)가 깨지므로 정수는 정수로 유지한다.
func scanRow(rows *sql.Rows, cols []string) (map[string]interface{}, error) {
	raw := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range raw {
		ptrs[i] = &raw[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return nil, err
	}
	item := make(map[string]interface{}, len(cols))
	for i, col := range cols {
		switch v := raw[i].(type) {
		case []byte:
			s := string(v)
			if n, err := strconv.ParseInt(s, 10, 64); err == nil {
				item[col] = n
			} else {
				item[col] = s
			}
		case nil:
			item[col] = ""
		default:
			item[col] = v
		}
	}
	return item, nil
}
