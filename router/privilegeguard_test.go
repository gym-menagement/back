package router

import "testing"

// 권한 상승 판정의 순수 로직 회귀. 실제 차단은 DB 의 현재 값과 비교하므로
// 여기서는 값 파싱(문자열 우회 포함)과 상한 규칙을 고정한다.

func TestPrivilegeInt(t *testing.T) {
	tests := []struct {
		name  string
		in    interface{}
		want  int
		valid bool
	}{
		{"숫자", float64(5), 5, true},
		{"문자열 우회", "5", 5, true},
		{"숫자 아닌 문자열", "admin", 0, false},
		{"불리언", true, 0, false},
		{"nil", nil, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := privilegeInt(tt.in)
			if ok != tt.valid || (ok && got != tt.want) {
				t.Errorf("privilegeInt(%v) = (%d, %v), want (%d, %v)", tt.in, got, ok, tt.want, tt.valid)
			}
		})
	}
}

// 가입 상한과 수정 상한이 같은 값을 써야 한다 — 한쪽만 바뀌면 우회로가 생긴다.
// (무인증 가입은 3/4 로 막아두고 인증 수정은 5/5 를 허용하던 것이 실제 취약점이었다)
func TestEscalationCapsMatchSignupCaps(t *testing.T) {
	if maxSignupLevel != 3 || maxSignupRole != 4 {
		t.Errorf("가입 상한이 바뀌었다 (level=%d role=%d) — privilegeguard.go 의 상한도 함께 확인할 것",
			maxSignupLevel, maxSignupRole)
	}
}
