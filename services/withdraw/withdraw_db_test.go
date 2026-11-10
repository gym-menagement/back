package withdraw

// 실 DB 스키마 대상 기능 테스트. 일회용 MySQL/MariaDB 가 있을 때만 돈다
// (services/billing/billing_db_test.go 와 같은 방식):
//
//	WITHDRAW_TEST_DB=1 DB_HOST=127.0.0.1 DB_PORT=33071 DB_USER=root \
//	  DB_PASS=e2etest DB_NAME=gym go test ./services/withdraw/ -v  # secret-scan: ok 일회용 도커 DB 예시값
//
// gym.sql + migrations/*.sql 이 적용된 DB 여야 한다. 운영 DB 로 돌리지 말 것 —
// 테스트 데이터가 실제로 INSERT/DELETE 된다.

import (
	"os"
	"testing"

	"gym/models"
	"gym/services/billing"
)

func count(t *testing.T, conn *models.Connection, query string, args ...interface{}) int {
	t.Helper()
	rows, err := conn.Query(query, args...)
	if err != nil {
		t.Fatal(query, err)
	}
	defer rows.Close()
	var n int
	if rows.Next() {
		_ = rows.Scan(&n)
	}
	return n
}

func TestWithdraw(t *testing.T) {
	if os.Getenv("WITHDRAW_TEST_DB") == "" {
		t.Skip("WITHDRAW_TEST_DB not set — DB 기능 테스트 건너뜀")
	}
	conn := models.NewConnection()
	defer conn.Close()

	mustExec := func(q string, args ...interface{}) int64 {
		t.Helper()
		res, err := conn.Exec(q, args...)
		if err != nil {
			t.Fatal(q, err)
		}
		id, _ := res.LastInsertId()
		return id
	}

	uid := mustExec(`INSERT INTO user_tb (u_loginid, u_name, u_tel) VALUES ('withdraw-me', '탈퇴자', '010-0000-0000')`)
	other := mustExec(`INSERT INTO user_tb (u_loginid, u_name) VALUES ('withdraw-other', '남는사람')`)

	for _, u := range []int64{uid, other} {
		mustExec(`INSERT INTO pushtoken_tb (pt_user, pt_token) VALUES (?, 'tok')`, u)
		mustExec(`INSERT INTO memberbody_tb (mb_user, mb_note) VALUES (?, '')`, u)
		mustExec(`INSERT INTO attendance_tb (at_user, at_note) VALUES (?, '')`, u)
		mustExec(`INSERT INTO workoutlog_tb (wl_user, wl_note) VALUES (?, '')`, u)
		mustExec(`INSERT INTO inquiry_tb (iq_user, iq_content, iq_answer) VALUES (?, '', '')`, u)
		mustExec(`INSERT INTO loginlog_tb (ll_user) VALUES (?)`, u)
		mustExec(`INSERT INTO order_tb (o_user) VALUES (?)`, u)
		mustExec(`INSERT INTO payment_tb (p_user, p_cost) VALUES (?, 1000)`, u)
		mustExec(`INSERT INTO contract_tb (ct_user, ct_provider, ct_document_id) VALUES (?, 'test', ?)`, u, u)
		mid := mustExec(`INSERT INTO membership_tb (m_user) VALUES (?)`, u)
		// uh_user 없이 membership 으로만 연결된 회원권/이용내역 (DEPRECATED 컬럼 대비)
		mustExec(`INSERT INTO usehealth_tb (uh_membership, uh_status, uh_qrcode) VALUES (?, ?, 'QR')`, mid, billing.UhStatusUse)
		mustExec(`INSERT INTO usehealthusage_tb (uhu_membership, uhu_note) VALUES (?, '')`, mid)
		bk := mustExec(`INSERT INTO billingkey_tb (bk_user, bk_customer_uid) VALUES (?, ?)`, u, u)
		mustExec(`INSERT INTO subscription_tb (s_user, s_billingkey, s_status, s_next_billing_date) VALUES (?, ?, ?, CURDATE())`,
			u, bk, billing.SubActive)
	}
	// 당사자 양쪽 관계: 남는 사람이 회원, 탈퇴자가 트레이너
	mustExec(`INSERT INTO ptreservation_tb (pr_member, pr_trainer, pr_note, pr_cancelreason) VALUES (?, ?, '', '')`, other, uid)
	mustExec(`INSERT INTO trainermember_tb (tm_member, tm_trainer, tm_note) VALUES (?, ?, '')`, other, uid)

	if err := Withdraw(conn, uid); err != nil {
		t.Fatal("Withdraw:", err)
	}

	deleted := map[string]string{
		"user_tb": "u_id", "pushtoken_tb": "pt_user", "memberbody_tb": "mb_user",
		"attendance_tb": "at_user", "workoutlog_tb": "wl_user", "inquiry_tb": "iq_user",
		"loginlog_tb": "ll_user", "billingkey_tb": "bk_user",
		"ptreservation_tb": "pr_trainer", "trainermember_tb": "tm_trainer",
	}
	for table, col := range deleted {
		if n := count(t, conn, `SELECT COUNT(*) FROM `+table+` WHERE `+col+` = ?`, uid); n != 0 {
			t.Errorf("%s: 탈퇴자 행 %d개 남음", table, n)
		}
	}
	retained := map[string]string{
		"order_tb": "o_user", "payment_tb": "p_user", "membership_tb": "m_user",
		"contract_tb": "ct_user", "subscription_tb": "s_user",
	}
	for table, col := range retained {
		if n := count(t, conn, `SELECT COUNT(*) FROM `+table+` WHERE `+col+` = ?`, uid); n != 1 {
			t.Errorf("%s: 보존 행 = %d, want 1", table, n)
		}
	}
	if n := count(t, conn, `SELECT COUNT(*) FROM subscription_tb WHERE s_user = ? AND s_status = ?`, uid, billing.SubCancelled); n != 1 {
		t.Error("구독이 해지되지 않음")
	}
	const uhOf = `SELECT COUNT(*) FROM usehealth_tb WHERE uh_membership IN (SELECT m_id FROM membership_tb WHERE m_user = ?)`
	if n := count(t, conn, uhOf+` AND uh_status = ? AND uh_qrcode = ''`, uid, billing.UhStatusTerminated); n != 1 {
		t.Error("회원권이 종료·QR 삭제 처리되지 않음")
	}
	if n := count(t, conn, `SELECT COUNT(*) FROM usehealthusage_tb WHERE uhu_membership IN (SELECT m_id FROM membership_tb WHERE m_user = ?)`, uid); n != 0 {
		t.Error("membership 경유 이용내역이 남음")
	}

	// 다른 사용자의 데이터는 그대로여야 한다.
	for table, col := range map[string]string{
		"user_tb": "u_id", "pushtoken_tb": "pt_user", "memberbody_tb": "mb_user", "attendance_tb": "at_user",
		"workoutlog_tb": "wl_user", "inquiry_tb": "iq_user", "loginlog_tb": "ll_user", "billingkey_tb": "bk_user",
	} {
		if n := count(t, conn, `SELECT COUNT(*) FROM `+table+` WHERE `+col+` = ?`, other); n != 1 {
			t.Errorf("%s: 다른 사용자 행 = %d, want 1", table, n)
		}
	}
	if n := count(t, conn, `SELECT COUNT(*) FROM subscription_tb WHERE s_user = ? AND s_status = ?`, other, billing.SubActive); n != 1 {
		t.Error("다른 사용자의 구독이 바뀜")
	}
	if n := count(t, conn, uhOf+` AND uh_status = ? AND uh_qrcode = 'QR'`, other, billing.UhStatusUse); n != 1 {
		t.Error("다른 사용자의 회원권이 바뀜")
	}

	if err := Withdraw(conn, uid); err != ErrNotFound {
		t.Errorf("재탈퇴 = %v, want ErrNotFound", err)
	}
}
