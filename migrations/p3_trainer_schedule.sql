-- P3 트레이너 근무 스케줄 (docs/roadmap-p1-p6.md, 스키마는 front/flow.md 제안 채택)
CREATE TABLE IF NOT EXISTS `trainerschedule_tb` (
  `ts_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '스케줄 ID',
  `ts_trainer` bigint(20) NOT NULL DEFAULT 0 COMMENT '트레이너 ID (user_tb 참조)',
  `ts_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `ts_day_of_week` int(11) NOT NULL DEFAULT 0 COMMENT '요일 (0:일 ~ 6:토)',
  `ts_start_time` time NOT NULL COMMENT '근무 시작',
  `ts_end_time` time NOT NULL COMMENT '근무 종료',
  `ts_break_start` time DEFAULT NULL COMMENT '휴게 시작',
  `ts_break_end` time DEFAULT NULL COMMENT '휴게 종료',
  `ts_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`ts_id`),
  INDEX idx_ts_trainer (ts_trainer, ts_day_of_week),
  INDEX idx_ts_gym (ts_gym)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='트레이너 근무 스케줄 테이블';
