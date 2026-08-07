-- PT 예약 이중예약 방지 (DB 레벨)
--
-- 애플리케이션 가드(router/reservationguard.go)는 SELECT 로 겹침을 확인한 뒤
-- INSERT 하므로 그 사이에 끼어드는 동시 요청을 막지 못한다 — 같은 슬롯에
-- 동시 5건을 던지면 5건이 전부 생성됐다(2026-08-07 5회차 점검).
--
-- 생성 컬럼 + 유니크 키로 DB 가 최종 판정하게 한다:
--   - 예약중(status=1) 행만 슬롯 키를 갖고 → 같은 트레이너·날짜·시각은 1건뿐
--   - 취소·완료·노쇼는 NULL → 유니크 제약에서 제외되므로 이력이 여러 건 남아도 무방
--     (같은 슬롯을 예약→취소→재예약 하는 정상 흐름이 막히지 않는다)
--
-- 적용 전 중복 정리가 필요하다. 아래 DELETE 는 각 슬롯의 가장 오래된 1건만
-- 남긴다 — 먼저 잡은 예약이 유효하다는 규칙.

DELETE r FROM ptreservation_tb r
JOIN (
  SELECT MIN(pr_id) AS keep_id, pr_trainer, pr_reservationdate, pr_starttime
  FROM ptreservation_tb
  WHERE pr_status = 1
  GROUP BY pr_trainer, pr_reservationdate, pr_starttime
  HAVING COUNT(*) > 1
) d ON r.pr_trainer = d.pr_trainer
   AND r.pr_reservationdate = d.pr_reservationdate
   AND r.pr_starttime = d.pr_starttime
   AND r.pr_status = 1
   AND r.pr_id <> d.keep_id;

ALTER TABLE ptreservation_tb
  ADD COLUMN `pr_active_slot` varchar(64)
    GENERATED ALWAYS AS (
      IF(pr_status = 1,
         CONCAT(pr_trainer, '_', pr_reservationdate, '_', pr_starttime),
         NULL)
    ) STORED COMMENT '예약중 슬롯 유니크 키 (취소·완료는 NULL)',
  ADD UNIQUE KEY `uq_pr_active_slot` (`pr_active_slot`);
