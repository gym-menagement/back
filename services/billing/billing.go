// Package billing 은 정기 구독결제(P1)를 담당한다.
//
// buildtool 생성 모델을 쓰지 않고 직접 SQL 을 쓴다 — billingkey/subscription 은
// model.json 에 등록돼 있어 나중에 buildtool 로 정식 모델을 뽑을 수 있지만,
// 청구 로직은 어차피 손코드라 여기 모아 둔다. 스키마는 migrations/p1_billing.sql.
package billing

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gym/global/log"
	"gym/models"
)

// 상태 상수 — DB int 값과 일치 (migrations/p1_billing.sql 주석 참조)
const (
	BkActive  = 1
	BkDeleted = 2

	SubActive    = 1
	SubPaused    = 2
	SubCancelled = 3
	SubFailed    = 4

	PaySuccess   = 1
	PayPending   = 2
	PayFail      = 3
	PayCancelled = 4
	PayRefunded  = 5

	MethodCard     = 1
	MethodTransfer = 2
	MethodCash     = 3
	MethodBilling  = 4

	// 연속 실패 상한. 넘으면 구독을 failed 로 전환하고 더 시도하지 않는다.
	MaxFailCount = 3
	// 결제 실패 후 회원권을 유지해 주는 유예일. 무인 센터 통상 관행.
	GraceDays = 3
)

type BillingKey struct {
	Id          int64  `json:"id"`
	User        int64  `json:"user"`
	Pg          string `json:"pg"`
	CustomerUid string `json:"customeruid"`
	CardName    string `json:"cardname"`
	CardLast4   string `json:"cardlast4"`
	Status      int    `json:"status"`
	Date        string `json:"date"`
}

type Subscription struct {
	Id              int64  `json:"id"`
	User            int64  `json:"user"`
	Gym             int64  `json:"gym"`
	Health          int64  `json:"health"`
	Usehealth       int64  `json:"usehealth"`
	Billingkey      int64  `json:"billingkey"`
	Cost            int    `json:"cost"`
	Status          int    `json:"status"`
	NextBillingDate string `json:"nextbillingdate"`
	FailCount       int    `json:"failcount"`
	Date            string `json:"date"`
}

// PG 는 실 결제 게이트웨이 추상화. 운영 구현은 pg_portone.go, 테스트는 mock.
type PG interface {
	// Charge 는 빌링키로 amount 를 청구하고 PG 거래 ID 를 돌려준다.
	Charge(customerUid string, merchantUid string, amount int, name string) (tid string, err error)
}

// ---- 빌링키 ----

func InsertBillingKey(conn *models.Connection, k *BillingKey) (int64, error) {
	res, err := conn.Exec(
		`INSERT INTO billingkey_tb (bk_user, bk_pg, bk_customer_uid, bk_card_name, bk_card_last4, bk_status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		k.User, k.Pg, k.CustomerUid, k.CardName, k.CardLast4, BkActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetBillingKey(conn *models.Connection, id int64) (*BillingKey, error) {
	rows, err := conn.Query(
		`SELECT bk_id, bk_user, bk_pg, bk_customer_uid, bk_card_name, bk_card_last4, bk_status, bk_date
		 FROM billingkey_tb WHERE bk_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	var k BillingKey
	if err := rows.Scan(&k.Id, &k.User, &k.Pg, &k.CustomerUid, &k.CardName, &k.CardLast4, &k.Status, &k.Date); err != nil {
		return nil, err
	}
	return &k, nil
}

// ---- 구독 ----

func InsertSubscription(conn *models.Connection, s *Subscription) (int64, error) {
	res, err := conn.Exec(
		`INSERT INTO subscription_tb
		 (s_user, s_gym, s_health, s_usehealth, s_billingkey, s_cost, s_status, s_next_billing_date)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.User, s.Gym, s.Health, s.Usehealth, s.Billingkey, s.Cost, SubActive, s.NextBillingDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetSubscription(conn *models.Connection, id int64) (*Subscription, error) {
	rows, err := conn.Query(subscriptionSelect+` WHERE s_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	return scanSubscription(rows)
}

const subscriptionSelect = `SELECT s_id, s_user, s_gym, s_health, s_usehealth, s_billingkey, s_cost,
	s_status, DATE_FORMAT(s_next_billing_date, '%Y-%m-%d'), s_fail_count, s_date FROM subscription_tb`

func scanSubscription(rows *sql.Rows) (*Subscription, error) {
	var s Subscription
	if err := rows.Scan(&s.Id, &s.User, &s.Gym, &s.Health, &s.Usehealth, &s.Billingkey, &s.Cost,
		&s.Status, &s.NextBillingDate, &s.FailCount, &s.Date); err != nil {
		return nil, err
	}
	return &s, nil
}

func FindSubscriptions(conn *models.Connection, where string, params ...interface{}) ([]*Subscription, error) {
	rows, err := conn.Query(subscriptionSelect+" "+where, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Subscription
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func CancelSubscription(conn *models.Connection, id int64) error {
	_, err := conn.Exec(`UPDATE subscription_tb SET s_status = ? WHERE s_id = ?`, SubCancelled, id)
	return err
}

// ---- 청구 ----

// NextBillingDate 는 이번 결제 성공 후의 다음 결제일을 계산한다.
// 단순 1개월 뒤가 아니라 "기준일 + 1개월" — 월말 넘침은 Go 의 AddDate 규칙을
// 그대로 따르되, 시각이 아닌 날짜만 쓴다.
func NextBillingDate(current time.Time) string {
	return current.AddDate(0, 1, 0).Format("2006-01-02")
}

// merchantUid 는 PG 로 보내는 주문 식별자. 구독 ID + 결제일이라 같은 날 같은
// 구독을 두 번 청구하려 하면 PG 단에서도 중복이 드러난다.
func merchantUid(subID int64, day string) string {
	return fmt.Sprintf("sub_%d_%s", subID, day)
}

// alreadyCharged 는 멱등성 가드 — 오늘 이 구독의 성공 결제가 이미 있으면
// cron 이 겹쳐 돌아도 이중 청구하지 않는다.
func alreadyCharged(conn *models.Connection, subID int64, day string) (bool, error) {
	rows, err := conn.Query(
		`SELECT p_id FROM payment_tb
		 WHERE p_subscription = ? AND p_status = ? AND DATE(p_date) = ? LIMIT 1`,
		subID, PaySuccess, day)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// ChargeSubscription 은 구독 1건을 청구한다. 성공 시 payment 기록 + 회원권
// (usehealth) 연장 + 다음 결제일 갱신, 실패 시 실패 기록 + fail_count 증가.
func ChargeSubscription(conn *models.Connection, pg PG, s *Subscription, today time.Time) error {
	day := today.Format("2006-01-02")

	if done, err := alreadyCharged(conn, s.Id, day); err != nil {
		return err
	} else if done {
		return nil
	}

	key, err := GetBillingKey(conn, s.Billingkey)
	if err != nil {
		return fmt.Errorf("billingkey %d: %w", s.Billingkey, err)
	}
	if key.Status != BkActive {
		return markFail(conn, s, "billing key deleted")
	}

	tid, err := pg.Charge(key.CustomerUid, merchantUid(s.Id, day), s.Cost, fmt.Sprintf("구독 갱신 #%d", s.Id))
	if err != nil {
		return markFail(conn, s, err.Error())
	}

	// 성공 경로는 한 트랜잭션으로 묶는다 — 결제 기록과 회원권 연장이
	// 따로 놀면 "돈은 나갔는데 회원권은 그대로" 가 된다.
	conn.Begin()
	_, err = conn.Exec(
		`INSERT INTO payment_tb (p_gym, p_order, p_user, p_cost, p_status, p_method, p_pg_tid, p_subscription)
		 VALUES (?, 0, ?, ?, ?, ?, ?, ?)`,
		s.Gym, s.User, s.Cost, PaySuccess, MethodBilling, tid, s.Id)
	if err == nil {
		// 만료 전 갱신이면 endday 에서, 만료 후(유예)면 오늘부터 1개월.
		_, err = conn.Exec(
			`UPDATE usehealth_tb
			 SET uh_endday = DATE_ADD(GREATEST(uh_endday, NOW()), INTERVAL 1 MONTH), uh_status = 1
			 WHERE uh_id = ?`, s.Usehealth)
	}
	if err == nil {
		_, err = conn.Exec(
			`UPDATE subscription_tb SET s_next_billing_date = ?, s_fail_count = 0, s_status = ? WHERE s_id = ?`,
			NextBillingDate(today), SubActive, s.Id)
	}
	if err != nil {
		conn.Rollback()
		// 결제는 이미 나갔는데 기록에 실패한 케이스 — 로그를 남겨 수동 대사 대상.
		log.Error().Int64("sub", s.Id).Str("tid", tid).Msg("charged but record failed: " + err.Error())
		return err
	}
	return conn.Commit()
}

func markFail(conn *models.Connection, s *Subscription, reason string) error {
	if len(reason) > 250 {
		reason = reason[:250]
	}
	_, err := conn.Exec(
		`INSERT INTO payment_tb (p_gym, p_order, p_user, p_cost, p_status, p_method, p_subscription, p_fail_reason)
		 VALUES (?, 0, ?, ?, ?, ?, ?, ?)`,
		s.Gym, s.User, s.Cost, PayFail, MethodBilling, s.Id, reason)
	if err != nil {
		return err
	}

	s.FailCount++
	status := SubActive
	if s.FailCount >= MaxFailCount {
		status = SubFailed
	}
	_, err = conn.Exec(
		`UPDATE subscription_tb SET s_fail_count = ?, s_status = ? WHERE s_id = ?`,
		s.FailCount, status, s.Id)
	if err != nil {
		return err
	}
	return errors.New("charge failed: " + reason)
}

// RunDailyBilling 은 cron 진입점 — 오늘이 결제일인 active 구독을 전부 청구한다.
// 개별 실패는 로그만 남기고 다음 건으로 진행한다 (한 건 실패가 전체를 막으면 안 됨).
func RunDailyBilling(pg PG) {
	conn := models.NewConnection()
	defer conn.Close()

	today := time.Now()
	subs, err := FindSubscriptions(conn,
		`WHERE s_status = ? AND s_next_billing_date <= ?`, SubActive, today.Format("2006-01-02"))
	if err != nil {
		log.Error().Msg("daily billing query failed: " + err.Error())
		return
	}

	log.Info().Int("count", len(subs)).Msg("daily billing start")
	for _, s := range subs {
		if err := ChargeSubscription(conn, pg, s, today); err != nil {
			log.Error().Int64("sub", s.Id).Msg("billing failed: " + err.Error())
		}
	}
}
