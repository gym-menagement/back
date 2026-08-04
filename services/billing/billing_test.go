package billing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"
	"time"
)

func TestNextBillingDate(t *testing.T) {
	tests := []struct {
		name string
		now  string
		want string
	}{
		{"평일", "2026-07-15", "2026-08-15"},
		{"월말 31→31", "2026-07-31", "2026-08-31"},
		// Go AddDate 규칙: 1/31 +1개월 = 3/2(평년) — 넘침이 다음 달로 밀린다.
		// PG 정기결제 관행과 다를 수 있으나 규칙이 결정적(deterministic)인 게 중요.
		{"1월 31일 넘침", "2026-01-31", "2026-03-03"},
		{"연말", "2026-12-05", "2027-01-05"},
	}
	for _, tt := range tests {
		now, _ := time.Parse("2006-01-02", tt.now)
		if got := NextBillingDate(now); got != tt.want {
			t.Errorf("%s: NextBillingDate(%s) = %s, want %s", tt.name, tt.now, got, tt.want)
		}
	}
}

func TestMerchantUid(t *testing.T) {
	// 같은 구독+같은 날 = 같은 uid (PG 단 중복 감지의 전제)
	a := merchantUid(7, "2026-07-28")
	b := merchantUid(7, "2026-07-28")
	if a != b {
		t.Errorf("same sub+day must produce same uid: %s != %s", a, b)
	}
	if merchantUid(7, "2026-07-29") == a {
		t.Error("different day must produce different uid")
	}
	if merchantUid(8, "2026-07-28") == a {
		t.Error("different sub must produce different uid")
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	body := []byte(`{"event":"paid"}`)
	secret := "webhook-secret"

	// 유효한 서명 만들기 (구현과 동일한 방식)
	valid := signForTest(body, secret)

	tests := []struct {
		name   string
		body   []byte
		sig    string
		secret string
		want   bool
	}{
		{"정상", body, valid, secret, true},
		{"본문 변조", []byte(`{"event":"refund"}`), valid, secret, false},
		{"서명 없음", body, "", secret, false},
		{"시크릿 미설정이면 무조건 거부", body, valid, "", false},
		{"엉뚱한 서명", body, "AAAA", secret, false},
	}
	for _, tt := range tests {
		if got := VerifyWebhookSignature(tt.body, tt.sig, tt.secret); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

// signForTest 는 VerifyWebhookSignature 와 동일한 HMAC-SHA256/base64 서명.
func signForTest(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
