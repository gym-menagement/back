#!/usr/bin/env bash
# gym DB 계정 비밀번호 로테이션 헬퍼. 서버에서 실행한다.
# 자세한 배경·롤백은 docs/db-password-rotation.md 참고.
#
# 시크릿은 인자로 받지 않고 환경변수/프롬프트로 받는다 (프로세스 목록·셸
# 히스토리 노출 방지). 아무것도 하드코딩하지 않는다.
#
# 필요 환경변수:
#   DB_CONTAINER   DB 컨테이너 이름 (예: gym_db 또는 mariadb)
#   DB_ROOT_PASS   DB root 비밀번호 (미설정 시 프롬프트)
#   ENV_FILE       배포 env 경로 (기본 /data/gym_go/.env)
#   NEW_PASS       새 비밀번호 (미설정 시 자동 생성)
#
# 사용:
#   DB_CONTAINER=gym_db ./rotate_db_password.sh            # 자동 생성 + 대화형
#   DB_CONTAINER=gym_db NEW_PASS="$MY_SECRET" ./rotate_db_password.sh --apply  # secret-scan: ok 예시 플레이스홀더
#
# 기본은 dry-run(계획만 출력). 실제 변경은 --apply 를 붙일 때만.

set -euo pipefail

APPLY=0
[ "${1:-}" = "--apply" ] && APPLY=1

DB_USER=gym
ENV_FILE=${ENV_FILE:-/data/gym_go/.env}
: "${DB_CONTAINER:?DB_CONTAINER 환경변수가 필요하다 (DB 컨테이너 이름)}"

# root 비밀번호 — 미설정 시 조용히 프롬프트(입력 감춤).
if [ -z "${DB_ROOT_PASS:-}" ]; then
  read -rsp "DB root 비밀번호: " DB_ROOT_PASS; echo
fi

# 새 비밀번호 — 미설정 시 셸/.env 안전 문자로만 32자 생성.
if [ -z "${NEW_PASS:-}" ]; then
  NEW_PASS=$(openssl rand -base64 48 | tr -dc 'A-Za-z0-9' | head -c 32)
  echo "새 비밀번호를 자동 생성했다 (32자)."
fi

mysql_root() { docker exec -i "$DB_CONTAINER" mysql -uroot -p"$DB_ROOT_PASS" "$@"; }

echo "== 1) 현재 '$DB_USER' 계정 호스트 그레인트 =="
HOSTS=$(mysql_root -N -e "SELECT host FROM mysql.user WHERE user='$DB_USER';") || {
  echo "root 접속 또는 조회 실패 — DB_ROOT_PASS/DB_CONTAINER 확인"; exit 1; }
if [ -z "$HOSTS" ]; then echo "'$DB_USER' 계정을 찾지 못했다."; exit 1; fi
echo "$HOSTS" | sed 's/^/  gym@/'

if [ "$APPLY" -eq 0 ]; then
  echo
  echo "[dry-run] 위 호스트 전부에 ALTER USER 를 적용하고 $ENV_FILE 를 갱신할 예정."
  echo "[dry-run] 실제로 바꾸려면 --apply 를 붙여 다시 실행한다."
  exit 0
fi

echo "== 2) 비밀번호 교체 =="
SQL=""
while IFS= read -r h; do
  [ -z "$h" ] && continue
  SQL+="ALTER USER '$DB_USER'@'$h' IDENTIFIED BY '$NEW_PASS'; "
done <<< "$HOSTS"
SQL+="FLUSH PRIVILEGES;"
mysql_root -e "$SQL"
echo "  DB 비밀번호 교체 완료."

echo "== 3) 배포 env 갱신: $ENV_FILE =="
if [ ! -f "$ENV_FILE" ]; then echo "경고: $ENV_FILE 없음 — 수동 갱신 필요"; else
  cp "$ENV_FILE" "$ENV_FILE.bak.$(date +%s)"
  # GYM_DB_PASS / DB_PASS 둘 다 있으면 둘 다 교체.
  for key in GYM_DB_PASS DB_PASS; do
    if grep -q "^$key=" "$ENV_FILE"; then
      sed -i "s|^$key=.*|$key=$NEW_PASS|" "$ENV_FILE"
      echo "  $key 갱신"
    fi
  done
fi

echo "== 4) 컨테이너 재기동 후 아래를 수동 확인 =="
echo "  cd \$(dirname $ENV_FILE) && docker-compose up -d gym_api"
echo "  docker logs --tail 30 gym_api        # DB 연결 에러 없는지"
echo "  curl -s 'localhost:8004/api/gym?page=1&pagesize=1' | head -c 200"
echo
echo "새 비밀번호는 비밀 관리처에만 보관하고 이 터미널 기록을 지운다."
