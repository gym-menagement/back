package router

import (
	"reflect"
	"testing"
)

func TestComputeSlots(t *testing.T) {
	// 14:00~22:00 근무, 18:00~19:00 휴게 (front/flow.md 예시)
	work := interval{14 * 60, 22 * 60}
	breaks := []interval{{18 * 60, 19 * 60}}

	t.Run("예약 없음 — 휴게만 빠진다", func(t *testing.T) {
		got := ComputeSlots(work, breaks, nil)
		want := []Slot{
			{"14:00", "15:00"}, {"15:00", "16:00"}, {"16:00", "17:00"}, {"17:00", "18:00"},
			{"19:00", "20:00"}, {"20:00", "21:00"}, {"21:00", "22:00"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("기존 예약이 슬롯을 막는다", func(t *testing.T) {
		booked := []interval{{15 * 60, 16 * 60}}
		got := ComputeSlots(work, breaks, booked)
		for _, s := range got {
			if s.Start == "15:00" {
				t.Error("booked slot must be excluded")
			}
		}
	})

	t.Run("경계 — 슬롯 끝과 예약 시작이 맞닿으면 열림 (반개구간)", func(t *testing.T) {
		booked := []interval{{16 * 60, 17 * 60}}
		got := ComputeSlots(work, breaks, booked)
		found1500, found1700 := false, false
		for _, s := range got {
			if s.Start == "15:00" {
				found1500 = true
			}
			if s.Start == "17:00" {
				found1700 = true
			}
		}
		if !found1500 || !found1700 {
			t.Errorf("adjacent slots must remain open: 15:00=%v 17:00=%v", found1500, found1700)
		}
	})

	t.Run("휴게 걸침 — 30분 걸치는 예약이 두 슬롯을 막는다", func(t *testing.T) {
		booked := []interval{{14*60 + 30, 15*60 + 30}}
		got := ComputeSlots(work, breaks, booked)
		for _, s := range got {
			if s.Start == "14:00" || s.Start == "15:00" {
				t.Errorf("overlapping slot %s must be excluded", s.Start)
			}
		}
	})

	t.Run("근무시간이 슬롯보다 짧으면 빈 목록", func(t *testing.T) {
		got := ComputeSlots(interval{14 * 60, 14*60 + 30}, nil, nil)
		if len(got) != 0 {
			t.Errorf("expected no slots, got %v", got)
		}
	})
}

func TestHHMMConversion(t *testing.T) {
	if hhmmToMin("14:30") != 14*60+30 {
		t.Error("14:30 conversion failed")
	}
	if hhmmToMin("14:30:00") != 14*60+30 {
		t.Error("HH:MM:SS must be accepted")
	}
	if hhmmToMin("bad") != -1 {
		t.Error("invalid time must return -1")
	}
	if minToHHMM(14*60+30) != "14:30" {
		t.Error("minToHHMM failed")
	}
}
