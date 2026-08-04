-- P5 전자계약 (docs/roadmap-p1-p6.md)
-- 서명 자체는 외부 서비스(모두싸인류)에 위임하고 우리는 연결 기록만 갖는다.
CREATE TABLE IF NOT EXISTS `contract_tb` (
  `ct_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '계약 ID',
  `ct_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '사용자 ID (user_tb 참조)',
  `ct_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `ct_order` bigint(20) NOT NULL DEFAULT 0 COMMENT '주문 ID (order_tb 참조)',
  `ct_provider` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '서명 서비스 (modusign 등)',
  `ct_document_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '외부 문서 ID',
  `ct_sign_url` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '서명 링크',
  `ct_status` int(11) NOT NULL DEFAULT 1 COMMENT '상태 (1:requested, 2:signed, 3:declined, 4:expired)',
  `ct_signed_date` datetime DEFAULT NULL COMMENT '서명 완료일',
  `ct_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`ct_id`),
  INDEX idx_ct_user (ct_user),
  INDEX idx_ct_gym (ct_gym),
  UNIQUE KEY uq_ct_document (ct_provider, ct_document_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='전자계약 테이블';
