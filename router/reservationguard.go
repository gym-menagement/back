package router

import (
	"encoding/json"
	"strings"

	"gym/global/log"
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// ReservationGuard 는 PT 예약 생성이 실제로 가능한 시간인지 서버에서 다시
// 판정한다.
//
// 가용 슬롯 API(/ptreservation/available)가 이미 근무시간 ∩ 휴게·기존예약을
// 계산해 주지만, 그건 화면에 보여줄 목록일 뿐이다. 생성 라우트는 buildtool
// 생성 코드라 아무 검증이 없어서, 슬롯 조회를 건너뛰고 POST 만 던지면
// 같은 시간에 몇 건이든 예약이 쌓였다 — 조회와 생성 사이의 경합에서도 마찬가지.
// (2026-08-07 프로세스 E2E 3회차에서 이중 예약이 실제로 생성됨)
//
// 생성 파일을 건드리지 않도록 미들웨어로 만들어 services/http.go 에서 건다.
// 판정 불가(본문 파싱 실패 등)면 통과시킨다 — 이 가드는 정상 흐름의 경합을
// 막는 장치지 인증·권한 경계가 아니다. 권한은 기존 JWT 미들웨어가 본다.
func ReservationGuard(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost || c.Path() != "/api/ptreservation" {
		return c.Next()
	}

	var body struct {
		Trainer         int64  `json:"trainer"`
		Reservationdate string `json:"reservationdate"`
		Starttime       string `json:"starttime"`
		Endtime         string `json:"endtime"`
		Status          int    `json:"status"`
	}
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Next()
	}
	if body.Trainer == 0 || body.Reservationdate == "" || body.Starttime == "" || body.Endtime == "" {
		return c.Next()
	}
	// 취소·완료 상태로 들어오는 기록성 생성은 막지 않는다.
	if body.Status != 0 && body.Status != PrStatusReserved {
		return c.Next()
	}

	day := body.Reservationdate
	if i := strings.IndexByte(day, 'T'); i > 0 {
		day = day[:i]
	}
	start, end := hhmmToMin(body.Starttime), hhmmToMin(body.Endtime)
	if start < 0 || end < 0 || start >= end {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad time range"})
	}

	conn := models.NewConnection()
	defer conn.Close()

	// 1) 이미 잡힌 예약과 겹치면 거부.
	rows, err := conn.Query(
		`SELECT TIME_FORMAT(pr_starttime, '%H:%i'), TIME_FORMAT(pr_endtime, '%H:%i')
		 FROM ptreservation_tb
		 WHERE pr_trainer = ? AND pr_reservationdate = ? AND pr_status = ?`,
		body.Trainer, day, PrStatusReserved)
	if err != nil {
		return c.Next() // 판정 불가 — 기존 동작 유지
	}
	for rows.Next() {
		var bs, be string
		if rows.Scan(&bs, &be) != nil {
			continue
		}
		if start < hhmmToMin(be) && hhmmToMin(bs) < end { // 반개구간 겹침
			rows.Close()
			log.Info().Int64("trainer", body.Trainer).Str("date", day).
				Str("start", body.Starttime).Msg("duplicate reservation blocked")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"code": "slot taken", "message": "이미 예약된 시간입니다.",
			})
		}
	}
	rows.Close()

	// 2) 트레이너 근무시간 밖이면 거부 — 근무 스케줄이 등록된 경우에만
	//    판정한다 (미등록 트레이너는 기존처럼 자유 예약).
	rows, err = conn.Query(
		`SELECT TIME_FORMAT(ts_start_time, '%H:%i'), TIME_FORMAT(ts_end_time, '%H:%i'),
		        IFNULL(TIME_FORMAT(ts_break_start, '%H:%i'), ''), IFNULL(TIME_FORMAT(ts_break_end, '%H:%i'), '')
		 FROM trainerschedule_tb
		 WHERE ts_trainer = ? AND ts_day_of_week = DAYOFWEEK(?) - 1`,
		body.Trainer, day)
	if err != nil {
		return c.Next()
	}
	defer rows.Close()

	hasSchedule, inShift := false, false
	for rows.Next() {
		var ws, we, bs, be string
		if rows.Scan(&ws, &we, &bs, &be) != nil {
			continue
		}
		hasSchedule = true
		if start >= hhmmToMin(ws) && end <= hhmmToMin(we) {
			// 휴게시간과 겹치면 그 근무 구간에서는 불가.
			if bs != "" && be != "" && start < hhmmToMin(be) && hhmmToMin(bs) < end {
				continue
			}
			inShift = true
			break
		}
	}
	if hasSchedule && !inShift {
		log.Info().Int64("trainer", body.Trainer).Str("date", day).
			Str("start", body.Starttime).Msg("reservation outside working hours blocked")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code": "outside working hours", "message": "트레이너 근무 시간이 아닙니다.",
		})
	}

	if err := c.Next(); err != nil {
		return err
	}

	// 위 검사는 SELECT 후 INSERT 라 그 사이에 끼어드는 동시 요청을 막지 못한다
	// (같은 슬롯 동시 5건이 전부 생성됐다). 최종 판정은 DB 유니크 제약
	// uq_pr_active_slot 이 하는데(migrations/p3b_reservation_unique.sql),
	// 생성 컨트롤러가 INSERT 에러를 삼키고 id 0 으로 200 을 돌려줘 클라이언트가
	// 예약 성공으로 오해한다. 여기서 409 로 바로잡는다.
	if insertFailed(c.Response().Body()) {
		log.Info().Int64("trainer", body.Trainer).Str("date", day).
			Str("start", body.Starttime).Msg("reservation lost race — rejected")
		c.Response().Reset()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code": "slot taken", "message": "이미 예약된 시간입니다.",
		})
	}
	return nil
}

// insertFailed 는 생성 응답이 실패를 감춘 200 인지 본다.
// 생성 컨트롤러는 INSERT 에러 시 {"code":"error", ...} 를 200 으로 돌려주고
// (id 필드 없음), 성공하면 {"code":"ok","id":N} 이다.
func insertFailed(body []byte) bool {
	var res struct {
		Code string `json:"code"`
		Id   *int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return false
	}
	if res.Code != "" && res.Code != "ok" {
		return true
	}
	return res.Id != nil && *res.Id == 0
}
