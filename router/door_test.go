package router

import (
	"strings"
	"testing"
	"time"
)

// 도어 QR 은 무인 출입의 보안 경계다 — 위조·만료·변조가 전부 거부되는지 본다.
func TestDoorQRRoundtrip(t *testing.T) {
	now := time.Date(2026, 7, 28, 10, 0, 0, 0, time.UTC)
	qr := MakeDoorQR(42, now, 30*time.Second)

	// 정상: TTL 안이면 통과
	if id, ok := parseDoorQR(qr, now.Add(10*time.Second)); !ok || id != 42 {
		t.Fatalf("valid QR rejected: id=%d ok=%v", id, ok)
	}

	// 만료: TTL 지나면 거부
	if _, ok := parseDoorQR(qr, now.Add(31*time.Second)); ok {
		t.Error("expired QR must be rejected")
	}

	// 변조: usehealth ID 바꿔치기
	tampered := strings.Replace(qr, "uh:42:", "uh:43:", 1)
	if _, ok := parseDoorQR(tampered, now); ok {
		t.Error("tampered QR must be rejected")
	}

	// 형식 오류들
	for _, bad := range []string{"", "uh:42", "xx:42:99999:aaaa", "uh:42:notanumber:aaaa", qr + ":extra"} {
		if _, ok := parseDoorQR(bad, now); ok {
			t.Errorf("malformed QR accepted: %q", bad)
		}
	}
}

func TestDoorSign(t *testing.T) {
	body := []byte(`{"qr":"uh:1:2:abcd"}`)
	a := doorSign("secret-a", body)
	if a != doorSign("secret-a", body) {
		t.Error("same secret+body must produce same signature")
	}
	if a == doorSign("secret-b", body) {
		t.Error("different secret must produce different signature")
	}
	if a == doorSign("secret-a", []byte(`{"qr":"other"}`)) {
		t.Error("different body must produce different signature")
	}
}
