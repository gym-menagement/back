#!/usr/bin/env bash
#
# gym Go 백엔드 라우팅 정렬 검증 (프론트 REST 계약 일치 확인)
# ---------------------------------------------------------------
# 실제 gym DB에 연결해 서버를 띄우고, 프론트가 호출하는 라우트가
# 404/405 없이 도달하는지(=라우팅 일치) 확인한다.
# 데이터 검증까지 하려면 GYM_LOGINID / GYM_PASSWD 를 넘긴다.
#
# DB 접속 정보는 환경변수로만 받는다. 예전엔 기본값에 실제 자격증명이 박혀
# 있었는데, 이 저장소는 공개라 그대로 유출됐다. 기본값을 두지 않고 없으면
# 즉시 중단해서 같은 실수가 반복되지 않게 한다.
#
#   DB_HOST=<host> DB_USER=<user> DB_PASS=<pass> DB_NAME=<db> ./verify_routes.sh
#   (선택) GYM_LOGINID=<id> GYM_PASSWD=<pw> 추가 시 로그인+목록까지 확인
set -uo pipefail
cd "$(dirname "$0")"

for v in DB_HOST DB_USER DB_PASS DB_NAME; do
  if [ -z "${!v:-}" ]; then
    echo "환경변수 $v 가 필요하다. 사용법은 이 파일 상단 주석 참고." >&2
    exit 2
  fi
done

PORT=${PORT:-8104}
BASE="http://127.0.0.1:${PORT}/api"
PID=""
pass=0; fail=0
ok(){ echo "  ✅ $1"; pass=$((pass+1)); }
bad(){ echo "  ❌ $1"; fail=$((fail+1)); }
cleanup(){ [ -n "$PID" ] && kill "$PID" 2>/dev/null; rm -f /tmp/gym_verify.bin /tmp/gym_verify.log; }
trap cleanup EXIT

echo "── 빌드 ──"
GOWORK=off CGO_ENABLED=0 go build -o /tmp/gym_verify.bin . || { echo "빌드 실패"; exit 1; }
ok "빌드"

echo "── 서버 실행 (DB: ${DB_HOST}/${DB_NAME}) ──"
APP_MODE=production PORT=$PORT \
  DB_TYPE=mysql DB_HOST="$DB_HOST" DB_PORT=${DB_PORT:-3306} \
  DB_NAME="$DB_NAME" DB_USER="$DB_USER" DB_PASS="$DB_PASS" \
  /tmp/gym_verify.bin >/tmp/gym_verify.log 2>&1 &
PID=$!
sleep 3
kill -0 "$PID" 2>/dev/null || { echo "서버 시작 실패:"; tail -20 /tmp/gym_verify.log; exit 1; }

# 라우팅 도달 확인: 404(라우트 없음)/405(메서드 불가)가 아니면 통과
code(){ curl -s -o /dev/null -w '%{http_code}' "$@"; }
reach(){ # $1 desc, 나머지 curl 인자
  local desc=$1; shift
  local c; c=$(code "$@")
  if [ "$c" = 404 ] || [ "$c" = 405 ]; then bad "$desc → $c (라우트 미일치)"; else ok "$desc → $c (도달)"; fi
}

echo "── 라우팅 일치 검증 (프론트 호출 패턴) ──"
reach "GET  /membership"            "$BASE/membership"
reach "GET  /membership/count"      "$BASE/membership/count"
reach "GET  /membership/1"          "$BASE/membership/1"
reach "PUT  /membership/1"          -X PUT  "$BASE/membership/1" -H 'Content-Type: application/json' -d '{}'
reach "PATCH /membership/1"         -X PATCH "$BASE/membership/1" -H 'Content-Type: application/json' -d '{}'
reach "DELETE /membership/1"        -X DELETE "$BASE/membership/1"
reach "GET  /user/search/loginid"   "$BASE/user/search/loginid?loginid=test"
reach "GET  /user/search/email"     "$BASE/user/search/email?email=test"
reach "GET  /gymadmin/search/gym"   "$BASE/gymadmin/search/gym?gym=1"
reach "GET  /jwt"                   "$BASE/jwt?loginid=x&passwd=y"

if [ -n "${GYM_LOGINID:-}" ] && [ -n "${GYM_PASSWD:-}" ]; then
  echo "── 인증 + 목록 데이터 검증 ──"
  R=$(curl -s "$BASE/jwt?loginid=${GYM_LOGINID}&passwd=${GYM_PASSWD}")
  TOKEN=$(echo "$R" | python3 -c "import sys,json;print(json.load(sys.stdin).get('token',''))" 2>/dev/null)
  [ -n "$TOKEN" ] && ok "로그인 토큰 발급" || bad "로그인 ($R)"
  if [ -n "$TOKEN" ]; then
    L=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE/membership?page=0&pagesize=5")
    echo "$L" | grep -q '"items"' && ok "목록 응답(items) — 인터셉터가 content로 매핑" || bad "목록 ($L)"
  fi
fi

echo ""
echo "통과 $pass / 실패 $fail"
[ "$fail" -eq 0 ] && echo "🎉 라우팅이 프론트 REST 계약과 일치" || { echo "서버로그:"; tail -15 /tmp/gym_verify.log; }
exit $fail
