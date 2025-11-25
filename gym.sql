-- MariaDB dump 10.19  Distrib 10.8.2-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: gym
-- ------------------------------------------------------
-- Server version	10.8.2-MariaDB-1:10.8.2+maria~focal

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `alarm_tb`
--

DROP TABLE IF EXISTS `alarm_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alarm_tb` (
  `al_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '알람 ID',
  `al_title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '알람 제목',
  `al_content` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '알람 내용',
  `al_type` int(11) NOT NULL COMMENT '알람 타입 (NOTICE:공지, WARNING:경고, ERROR:에러, INFO:정보)',
  `al_status` int(11) NOT NULL COMMENT '알람 상태 (SUCCESS:성공, FAIL:실패, PENDING:대기)',
  `al_user` bigint(20) NOT NULL COMMENT '대상 회원 ID (user_tb 참조)',
  `al_date` datetime NOT NULL COMMENT '등록일',
  PRIMARY KEY (`al_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='알람 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `daytype_tb`
--

DROP TABLE IF EXISTS `daytype_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `daytype_tb` (
  `dt_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '요일 타입 ID',
  `dt_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `dt_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '요일 타입명 (평일/주말/공휴일 등)',
  `dt_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`dt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='요일 타입 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `discount_tb`
--

DROP TABLE IF EXISTS `discount_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `discount_tb` (
  `d_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '할인 ID',
  `d_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '할인명',
  `d_discount` int(11) NOT NULL DEFAULT 0 COMMENT '할인율(%)',
  `d_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`d_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='할인 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gym_tb`
--

DROP TABLE IF EXISTS `gym_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gym_tb` (
  `g_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '헬스장 ID',
  `g_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '헬스장명',
  `g_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '헬스장 주소',
  `g_tel` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '헬스장 전화번호',
  `g_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '소유자 ID (user_tb 참조)',
  `g_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`g_id`),
  INDEX idx_user (g_user)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='헬스장 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `health_tb`
--

DROP TABLE IF EXISTS `health_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `health_tb` (
  `h_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '운동권 ID',
  `h_category` bigint(20) NOT NULL DEFAULT 0 COMMENT '운동권 카테고리 ID (healthcategory_tb 참조)',
  `h_term` bigint(20) NOT NULL DEFAULT 0 COMMENT '기간 ID (term_tb 참조)',
  `h_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '운동권명',
  `h_count` int(11) NOT NULL DEFAULT 0 COMMENT '이용 가능 횟수',
  `h_cost` int(11) NOT NULL DEFAULT 0 COMMENT '원가',
  `h_discount` bigint(20) NOT NULL DEFAULT 0 COMMENT '할인 ID (discount_tb 참조)',
  `h_costdiscount` int(11) NOT NULL DEFAULT 0 COMMENT '할인가',
  `h_content` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '운동권 설명',
  `h_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `h_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`h_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='운동권 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `healthcategory_tb`
--

DROP TABLE IF EXISTS `healthcategory_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `healthcategory_tb` (
  `hc_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '운동권 카테고리 ID',
  `hc_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `hc_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '카테고리명 (헬스/PT/요가 등)',
  `hc_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`hc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='운동권 카테고리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ipblock_tb`
--

DROP TABLE IF EXISTS `ipblock_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipblock_tb` (
  `ib_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'IP 차단 ID',
  `ib_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'IP 주소',
  `ib_type` int(11) NOT NULL DEFAULT 0 COMMENT 'IP 타입 (, admin:관리자 접근, normal:일반 접근)',
  `ib_policy` int(11) NOT NULL DEFAULT 0 COMMENT '정책 (grant:허용, DENY:거부)',
  `ib_use` int(11) NOT NULL DEFAULT 0 COMMENT '사용 여부 (USE:사용, NOTUSE:미사용)',
  `ib_order` int(11) NOT NULL DEFAULT 0 COMMENT '우선순위',
  `ib_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`ib_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='IP 차단 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `loginlog_tb`
--

DROP TABLE IF EXISTS `loginlog_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `loginlog_tb` (
  `ll_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '로그인 로그 ID',
  `ll_ip` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '접속 IP 주소',
  `ll_ipvalue` bigint(20) NOT NULL DEFAULT 0 COMMENT 'IP 주소 숫자값',
  `ll_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
  `ll_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '로그인 일시',
  PRIMARY KEY (`ll_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='로그인 로그 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `membership_tb`
--

DROP TABLE IF EXISTS `membership_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `membership_tb` (
  `m_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '회원 ID',
  `m_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `m_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '사용자 ID (user_tb 참조)',
  `m_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '회원명',
  `m_sex` int(11) NOT NULL DEFAULT 0 COMMENT '성별 (MALE:남성, FEMALE:여성)',
  `m_birth` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '생년월일',
  `m_phonenum` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '전화번호',
  `m_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '주소',
  `m_image` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '프로필 이미지',
  `m_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`m_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='회원 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_tb`
--

DROP TABLE IF EXISTS `order_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_tb` (
  `o_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '주문 ID',
  `o_membership` bigint(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (membership_tb 참조)',
  `o_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '주문일',
  PRIMARY KEY (`o_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='주문 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payment_tb`
--

DROP TABLE IF EXISTS `payment_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_tb` (
  `p_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '결제 ID',
  `p_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `p_order` bigint(20) NOT NULL DEFAULT 0 COMMENT '주문 ID (order_tb 참조)',
  `p_membership` bigint(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (membership_tb 참조)',
  `p_cost` int(11) NOT NULL DEFAULT 0 COMMENT '결제 금액',
  `p_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '결제일',
  PRIMARY KEY (`p_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='결제 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `paymentform_tb`
--

DROP TABLE IF EXISTS `paymentform_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `paymentform_tb` (
  `pf_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '결제 상세 ID',
  `pf_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `pf_payment` bigint(20) NOT NULL DEFAULT 0 COMMENT '결제 ID (payment_tb 참조)',
  `pf_type` bigint(20) NOT NULL DEFAULT 0 COMMENT '결제 타입 ID (paymenttype_tb 참조)',
  `pf_cost` int(11) NOT NULL DEFAULT 0 COMMENT '결제 금액',
  `pf_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`pf_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='결제 상세 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `paymenttype_tb`
--

DROP TABLE IF EXISTS `paymenttype_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `paymenttype_tb` (
  `pt_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '결제 타입 ID',
  `pt_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `pt_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '결제 타입명 (현금/카드/계좌이체 등)',
  `pt_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`pt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='결제 타입 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rocker_tb`
--

DROP TABLE IF EXISTS `rocker_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rocker_tb` (
  `r_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '락커 ID',
  `r_group` bigint(20) NOT NULL DEFAULT 0 COMMENT '락커 그룹 ID (rockergroup_tb 참조)',
  `r_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '락커 번호/이름',
  `r_available` int(11) NOT NULL DEFAULT 0 COMMENT '사용 가능 여부 (IN_USE:사용중, AVAILABLE:사용가능)',
  `r_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`r_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='락커 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rockergroup_tb`
--

DROP TABLE IF EXISTS `rockergroup_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rockergroup_tb` (
  `rg_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '락커 그룹 ID',
  `rg_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `rg_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '락커 그룹명 (남자락커실/여자락커실 등)',
  `rg_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`rg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='락커 그룹 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role_tb`
--

DROP TABLE IF EXISTS `role_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_tb` (
  `r_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '역할 ID',
  `r_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `r_roleid` int(11) NOT NULL DEFAULT 0 COMMENT '역할 타입 (MEMBER:회원, TRAINER:트레이너, STAFF:직원, GYM_ADMIN:헬스장관리자, PLATFORM_ADMIN:플랫폼관리자)',
  `r_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '역할명',
  `r_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`r_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='역할 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `setting_tb`
--

DROP TABLE IF EXISTS `setting_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `setting_tb` (
  `se_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '설정 ID',
  `se_category` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '설정 카테고리',
  `se_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '설정명',
  `se_key` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '설정 키',
  `se_value` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '설정 값',
  `se_remark` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '설명',
  `se_type` int(11) NOT NULL DEFAULT 0 COMMENT '설정 타입 (STRING:문자열, NUMBER:숫자, BOOLEAN:참거짓)',
  `se_data` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '추가 데이터 (JSON 등)',
  `se_order` int(11) NOT NULL DEFAULT 0 COMMENT '정렬 순서',
  `se_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`se_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='설정 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stop_tb`
--

DROP TABLE IF EXISTS `stop_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stop_tb` (
  `s_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '일시정지 ID',
  `s_usehelth` bigint(20) NOT NULL DEFAULT 0 COMMENT '운동권 사용 ID (usehealth_tb 참조)',
  `s_startday` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '일시정지 시작일',
  `s_endday` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '일시정지 종료일',
  `s_count` int(11) NOT NULL DEFAULT 0 COMMENT '일시정지 일수',
  `s_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='회원권 일시정지 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `systemlog_tb`
--

DROP TABLE IF EXISTS `systemlog_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `systemlog_tb` (
  `sl_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '시스템 로그 ID',
  `sl_type` int(11) NOT NULL DEFAULT 0 COMMENT '로그 타입 (LOGIN:로그인, CRAWLING:크롤링)',
  `sl_content` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '로그 내용',
  `sl_result` int(11) NOT NULL DEFAULT 0 COMMENT '결과 (SUCCESS:성공, FAIL:실패)',
  `sl_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '발생일시',
  PRIMARY KEY (`sl_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='시스템 로그 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_tb`
--

DROP TABLE IF EXISTS `term_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_tb` (
  `t_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '기간 ID',
  `t_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `t_daytype` bigint(20) NOT NULL DEFAULT 0 COMMENT '요일 타입 ID (daytype_tb 참조)',
  `t_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '기간명 (1개월/3개월/6개월 등)',
  `t_term` int(11) NOT NULL DEFAULT 0 COMMENT '기간 일수',
  `t_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`t_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='기간 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `token_tb`
--

DROP TABLE IF EXISTS `token_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `token_tb` (
  `to_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '토큰 ID',
  `to_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
  `to_token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'JWT 토큰',
  `to_status` int(11) NOT NULL DEFAULT 0 COMMENT '토큰 상태 (ACTIVE:활성, EXPIRED:만료, REVOKED:폐기)',
  `to_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '발급일시',
  PRIMARY KEY (`to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='인증 토큰 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `usehealth_tb`
--

DROP TABLE IF EXISTS `usehealth_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `usehealth_tb` (
  `uh_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '운동권 사용 ID',
  `uh_order` bigint(20) NOT NULL DEFAULT 0 COMMENT '주문 ID (order_tb 참조)',
  `uh_health` bigint(20) NOT NULL DEFAULT 0 COMMENT '운동권 ID (health_tb 참조)',
  `uh_user` bigint(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
  `uh_rocker` bigint(20) NOT NULL DEFAULT 0 COMMENT '락커 ID (rocker_tb 참조)',
  `uh_term` bigint(20) NOT NULL DEFAULT 0 COMMENT '기간 ID (term_tb 참조)',
  `uh_discount` bigint(20) NOT NULL DEFAULT 0 COMMENT '할인 ID (discount_tb 참조)',
  `uh_startday` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '시작일',
  `uh_endday` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '종료일',
  `uh_gym` bigint(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
  `uh_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
  PRIMARY KEY (`uh_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='운동권 사용 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_tb`
--

DROP TABLE IF EXISTS `user_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_tb` (
  `u_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '사용자 ID',
  `u_loginid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '로그인 ID',
  `u_passwd` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '비밀번호 (암호화)',
  `u_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '이메일',
  `u_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '이름',
  `u_tel` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '전화번호',
  `u_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '주소',
  `u_image` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '프로필 이미지',
  `u_sex` int(11) NOT NULL DEFAULT 0 COMMENT '성별 (MALE:남성, FEMALE:여성)',
  `u_birth` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '생년월일',
  `u_type` int(11) NOT NULL DEFAULT 0 COMMENT '가입 타입 (NORMAL:일반, KAKAO:카카오, NAVER:네이버, GOOGLE:구글, APPLE:애플)',
  `u_connectid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'SNS 연동 ID',
  `u_level` int(11) NOT NULL DEFAULT 0 COMMENT '권한 레벨 (normal:일반회원, manager:트레이너/직원, admin:헬스장관리자, superadmin:플랫폼관리자, rootadmin:최고관리자)',
  `u_role` int(11) NOT NULL DEFAULT 0 COMMENT '역할 (MEMBER:회원, TRAINER:트레이너, STAFF:직원, GYM_ADMIN:헬스장관리자, PLATFORM_ADMIN:플랫폼관리자)',
  `u_use` int(11) NOT NULL DEFAULT 0 COMMENT '사용 여부 (USE:사용, NOTUSE:미사용)',
  `u_logindate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '최근 로그인 일시',
  `u_lastchangepasswddate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '마지막 비밀번호 변경일',
  `u_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '가입일',
  PRIMARY KEY (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='사용자 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-10-01  3:16:46


-- ============================================
-- 헬스장 관리 어플리케이션 추가 테이블
-- ============================================

-- 1. 출석 체크 테이블 (attendance_tb)
-- 회원들의 출석/입장 기록을 저장
CREATE TABLE IF NOT EXISTS attendance_tb (
    at_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '출석 ID',
    at_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
    at_membership BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원권 ID (membership_tb 참조)',
    at_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
    at_type INT(11) NOT NULL DEFAULT 0 COMMENT '출석 타입 (ENTRY:입장, EXIT:퇴장)',
    at_method INT(11) NOT NULL DEFAULT 0 COMMENT '체크 방법 (QR_CODE:QR코드, MANUAL:수동, CARD:카드)',
    at_checkintime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '입장 시간',
    at_checkouttime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '퇴장 시간',
    at_duration INT(11) NOT NULL DEFAULT 0 COMMENT '체류 시간(분)',
    at_status INT(11) NOT NULL DEFAULT 0 COMMENT '상태 (NORMAL:정상, LATE:지각, UNAUTHORIZED:무단입장)',
    at_note TEXT NOT NULL COMMENT '비고',
    at_ip VARCHAR(50) NOT NULL DEFAULT '' COMMENT '접속 IP',
    at_device VARCHAR(100) NOT NULL DEFAULT '' COMMENT '디바이스 정보',
    at_createdby BIGINT(20) NOT NULL DEFAULT 0 COMMENT '등록자 ID',
    at_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (at_id),
    INDEX idx_user (at_user),
    INDEX idx_gym (at_gym),
    INDEX idx_checkintime (at_checkintime),
    INDEX idx_date (at_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='출석 체크 테이블';

-- 2. 회원 QR 코드 테이블 (memberqr_tb)
-- 각 회원의 고유 QR 코드 정보
CREATE TABLE IF NOT EXISTS memberqr_tb (
    mq_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'QR ID',
    mq_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
    mq_code VARCHAR(255) NOT NULL DEFAULT '' UNIQUE COMMENT 'QR 코드 값 (UUID)',
    mq_imageurl VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'QR 이미지 URL',
    mq_isactive INT(11) NOT NULL DEFAULT 1 COMMENT '활성 상태 (INACTIVE:비활성, ACTIVE:활성)',
    mq_expiredate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '만료일',
    mq_generateddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    mq_lastuseddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '마지막 사용일',
    mq_usecount INT(11) NOT NULL DEFAULT 0 COMMENT '사용 횟수',
    mq_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (mq_id),
    UNIQUE KEY uk_user (mq_user),
    UNIQUE KEY uk_code (mq_code),
    INDEX idx_code (mq_code),
    INDEX idx_active (mq_isactive)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='회원 QR 코드 테이블';

-- 3. 운동 기록 테이블 (workoutlog_tb)
-- 회원들의 운동 기록 (세트, 무게, 횟수 등)
CREATE TABLE IF NOT EXISTS workoutlog_tb (
    wl_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '운동기록 ID',
    wl_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    wl_attendance BIGINT(20) NOT NULL DEFAULT 0 COMMENT '출석 ID (attendance_tb 참조)',
    wl_health BIGINT(20) NOT NULL DEFAULT 0 COMMENT '운동기구 ID (health_tb 참조)',
    wl_exercisename VARCHAR(200) NOT NULL DEFAULT '' COMMENT '운동명',
    wl_sets INT(11) NOT NULL DEFAULT 0 COMMENT '세트 수',
    wl_reps INT(11) NOT NULL DEFAULT 0 COMMENT '반복 횟수',
    wl_weight DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '중량(kg)',
    wl_duration INT(11) NOT NULL DEFAULT 0 COMMENT '운동 시간(분)',
    wl_calories INT(11) NOT NULL DEFAULT 0 COMMENT '소모 칼로리',
    wl_note TEXT NOT NULL COMMENT '메모',
    wl_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (wl_id),
    INDEX idx_user (wl_user),
    INDEX idx_attendance (wl_attendance),
    INDEX idx_date (wl_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='운동 기록 테이블';

-- 4. 회원권 이용 내역 테이블 (membershipusage_tb)
-- 회원권별 사용 내역 (남은 횟수, 기간 등 추적)
CREATE TABLE IF NOT EXISTS membershipusage_tb (
    mu_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '이용내역 ID',
    mu_membership BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원권 ID',
    mu_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    mu_type INT(11) NOT NULL DEFAULT 0 COMMENT '회원권 타입 (PERIOD_BASED:기간제, COUNT_BASED:횟수제)',
    mu_totaldays INT(11) NOT NULL DEFAULT 0 COMMENT '총 이용 가능 일수',
    mu_useddays INT(11) NOT NULL DEFAULT 0 COMMENT '사용한 일수',
    mu_remainingdays INT(11) NOT NULL DEFAULT 0 COMMENT '남은 일수',
    mu_totalcount INT(11) NOT NULL DEFAULT 0 COMMENT '총 이용 가능 횟수',
    mu_usedcount INT(11) NOT NULL DEFAULT 0 COMMENT '사용한 횟수',
    mu_remainingcount INT(11) NOT NULL DEFAULT 0 COMMENT '남은 횟수',
    mu_startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '시작일',
    mu_enddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '종료일',
    mu_status INT(11) NOT NULL DEFAULT 0 COMMENT '상태 (IN_USE:사용중, PAUSED:일시정지, EXPIRED:만료, REFUNDED:환불)',
    mu_pausedays INT(11) NOT NULL DEFAULT 0 COMMENT '일시정지 일수',
    mu_lastuseddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '마지막 사용일',
    mu_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (mu_id),
    INDEX idx_membership (mu_membership),
    INDEX idx_user (mu_user),
    INDEX idx_status (mu_status),
    INDEX idx_enddate (mu_enddate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='회원권 이용 내역 테이블';

-- 5. 공지사항 테이블 (notice_tb)
-- 헬스장 공지사항 (웹/앱에서 조회)
CREATE TABLE IF NOT EXISTS notice_tb (
    nt_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '공지사항 ID',
    nt_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID',
    nt_title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '제목',
    nt_content TEXT NOT NULL COMMENT '내용',
    nt_type INT(11) NOT NULL DEFAULT 0 COMMENT '타입 (GENERAL:일반, IMPORTANT:중요, EVENT:이벤트)',
    nt_ispopup INT(11) NOT NULL DEFAULT 0 COMMENT '팝업 여부 (NO:아니오, YES:예)',
    nt_ispush INT(11) NOT NULL DEFAULT 0 COMMENT '푸시알림 여부 (NO:아니오, YES:예)',
    nt_target INT(11) NOT NULL DEFAULT 0 COMMENT '대상 (ALL:전체, MEMBERS_ONLY:회원만, SPECIFIC_MEMBERS:특정회원)',
    nt_viewcount INT(11) NOT NULL DEFAULT 0 COMMENT '조회수',
    nt_startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '게시 시작일',
    nt_enddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '게시 종료일',
    nt_status INT(11) NOT NULL DEFAULT 1 COMMENT '상태 (PRIVATE:비공개, PUBLIC:공개)',
    nt_createdby BIGINT(20) NOT NULL DEFAULT 0 COMMENT '작성자 ID',
    nt_createddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '작성일',
    nt_updateddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    nt_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (nt_id),
    INDEX idx_gym (nt_gym),
    INDEX idx_type (nt_type),
    INDEX idx_status (nt_status),
    INDEX idx_createddate (nt_createddate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='공지사항 테이블';

-- 6. 트레이너-회원 매칭 테이블 (trainermember_tb)
-- 트레이너와 담당 회원 매칭 관리
CREATE TABLE IF NOT EXISTS trainermember_tb (
    tm_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '매칭 ID',
    tm_trainer BIGINT(20) NOT NULL DEFAULT 0 COMMENT '트레이너 ID (user_tb 참조)',
    tm_member BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID (user_tb 참조)',
    tm_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID',
    tm_startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '담당 시작일',
    tm_enddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '담당 종료일',
    tm_status INT(11) NOT NULL DEFAULT 1 COMMENT '상태 (TERMINATED:종료, IN_PROGRESS:진행중)',
    tm_note TEXT NOT NULL COMMENT '비고',
    tm_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (tm_id),
    INDEX idx_trainer (tm_trainer),
    INDEX idx_member (tm_member),
    INDEX idx_status (tm_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='트레이너-회원 매칭 테이블';

-- 7. PT 수업 예약 테이블 (ptreservation_tb)
-- 개인 트레이닝 수업 예약 관리
CREATE TABLE IF NOT EXISTS ptreservation_tb (
    pr_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '예약 ID',
    pr_trainer BIGINT(20) NOT NULL DEFAULT 0 COMMENT '트레이너 ID',
    pr_member BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    pr_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID',
    pr_reservationdate DATE NOT NULL DEFAULT '1900-01-01' COMMENT '예약 날짜',
    pr_starttime TIME NOT NULL DEFAULT '00:00:00' COMMENT '시작 시간',
    pr_endtime TIME NOT NULL DEFAULT '00:00:00' COMMENT '종료 시간',
    pr_duration INT(11) NOT NULL DEFAULT 60 COMMENT '수업 시간(분)',
    pr_status INT(11) NOT NULL DEFAULT 0 COMMENT '상태 (RESERVED:예약, COMPLETED:완료, CANCELLED:취소, NO_SHOW:노쇼)',
    pr_note TEXT NOT NULL COMMENT '메모',
    pr_cancelreason TEXT NOT NULL COMMENT '취소 사유',
    pr_createddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '예약일',
    pr_updateddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    pr_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (pr_id),
    INDEX idx_trainer (pr_trainer),
    INDEX idx_member (pr_member),
    INDEX idx_reservationdate (pr_reservationdate),
    INDEX idx_status (pr_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PT 수업 예약 테이블';

-- 8. 회원 신체 정보 테이블 (memberbody_tb)
-- 회원들의 신체 정보 변화 추적 (체중, 체지방률 등)
CREATE TABLE IF NOT EXISTS memberbody_tb (
    mb_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '신체정보 ID',
    mb_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    mb_height DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '키(cm)',
    mb_weight DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '체중(kg)',
    mb_bodyfat DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '체지방률(%)',
    mb_musclemass DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '근육량(kg)',
    mb_bmi DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT 'BMI',
    mb_skeletalmuscle DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '골격근량(kg)',
    mb_bodywater DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '체수분(%)',
    mb_chest DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '가슴둘레(cm)',
    mb_waist DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '허리둘레(cm)',
    mb_hip DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '엉덩이둘레(cm)',
    mb_arm DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '팔둘레(cm)',
    mb_thigh DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '허벅지둘레(cm)',
    mb_note TEXT NOT NULL COMMENT '메모',
    mb_measureddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '측정일',
    mb_measuredby BIGINT(20) NOT NULL DEFAULT 0 COMMENT '측정자 ID',
    mb_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (mb_id),
    INDEX idx_user (mb_user),
    INDEX idx_measureddate (mb_measureddate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='회원 신체 정보 테이블';

-- 9. 문의/상담 테이블 (inquiry_tb)
-- 회원 문의 및 상담 내역
CREATE TABLE IF NOT EXISTS inquiry_tb (
    iq_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '문의 ID',
    iq_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    iq_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID',
    iq_type INT(11) NOT NULL DEFAULT 0 COMMENT '문의 유형 (GENERAL:일반, MEMBERSHIP:회원권, REFUND:환불, FACILITY:시설, OTHER:기타)',
    iq_title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '제목',
    iq_content TEXT NOT NULL COMMENT '내용',
    iq_status INT(11) NOT NULL DEFAULT 0 COMMENT '상태 (WAITING:대기, ANSWERED:답변완료)',
    iq_answer TEXT NOT NULL COMMENT '답변',
    iq_answeredby BIGINT(20) NOT NULL DEFAULT 0 COMMENT '답변자 ID',
    iq_answereddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '답변일',
    iq_createddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '작성일',
    iq_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (iq_id),
    INDEX idx_user (iq_user),
    INDEX idx_gym (iq_gym),
    INDEX idx_status (iq_status),
    INDEX idx_createddate (iq_createddate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='문의/상담 테이블';

-- 10. 락커 사용 내역 테이블 (rockerusage_tb)
-- 락커 배정 및 사용 내역
CREATE TABLE IF NOT EXISTS rockerusage_tb (
    ru_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '사용내역 ID',
    ru_rocker BIGINT(20) NOT NULL DEFAULT 0 COMMENT '락커 ID',
    ru_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    ru_membership BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원권 ID',
    ru_startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '시작일',
    ru_enddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '종료일',
    ru_status INT(11) NOT NULL DEFAULT 1 COMMENT '상태 (TERMINATED:종료, IN_USE:사용중, OVERDUE:연체)',
    ru_deposit DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '보증금',
    ru_monthlyfee DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '월 이용료',
    ru_note TEXT NOT NULL COMMENT '비고',
    ru_assignedby BIGINT(20) NOT NULL DEFAULT 0 COMMENT '배정자 ID',
    ru_assigneddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '배정일',
    ru_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (ru_id),
    INDEX idx_rocker (ru_rocker),
    INDEX idx_user (ru_user),
    INDEX idx_status (ru_status),
    INDEX idx_enddate (ru_enddate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='락커 사용 내역 테이블';

-- 11. 푸시 알림 토큰 테이블 (pushtoken_tb)
-- 모바일 앱 푸시 알림용 디바이스 토큰
CREATE TABLE IF NOT EXISTS pushtoken_tb (
    pt_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '토큰 ID',
    pt_user BIGINT(20) NOT NULL DEFAULT 0 COMMENT '회원 ID',
    pt_token VARCHAR(500) NOT NULL DEFAULT '' COMMENT '푸시 토큰',
    pt_devicetype VARCHAR(20) NOT NULL DEFAULT '' COMMENT '디바이스 타입 (ios, android)',
    pt_deviceid VARCHAR(255) NOT NULL DEFAULT '' COMMENT '디바이스 ID',
    pt_appversion VARCHAR(50) NOT NULL DEFAULT '' COMMENT '앱 버전',
    pt_isactive INT(11) NOT NULL DEFAULT 1 COMMENT '활성 상태 (INACTIVE:비활성, ACTIVE:활성)',
    pt_createddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    pt_updateddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '갱신일',
    pt_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (pt_id),
    INDEX idx_user (pt_user),
    INDEX idx_token (pt_token(255)),
    INDEX idx_active (pt_isactive)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='푸시 알림 토큰 테이블';

-- 12. 앱 버전 관리 테이블 (appversion_tb)
-- 모바일 앱 버전 관리 및 강제 업데이트
CREATE TABLE IF NOT EXISTS appversion_tb (
    av_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '버전 ID',
    av_platform VARCHAR(20) NOT NULL DEFAULT '' COMMENT '플랫폼 (ios, android)',
    av_version VARCHAR(50) NOT NULL DEFAULT '' COMMENT '버전',
    av_minversion VARCHAR(50) NOT NULL DEFAULT '' COMMENT '최소 요구 버전',
    av_forceupdate INT(11) NOT NULL DEFAULT 0 COMMENT '강제 업데이트 (NO:아니오, YES:예)',
    av_updatemessage TEXT NOT NULL COMMENT '업데이트 안내 메시지',
    av_downloadurl VARCHAR(500) NOT NULL DEFAULT '' COMMENT '다운로드 URL',
    av_status INT(11) NOT NULL DEFAULT 1 COMMENT '상태 (INACTIVE:비활성, ACTIVE:활성)',
    av_releasedate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '배포일',
    av_createddate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    av_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (av_id),
    INDEX idx_platform (av_platform),
    INDEX idx_status (av_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='앱 버전 관리 테이블';

-- 13. 헬스장-트레이너 매칭 테이블 (gymtrainer_tb)
-- 트레이너가 어느 헬스장에서 근무하는지 관리 (트레이너는 여러 헬스장에서 근무 가능)
CREATE TABLE IF NOT EXISTS gymtrainer_tb (
    gt_id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '헬스장-트레이너 매칭 ID',
    gt_gym BIGINT(20) NOT NULL DEFAULT 0 COMMENT '헬스장 ID (gym_tb 참조)',
    gt_trainer BIGINT(20) NOT NULL DEFAULT 0 COMMENT '트레이너 ID (user_tb 참조)',
    gt_startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '근무 시작일',
    gt_enddate DATETIME NULL COMMENT '근무 종료일 (NULL: 현재 근무 중)',
    gt_status INT(11) NOT NULL DEFAULT 1 COMMENT '상태 (TERMINATED:종료, IN_PROGRESS:진행중)',
    gt_position VARCHAR(100) NOT NULL DEFAULT '' COMMENT '직책 (헤드 트레이너, 일반 트레이너 등)',
    gt_note TEXT NOT NULL COMMENT '비고',
    gt_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '등록일',
    PRIMARY KEY (gt_id),
    INDEX idx_gym (gt_gym),
    INDEX idx_trainer (gt_trainer),
    INDEX idx_status (gt_status),
    INDEX idx_gym_trainer (gt_gym, gt_trainer)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='헬스장-트레이너 매칭 테이블';

-- ============================================
-- 샘플 데이터 삽입 (테스트용)
-- ============================================

-- -- QR 코드 샘플 (user_tb의 ID 1번 회원)
-- INSERT INTO memberqr_tb (mq_user, mq_code, mq_isactive, mq_generateddate)
-- VALUES (1, UUID(), 1, NOW());

-- -- 공지사항 샘플
-- INSERT INTO notice_tb (nt_gym, nt_title, nt_content, nt_type, nt_status, nt_createddate)
-- VALUES
-- (1, '신규 오픈 이벤트', '헬스장 오픈 기념 첫 달 회원권 50% 할인!', 2, 1, NOW()),
-- (1, '운영 시간 안내', '평일: 06:00-23:00, 주말: 08:00-20:00', 0, 1, NOW());

-- -- 앱 버전 샘플
-- INSERT INTO appversion_tb (av_platform, av_version, av_minversion, av_forceupdate, av_status, av_releasedate)
-- VALUES
-- ('android', '1.0.0', '1.0.0', 0, 1, NOW()),
-- ('ios', '1.0.0', '1.0.0', 0, 1, NOW());
