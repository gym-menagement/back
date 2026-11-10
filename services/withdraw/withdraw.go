// Package withdraw 는 회원 탈퇴(계정 삭제) 시 개인정보 정리를 담당한다.
//
// buildtool 이 생성한 UserManager.Delete 는 user_tb 행 하나만 지운다. 개인정보
// 처리방침은 "탈퇴 시 계정과 이용 기록을 삭제한다"고 고지하므로, user 를
// 참조하는 개인 기록을 한 트랜잭션으로 함께 지운다. 생성 모델은 make server 가
// 덮어쓰므로 services/billing 과 같은 방식으로 직접 SQL 을 쓴다.
//
// 스키마에 FK 가 없어 순서 제약은 없지만, user_tb 삭제를 마지막에 둔다.
package withdraw

import (
	"errors"

	"gym/models"
	"gym/services/billing"
)

var (
	ErrNotFound = errors.New("user not found")
	ErrGymOwner = errors.New("user owns a gym")
)

// membershipOf 는 uh_user/uhu_user 가 DEPRECATED 라 membership_tb 경유로만
// 연결된 행까지 잡기 위한 서브쿼리 (gym.sql 의 usehealth_tb 주석 참조).
const membershipOf = `(SELECT m_id FROM membership_tb WHERE m_user = ?)`

type step struct {
	table string
	query string
	args  func(uid int64) []interface{}
}

// ids 는 user id 를 n 번 바인딩하는 인자 생성기.
func ids(n int) func(int64) []interface{} {
	return func(uid int64) []interface{} {
		args := make([]interface{}, n)
		for i := range args {
			args[i] = uid
		}
		return args
	}
}

// purgeSteps — 삭제되는 개인 기록. 탈퇴자가 회원인 경우와 트레이너·직원인
// 경우(담당 관계·근무 스케줄·관리자 권한)를 모두 포함한다.
var purgeSteps = []step{
	{table: "pushtoken_tb", query: `DELETE FROM pushtoken_tb WHERE pt_user = ?`, args: ids(1)},
	{table: "token_tb", query: `DELETE FROM token_tb WHERE to_user = ?`, args: ids(1)},
	{table: "loginlog_tb", query: `DELETE FROM loginlog_tb WHERE ll_user = ?`, args: ids(1)},
	{table: "notificationsetting_tb", query: `DELETE FROM notificationsetting_tb WHERE ns_user = ?`, args: ids(1)},
	{table: "notificationhistory_tb", query: `DELETE FROM notificationhistory_tb WHERE nh_receiver = ?`, args: ids(1)},
	{table: "alarm_tb", query: `DELETE FROM alarm_tb WHERE al_user = ?`, args: ids(1)},
	{table: "memberbody_tb", query: `DELETE FROM memberbody_tb WHERE mb_user = ?`, args: ids(1)},
	{table: "workoutlog_tb", query: `DELETE FROM workoutlog_tb WHERE wl_user = ?`, args: ids(1)},
	{table: "usehealthusage_tb", query: `DELETE FROM usehealthusage_tb WHERE uhu_user = ? OR uhu_membership IN ` + membershipOf, args: ids(2)},
	{table: "attendance_tb", query: `DELETE FROM attendance_tb WHERE at_user = ?`, args: ids(1)},
	{table: "ptreservation_tb", query: `DELETE FROM ptreservation_tb WHERE pr_member = ? OR pr_trainer = ?`, args: ids(2)},
	{table: "inquiry_tb", query: `DELETE FROM inquiry_tb WHERE iq_user = ?`, args: ids(1)},
	{table: "trainermember_tb", query: `DELETE FROM trainermember_tb WHERE tm_member = ? OR tm_trainer = ?`, args: ids(2)},
	{table: "qrcode_tb", query: `DELETE FROM qrcode_tb WHERE qr_user = ?`, args: ids(1)},
	{table: "gymtrainer_tb", query: `DELETE FROM gymtrainer_tb WHERE gt_trainer = ?`, args: ids(1)},
	{table: "trainerschedule_tb", query: `DELETE FROM trainerschedule_tb WHERE ts_trainer = ?`, args: ids(1)},
	{table: "gymadmin_tb", query: `DELETE FROM gymadmin_tb WHERE ga_user = ?`, args: ids(1)},
	// 빌링키는 결제 "기록"이 아니라 결제 "수단"(PG 식별자·카드 뒤 4자리)이다.
	{table: "billingkey_tb", query: `DELETE FROM billingkey_tb WHERE bk_user = ?`, args: ids(1)},
}

// retainSteps — 전자상거래법상 보존 대상(계약·결제 기록)이라 행은 남기되,
// 탈퇴 후 더 동작하면 안 되는 상태만 정리한다. order/payment/paymentform/
// membership/stop/contract 는 손대지 않는다 — 이름·연락처 사본 컬럼이 없고
// user id 만 참조하므로 user_tb 행이 사라지면 식별할 수 없는 번호만 남는다.
var retainSteps = []step{
	// 탈퇴자 카드에 정기결제가 계속 청구되지 않도록 구독을 해지한다.
	{table: "subscription_tb", query: `UPDATE subscription_tb SET s_status = ? WHERE s_user = ? AND s_status <> ?`,
		args: func(uid int64) []interface{} {
			return []interface{}{billing.SubCancelled, uid, billing.SubCancelled}
		}},
	// 회원권은 남기되 종료 처리하고 출입용 QR 값은 지운다 (도어 verify 의
	// uh_qrcode 폴백으로 탈퇴 후에도 문이 열리는 것을 막는다).
	{table: "usehealth_tb", query: `UPDATE usehealth_tb SET uh_qrcode = '',
		uh_status = IF(uh_status IN (?, ?), ?, uh_status)
		WHERE uh_user = ? OR uh_membership IN ` + membershipOf,
		args: func(uid int64) []interface{} {
			return []interface{}{billing.UhStatusUse, billing.UhStatusPaused, billing.UhStatusTerminated, uid, uid}
		}},
}

// OwnsGym 은 사용자가 헬스장 소유자(gym_tb.g_user)인지 본다.
func OwnsGym(conn *models.Connection, userID int64) (bool, error) {
	rows, err := conn.Query(`SELECT 1 FROM gym_tb WHERE g_user = ? LIMIT 1`, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), rows.Err()
}

// existingTables 는 현재 DB 에 실제로 있는 테이블 집합을 돌려준다.
// notificationsetting/notificationhistory 는 모델만 있고 gym.sql 에 없으며,
// migrations/*.sql 은 환경에 따라 미적용일 수 있다 — 없는 테이블 때문에
// 탈퇴 전체가 실패하지 않도록 있는 것만 처리한다.
func existingTables(conn *models.Connection) (map[string]bool, error) {
	rows, err := conn.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables[name] = true
	}
	return tables, rows.Err()
}

func exists(conn *models.Connection, userID int64) (bool, error) {
	rows, err := conn.Query(`SELECT 1 FROM user_tb WHERE u_id = ?`, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), rows.Err()
}

// Withdraw 는 사용자의 개인 기록과 user_tb 행을 한 트랜잭션으로 지운다.
// 중간에 하나라도 실패하면 전부 롤백한다. 권한 판정은 호출자(router) 몫이다.
func Withdraw(conn *models.Connection, userID int64) error {
	if userID <= 0 {
		return ErrNotFound
	}
	ok, err := exists(conn, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	tables, err := existingTables(conn)
	if err != nil {
		return err
	}

	conn.Begin()
	if conn.Tx == nil {
		conn.Transaction = false
		return errors.New("begin transaction failed")
	}
	defer conn.Rollback() // Commit 후에는 no-op

	for _, steps := range [][]step{purgeSteps, retainSteps} {
		for _, s := range steps {
			if !tables[s.table] {
				continue
			}
			if _, err := conn.Exec(s.query, s.args(userID)...); err != nil {
				return err
			}
		}
	}

	if _, err := conn.Exec(`DELETE FROM user_tb WHERE u_id = ?`, userID); err != nil {
		return err
	}
	return conn.Commit()
}
