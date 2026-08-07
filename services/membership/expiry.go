// Package membership 은 회원권 수명주기 배치를 담당한다.
//
// 스키마상 회원권(usehealth_tb)은 상태(uh_status)와 종료일(uh_endday)을 따로
// 갖는데, 종료일이 지나도 상태를 바꿔주는 주체가 없었다. 출입(P2)·QR 발급은
// 종료일을 함께 보므로 실제 입장은 정확히 막히지만, 상태값만 읽는 화면·집계는
// 만료 회원을 계속 '사용중'으로 표시했다 (2026-08-07 4회차 점검에서 확인,
// 운영 DB 에도 해당 행 존재).
package membership

import (
	"fmt"

	"gym/global/log"
	"gym/models"
)

const (
	// usehealth_tb.uh_status (buildtool enum, 1부터 시작)
	StatusTerminated = 1
	StatusUse        = 2
	StatusPaused     = 3
	StatusExpired    = 4

	// alarm_tb.al_type / al_status
	alarmTypeNotice  = 1
	alarmStatusReady = 3 // pending — 발송 주체가 집어가기 전 상태
)

// ExpireOverdue 는 종료일이 지난 '사용중' 회원권을 '만료'로 바꾼다.
//
// 일시정지(3)는 건드리지 않는다 — 정지 기간만큼 종료일을 미뤄주는 운영이
// 일반적이라, 정지 중 종료일이 지났다고 만료시키면 정지 해제 시 되돌릴 수
// 없다. 종료(1)·이미 만료(4)도 대상이 아니다.
func ExpireOverdue() {
	conn := models.NewConnection()
	defer conn.Close()

	res, err := conn.Exec(
		`UPDATE usehealth_tb SET uh_status = ?
		 WHERE uh_status = ? AND uh_endday < NOW()`,
		StatusExpired, StatusUse)
	if err != nil {
		log.Error().Msg("expire batch failed: " + err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n > 0 {
		log.Info().Int64("count", n).Msg("memberships expired")
	}
}

// NotifyExpiring 은 만료 D-7/D-3/D-1 회원에게 알림을 남긴다.
//
// alarm_tb 에 대기(pending) 상태로 적재만 하고 실제 발송은 하지 않는다 —
// 푸시 발송은 FCM 설정이 있는 환경에서 별도 워커가 집어가는 구조이고,
// 설정이 없는 환경에서 배치가 실패로 끝나면 만료 처리까지 못 돌기 때문이다.
//
// 하루 한 번 도는 것을 전제로 같은 날 중복 적재를 막는다 — 배치가 두 번
// 돌아도 회원이 같은 알림을 두 번 받지 않는다.
func NotifyExpiring() {
	conn := models.NewConnection()
	defer conn.Close()

	for _, days := range []int{7, 3, 1} {
		title := "회원권 만료 안내"
		content := fmt.Sprintf("회원권이 %d일 후 만료됩니다. 갱신을 원하시면 센터에 문의해주세요.", days)

		// DATE(uh_endday) = 오늘 + N일 인 사용중 회원권의 소유자.
		// 같은 날 같은 제목의 알림이 이미 있으면 건너뛴다(멱등).
		res, err := conn.Exec(
			`INSERT INTO alarm_tb (al_title, al_content, al_type, al_status, al_user, al_date)
			 SELECT ?, ?, ?, ?, uh_user, NOW()
			 FROM usehealth_tb u
			 WHERE u.uh_status = ?
			   AND DATE(u.uh_endday) = DATE_ADD(CURDATE(), INTERVAL ? DAY)
			   AND NOT EXISTS (
			     SELECT 1 FROM alarm_tb a
			     WHERE a.al_user = u.uh_user AND a.al_content = ? AND DATE(a.al_date) = CURDATE()
			   )
			 GROUP BY uh_user`,
			title, content, alarmTypeNotice, alarmStatusReady, StatusUse, days, content)
		if err != nil {
			log.Error().Int("days", days).Msg("expiry notice failed: " + err.Error())
			continue
		}
		if n, _ := res.RowsAffected(); n > 0 {
			log.Info().Int("days", days).Int64("count", n).Msg("expiry notices queued")
		}
	}
}

// RunDaily 는 cron 진입점 — 알림을 먼저 보내고 만료 처리한다.
// 순서가 중요하다: 만료를 먼저 하면 D-1 대상이 이미 만료로 바뀌어 사라진다.
func RunDaily() {
	NotifyExpiring()
	ExpireOverdue()
}
