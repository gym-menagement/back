package router

import (
	"time"

	"gym/global/log"
	"gym/models"

	"github.com/gofiber/fiber/v2"
)

// P3 트레이너 근무 스케줄 (docs/roadmap-p1-p6.md).
// 스키마는 front/flow.md 의 trainer_schedule_tb 제안을 그대로 채택했다
// (migrations/p3_trainer_schedule.sql). 슬롯 계산은 서버가 한다 —
// 근무시간 ∩ (휴게·기존 예약 제외) 를 앱/웹은 그대로 그리기만 하면 된다.

type TrainerSchedule struct {
	Id         int64  `json:"id"`
	Trainer    int64  `json:"trainer"`
	Gym        int64  `json:"gym"`
	DayOfWeek  int    `json:"dayofweek"`
	StartTime  string `json:"starttime"`
	EndTime    string `json:"endtime"`
	BreakStart string `json:"breakstart"`
	BreakEnd   string `json:"breakend"`
	Date       string `json:"date"`
}

// Slot 은 예약 가능한 시간 구간 (HH:MM). SlotMinutes 단위로 자른다.
type Slot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

const SlotMinutes = 60 // PT 기본 단위. pr_duration 기본값(60분)과 일치.

// interval 은 슬롯 계산용 분 단위 반개구간 [start, end).
type interval struct{ start, end int }

// ComputeSlots 는 근무시간에서 휴게·기존 예약을 뺀 뒤 SlotMinutes 단위
// 슬롯으로 자른다. 순수 함수라 DB 없이 테스트한다 (schedule_test.go).
func ComputeSlots(work interval, breaks []interval, booked []interval) []Slot {
	blocked := append(append([]interval{}, breaks...), booked...)

	var slots []Slot
	for s := work.start; s+SlotMinutes <= work.end; s += SlotMinutes {
		e := s + SlotMinutes
		ok := true
		for _, b := range blocked {
			// 반개구간 겹침: [s,e) ∩ [b.start,b.end) ≠ ∅
			if s < b.end && b.start < e {
				ok = false
				break
			}
		}
		if ok {
			slots = append(slots, Slot{Start: minToHHMM(s), End: minToHHMM(e)})
		}
	}
	return slots
}

func minToHHMM(m int) string {
	return time.Date(2000, 1, 1, m/60, m%60, 0, 0, time.UTC).Format("15:04")
}

// hhmmToMin 은 "HH:MM" 또는 "HH:MM:SS" 를 분으로. 실패 시 -1.
func hhmmToMin(s string) int {
	if len(s) >= 5 {
		s = s[:5]
	}
	t, err := time.Parse("15:04", s)
	if err != nil {
		return -1
	}
	return t.Hour()*60 + t.Minute()
}

func SetupScheduleRoutes(group fiber.Router) {

	// 근무 스케줄 CRUD. 등록·수정은 트레이너 본인 또는 관리자.
	group.Get("/trainerschedule", func(c *fiber.Ctx) error {
		trainer := c.QueryInt("trainer")
		conn := models.NewConnection()
		defer conn.Close()
		rows, err := conn.Query(
			`SELECT ts_id, ts_trainer, ts_gym, ts_day_of_week,
			        TIME_FORMAT(ts_start_time, '%H:%i'), TIME_FORMAT(ts_end_time, '%H:%i'),
			        IFNULL(TIME_FORMAT(ts_break_start, '%H:%i'), ''), IFNULL(TIME_FORMAT(ts_break_end, '%H:%i'), ''),
			        ts_date
			 FROM trainerschedule_tb WHERE ts_trainer = ? ORDER BY ts_day_of_week, ts_start_time`, trainer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		defer rows.Close()
		items := []TrainerSchedule{}
		for rows.Next() {
			var t TrainerSchedule
			if err := rows.Scan(&t.Id, &t.Trainer, &t.Gym, &t.DayOfWeek,
				&t.StartTime, &t.EndTime, &t.BreakStart, &t.BreakEnd, &t.Date); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
			}
			items = append(items, t)
		}
		return c.JSON(fiber.Map{"code": "ok", "items": items, "total": len(items)})
	})

	group.Post("/trainerschedule", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		var t TrainerSchedule
		if err := c.BodyParser(&t); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		// 본인 스케줄이 아니면 관리자만. 자정 넘김은 지원하지 않는다(roadmap).
		if t.Trainer != me.Id && !isGymAdmin(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		s, e := hhmmToMin(t.StartTime), hhmmToMin(t.EndTime)
		if t.DayOfWeek < 0 || t.DayOfWeek > 6 || s < 0 || e < 0 || s >= e {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		bs, be := hhmmToMin(t.BreakStart), hhmmToMin(t.BreakEnd)
		hasBreak := t.BreakStart != "" && t.BreakEnd != ""
		if hasBreak && (bs < 0 || be < 0 || bs >= be) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		conn := models.NewConnection()
		defer conn.Close()
		var res interface{ LastInsertId() (int64, error) }
		var err error
		if hasBreak {
			res, err = conn.Exec(
				`INSERT INTO trainerschedule_tb
				 (ts_trainer, ts_gym, ts_day_of_week, ts_start_time, ts_end_time, ts_break_start, ts_break_end)
				 VALUES (?, ?, ?, ?, ?, ?, ?)`,
				t.Trainer, t.Gym, t.DayOfWeek, t.StartTime, t.EndTime, t.BreakStart, t.BreakEnd)
		} else {
			res, err = conn.Exec(
				`INSERT INTO trainerschedule_tb
				 (ts_trainer, ts_gym, ts_day_of_week, ts_start_time, ts_end_time)
				 VALUES (?, ?, ?, ?, ?)`,
				t.Trainer, t.Gym, t.DayOfWeek, t.StartTime, t.EndTime)
		}
		if err != nil {
			log.Error().Msg("trainerschedule insert: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		id, _ := res.LastInsertId()
		return c.JSON(fiber.Map{"code": "ok", "id": id})
	})

	group.Delete("/trainerschedule/:id<int>", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		conn := models.NewConnection()
		defer conn.Close()

		// 본인 것 또는 관리자.
		id := gymParamID(c)
		rows, err := conn.Query(`SELECT ts_trainer FROM trainerschedule_tb WHERE ts_id = ?`, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		if !rows.Next() {
			rows.Close()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
		var owner int64
		_ = rows.Scan(&owner)
		rows.Close()
		if owner != me.Id && !isGymAdmin(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
		}
		if _, err := conn.Exec(`DELETE FROM trainerschedule_tb WHERE ts_id = ?`, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})

	// 가용 슬롯: 해당 날짜의 근무시간에서 휴게·확정 예약을 뺀 목록.
	group.Get("/ptreservation/available", func(c *fiber.Ctx) error {
		trainer := int64(c.QueryInt("trainer"))
		dateStr := c.Query("date") // YYYY-MM-DD
		day, err := time.Parse("2006-01-02", dateStr)
		if trainer <= 0 || err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}

		conn := models.NewConnection()
		defer conn.Close()

		rows, err := conn.Query(
			`SELECT TIME_FORMAT(ts_start_time, '%H:%i'), TIME_FORMAT(ts_end_time, '%H:%i'),
			        IFNULL(TIME_FORMAT(ts_break_start, '%H:%i'), ''), IFNULL(TIME_FORMAT(ts_break_end, '%H:%i'), '')
			 FROM trainerschedule_tb WHERE ts_trainer = ? AND ts_day_of_week = ?`,
			trainer, int(day.Weekday()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		type sched struct {
			work   interval
			breaks []interval
		}
		var schedules []sched
		for rows.Next() {
			var ws, we, bs, be string
			if err := rows.Scan(&ws, &we, &bs, &be); err != nil {
				rows.Close()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
			}
			s := sched{work: interval{hhmmToMin(ws), hhmmToMin(we)}}
			if bs != "" && be != "" {
				s.breaks = append(s.breaks, interval{hhmmToMin(bs), hhmmToMin(be)})
			}
			schedules = append(schedules, s)
		}
		rows.Close()

		// 확정 예약(RESERVED=1 기준: model.json ptreservation status 1:reserved)
		rows, err = conn.Query(
			`SELECT TIME_FORMAT(pr_starttime, '%H:%i'), TIME_FORMAT(pr_endtime, '%H:%i')
			 FROM ptreservation_tb
			 WHERE pr_trainer = ? AND pr_reservationdate = ? AND pr_status = 1`,
			trainer, dateStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
		var booked []interval
		for rows.Next() {
			var bs, be string
			if err := rows.Scan(&bs, &be); err != nil {
				rows.Close()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
			}
			booked = append(booked, interval{hhmmToMin(bs), hhmmToMin(be)})
		}
		rows.Close()

		slots := []Slot{}
		for _, s := range schedules {
			slots = append(slots, ComputeSlots(s.work, s.breaks, booked)...)
		}
		return c.JSON(fiber.Map{"code": "ok", "items": slots, "total": len(slots)})
	})
}
