package billing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"

	"resty.dev/v3"
)

// PortOne V2 REST 구현. 키는 환경변수로만 받는다:
//
//	PORTONE_API_SECRET  — 결제 API 시크릿 (없으면 청구 시도가 즉시 실패)
//	PORTONE_WEBHOOK_SECRET — 웹훅 서명 검증 키
//	PORTONE_STORE_ID / PORTONE_CHANNEL_KEY — 상점·채널 식별
//
// 아직 PG 계약 전이라 엔드포인트는 V2 문서 기준으로 잡아 두고, 계약 후
// 샌드박스에서 E2E 로 확정한다 (roadmap P1 검증 절차).
type PortOne struct {
	APISecret  string
	StoreID    string
	ChannelKey string
	BaseURL    string
}

func NewPortOneFromEnv() *PortOne {
	base := os.Getenv("PORTONE_BASE_URL")
	if base == "" {
		base = "https://api.portone.io"
	}
	return &PortOne{
		APISecret:  os.Getenv("PORTONE_API_SECRET"),
		StoreID:    os.Getenv("PORTONE_STORE_ID"),
		ChannelKey: os.Getenv("PORTONE_CHANNEL_KEY"),
		BaseURL:    base,
	}
}

func (p *PortOne) Charge(customerUid string, merchantUid string, amount int, name string) (string, error) {
	if p.APISecret == "" {
		return "", errors.New("PORTONE_API_SECRET not configured")
	}

	client := resty.New()
	defer client.Close()

	var result struct {
		Payment struct {
			TransactionId string `json:"transactionId"`
		} `json:"payment"`
		Message string `json:"message"`
	}

	// V2: 빌링키 결제. paymentId 가 우리 쪽 멱등성 키(merchantUid)다.
	resp, err := client.R().
		SetHeader("Authorization", "PortOne "+p.APISecret).
		SetBody(map[string]interface{}{
			"billingKey": customerUid,
			"orderName":  name,
			"amount":     map[string]int{"total": amount},
			"currency":   "KRW",
			"storeId":    p.StoreID,
			"channelKey": p.ChannelKey,
		}).
		SetResult(&result).
		Post(p.BaseURL + "/payments/" + merchantUid + "/billing-key")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", errors.New("pg error: " + resp.String())
	}
	tid := result.Payment.TransactionId
	if tid == "" {
		// 응답 형태가 예상과 달라도 성공 코드면 merchantUid 로 추적 가능하게 둔다.
		tid = merchantUid
	}
	return tid, nil
}

// VerifyWebhookSignature 는 PG 웹훅 본문의 HMAC-SHA256 서명을 검증한다.
// 서명이 다르면 위조 — 무인증 경로이므로 이 검증이 유일한 방어선이다.
func VerifyWebhookSignature(body []byte, signature string, secret string) bool {
	if secret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
