-- P2 출입 도어 제어 (docs/roadmap-p1-p6.md)
-- 키오스크 없이 자동문 + 리더기. 리더기는 도어별 시크릿(HMAC)으로 인증한다 —
-- JWT 를 리더기에 주면 탈취 시 회원 API 전체가 열리므로 전용 시크릿으로 격리.

CREATE TABLE IF NOT EXISTS `door_tb` (
  `d_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '도어 ID',
  `d_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `d_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '도어 이름 (정문 등)',
  `d_secret` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '리더기 HMAC 시크릿',
  `d_status` int(11) NOT NULL DEFAULT 1 COMMENT '상태 (1:active, 2:disabled)',
  `d_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`d_id`),
  INDEX idx_d_gym (d_gym)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='출입 도어 테이블';
