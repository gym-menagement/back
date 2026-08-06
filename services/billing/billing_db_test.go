package billing

// 실 DB 스키마 대상 기능 테스트. 일회용 MySQL 이 있을 때만 돈다:
//
//	BILLING_TEST_DB=1 DB_HOST=127.0.0.1 DB_PORT=33061 DB_USER=root \
//	  DB_PASS=e2etest DB_NAME=gym go test ./services/billing/ -run TestChargeFlow -v  # secret-scan: ok 일회용 도커 MySQL 예시값
//
// migrations/p1_billing.sql 이 적용된 DB 여야 한다. 운영 DB 로 돌리지 말 것 —
// 테스트 데이터가 실제로 INSERT 된다.

import (
	"os"
	"testing"
	"time"

	"gym/models"
)

type fakePG struct {
	tid  string
	fail bool
}

func (f *fakePG) Charge(customerUid, merchantUid string, amount int, name string) (string, error) {
	if f.fail {
		return "", errFake
	}
	if f.tid == "" {
		return merchantUid, nil
	}
	return f.tid, nil
}

var errFake = &fakeErr{}

type fakeErr struct{}

func (e *fakeErr) Error() string { return "card declined (test)" }

func TestChargeFlow(t *testing.T) {
	if os.Getenv("BILLING_TEST_DB") == "" {
		t.Skip("BILLING_TEST_DB not set — DB 기능 테스트 건너뜀")
	}

	conn := models.NewConnection()
	defer conn.Close()
	if !conn.IsConnect() {
		t.Fatal("DB 연결 실패 — 일회용 MySQL 기동 여부 확인")
	}

	// 픽스처: 회원권(usehealth) + 빌링키 + 구독
	res, err := conn.Exec(
		`INSERT INTO usehealth_tb (uh_user, uh_gym, uh_status, uh_startday, uh_endday, uh_totalcount)
		 VALUES (9001, 9001, 2, NOW(), DATE_ADD(NOW(), INTERVAL 3 DAY), 0)`)
	if err != nil {
		t.Fatal("usehealth fixture:", err)
	}
	uhID, _ := res.LastInsertId()

	bkID, err := InsertBillingKey(conn, &BillingKey{
		User: 9001, Pg: "portone", CustomerUid: "test-cuid-9001", CardName: "테스트", CardLast4: "1234",
	})
	if err != nil {
		t.Fatal("billingkey fixture:", err)
	}

	today := time.Now()
	subID, err := InsertSubscription(conn, &Subscription{
		User: 9001, Gym: 9001, Health: 1, Usehealth: uhID, Billingkey: bkID,
		Cost: 50000, NextBillingDate: today.Format("2006-01-02"),
	})
	if err != nil {
		t.Fatal("subscription fixture:", err)
	}

	sub, err := GetSubscription(conn, subID)
	if err != nil {
		t.Fatal("get subscription:", err)
	}

	// 1) 청구 성공 → payment 기록 + endday 연장 + next_billing_date 갱신
	if err := ChargeSubscription(conn, &fakePG{}, sub, today); err != nil {
		t.Fatal("charge:", err)
	}

	var payCount int
	rows, _ := conn.Query(`SELECT COUNT(*) FROM payment_tb WHERE p_subscription = ? AND p_status = 1`, subID)
	if rows.Next() {
		_ = rows.Scan(&payCount)
	}
	rows.Close()
	if payCount != 1 {
		t.Fatalf("성공 결제 1건이어야 함, got %d", payCount)
	}

	var extended int
	rows, _ = conn.Query(
		`SELECT uh_endday > DATE_ADD(NOW(), INTERVAL 25 DAY) FROM usehealth_tb WHERE uh_id = ?`, uhID)
	if rows.Next() {
		_ = rows.Scan(&extended)
	}
	rows.Close()
	if extended != 1 {
		t.Error("회원권이 1개월 연장돼야 함")
	}

	after, _ := GetSubscription(conn, subID)
	if after.NextBillingDate != NextBillingDate(today) {
		t.Errorf("다음 결제일 %s, want %s", after.NextBillingDate, NextBillingDate(today))
	}

	// 2) 멱등성: 같은 날 다시 청구해도 결제가 늘지 않는다
	if err := ChargeSubscription(conn, &fakePG{}, sub, today); err != nil {
		t.Fatal("idempotent charge:", err)
	}
	rows, _ = conn.Query(`SELECT COUNT(*) FROM payment_tb WHERE p_subscription = ? AND p_status = 1`, subID)
	if rows.Next() {
		_ = rows.Scan(&payCount)
	}
	rows.Close()
	if payCount != 1 {
		t.Errorf("멱등성 위반: 성공 결제 %d건 (want 1)", payCount)
	}

	// 3) 실패 경로: 실패 기록 + fail_count 증가, MaxFailCount 도달 시 failed
	failSub, err := GetSubscription(conn, subID)
	if err != nil {
		t.Fatal(err)
	}
	tomorrow := today.AddDate(0, 1, 0) // 다음 결제일로 이동해 멱등 가드 회피
	for i := 1; i <= MaxFailCount; i++ {
		_ = ChargeSubscription(conn, &fakePG{fail: true}, failSub, tomorrow)
	}
	final, _ := GetSubscription(conn, subID)
	if final.Status != SubFailed {
		t.Errorf("연속 %d회 실패 후 status=failed 여야 함, got %d", MaxFailCount, final.Status)
	}
	if final.FailCount < MaxFailCount {
		t.Errorf("fail_count %d, want >= %d", final.FailCount, MaxFailCount)
	}

	// 정리: 테스트 픽스처 제거
	_, _ = conn.Exec(`DELETE FROM payment_tb WHERE p_subscription = ?`, subID)
	_, _ = conn.Exec(`DELETE FROM subscription_tb WHERE s_id = ?`, subID)
	_, _ = conn.Exec(`DELETE FROM billingkey_tb WHERE bk_id = ?`, bkID)
	_, _ = conn.Exec(`DELETE FROM usehealth_tb WHERE uh_id = ?`, uhID)
}
