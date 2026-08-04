-- P1 정기 구독결제 스키마 (docs/roadmap-p1-p6.md)
-- 실 DB 적용 전 백업 필수. 모두 additive 라 롤백은 DROP/컬럼 제거로 충분하다.

-- 회원 결제수단 (PG 빌링키). 카드번호 원문은 절대 저장하지 않는다 —
-- PG 의 customer_uid 와 표시용 메타(카드사/뒤 4자리)만 둔다.
CREATE TABLE IF NOT EXISTS `billingkey_tb` (
  `bk_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '빌링키 ID',
  `bk_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '사용자 ID (user_tb 참조)',
  `bk_pg` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'PG 구분 (portone/toss)',
  `bk_customer_uid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'PG 빌링키 식별자',
  `bk_card_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '카드사명 (표시용)',
  `bk_card_last4` varchar(4) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '카드 뒤 4자리 (표시용)',
  `bk_status` int(11) NOT NULL DEFAULT 1 COMMENT '상태 (1:active, 2:deleted)',
  `bk_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`bk_id`),
  INDEX idx_bk_user (bk_user),
  UNIQUE KEY uq_bk_customer_uid (bk_customer_uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='결제수단(빌링키) 테이블';

-- 구독 계약. 갱신 대상은 usehealth_tb(실제 회원권) — uh_endday 를 연장한다.
CREATE TABLE IF NOT EXISTS `subscription_tb` (
  `s_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '구독 ID',
  `s_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '사용자 ID (user_tb 참조)',
  `s_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `s_health` bigint(20) NOT NULL DEFAULT 0 COMMENT '운동권 ID (health_tb 참조)',
  `s_usehealth` bigint(20) NOT NULL DEFAULT 0 COMMENT '갱신 대상 회원권 ID (usehealth_tb 참조)',
  `s_billingkey` bigint(20) NOT NULL DEFAULT 0 COMMENT '빌링키 ID (billingkey_tb 참조)',
  `s_cost` int(11) NOT NULL DEFAULT 0 COMMENT '회차 결제 금액',
  `s_status` int(11) NOT NULL DEFAULT 1 COMMENT '상태 (1:active, 2:paused, 3:cancelled, 4:failed)',
  `s_next_billing_date` date NOT NULL COMMENT '다음 결제일',
  `s_fail_count` int(11) NOT NULL DEFAULT 0 COMMENT '연속 실패 횟수 (3회면 failed 전환)',
  `s_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`s_id`),
  INDEX idx_s_user (s_user),
  INDEX idx_s_gym (s_gym),
  INDEX idx_s_billing (s_status, s_next_billing_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='구독 계약 테이블';

-- payment_tb 확장: 상태·수단·PG 거래 ID·구독 참조·실패 사유.
-- 기존 행은 기본값(성공/카드)으로 남는다 — 지금까지의 수동 기록과 의미가 같다.
ALTER TABLE `payment_tb`
  ADD COLUMN `p_status` int(11) NOT NULL DEFAULT 1 COMMENT '결제 상태 (1:success, 2:pending, 3:fail, 4:cancelled, 5:refunded)',
  ADD COLUMN `p_method` int(11) NOT NULL DEFAULT 1 COMMENT '결제 수단 (1:card, 2:transfer, 3:cash, 4:billing)',
  ADD COLUMN `p_pg_tid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'PG 거래 ID (멱등성 키)',
  ADD COLUMN `p_subscription` bigint(20) NOT NULL DEFAULT 0 COMMENT '구독 결제면 subscription_tb 참조',
  ADD COLUMN `p_fail_reason` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '실패 사유';

-- 멱등성: 같은 PG 거래가 두 번 기록되지 않게 한다 (빈 문자열은 유니크 제외 불가라
-- 인덱스만 — 중복 차단은 앱 레벨에서 p_pg_tid 조회로 수행).
ALTER TABLE `payment_tb` ADD INDEX idx_p_pg_tid (p_pg_tid);
ALTER TABLE `payment_tb` ADD INDEX idx_p_subscription (p_subscription);
