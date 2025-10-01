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
  `al_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `al_title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `al_content` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `al_type` int(11) NOT NULL,
  `al_status` int(11) NOT NULL,
  `al_user` bigint(20) NOT NULL,
  `al_date` datetime NOT NULL,
  PRIMARY KEY (`al_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `daytype_tb`
--

DROP TABLE IF EXISTS `daytype_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `daytype_tb` (
  `dt_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `dt_gym` bigint(20) DEFAULT NULL,
  `dt_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `dt_date` datetime DEFAULT NULL,
  PRIMARY KEY (`dt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `discount_tb`
--

DROP TABLE IF EXISTS `discount_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `discount_tb` (
  `d_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `d_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `d_discount` int(11) DEFAULT NULL,
  `d_date` datetime DEFAULT NULL,
  PRIMARY KEY (`d_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gym_tb`
--

DROP TABLE IF EXISTS `gym_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gym_tb` (
  `g_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `g_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `g_date` datetime DEFAULT NULL,
  PRIMARY KEY (`g_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `health_tb`
--

DROP TABLE IF EXISTS `health_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `health_tb` (
  `h_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `h_category` bigint(20) DEFAULT NULL,
  `h_term` bigint(20) DEFAULT NULL,
  `h_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `h_count` int(11) DEFAULT NULL,
  `h_cost` int(11) DEFAULT NULL,
  `h_discount` bigint(20) DEFAULT NULL,
  `h_costdiscount` int(11) DEFAULT NULL,
  `h_content` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `h_date` datetime DEFAULT NULL,
  PRIMARY KEY (`h_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `healthcategory_tb`
--

DROP TABLE IF EXISTS `healthcategory_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `healthcategory_tb` (
  `hc_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `hc_gym` bigint(20) DEFAULT NULL,
  `hc_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `hc_date` datetime DEFAULT NULL,
  PRIMARY KEY (`hc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ipblock_tb`
--

DROP TABLE IF EXISTS `ipblock_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipblock_tb` (
  `ib_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ib_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ib_type` int(11) NOT NULL,
  `ib_policy` int(11) NOT NULL,
  `ib_use` int(11) NOT NULL,
  `ib_order` int(11) NOT NULL,
  `ib_date` datetime NOT NULL,
  PRIMARY KEY (`ib_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `loginlog_tb`
--

DROP TABLE IF EXISTS `loginlog_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `loginlog_tb` (
  `ll_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ll_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ll_ipvalue` bigint(20) NOT NULL,
  `ll_user` bigint(20) NOT NULL,
  `ll_date` datetime NOT NULL,
  PRIMARY KEY (`ll_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `membership_tb`
--

DROP TABLE IF EXISTS `membership_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `membership_tb` (
  `m_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `m_gym` bigint(20) DEFAULT NULL,
  `m_user` bigint(20) DEFAULT NULL,
  `m_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `m_sex` int(11) DEFAULT NULL,
  `m_birth` datetime DEFAULT NULL,
  `m_phonenum` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `m_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `m_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `m_date` datetime DEFAULT NULL,
  PRIMARY KEY (`m_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_tb`
--

DROP TABLE IF EXISTS `order_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_tb` (
  `o_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `o_membership` bigint(20) DEFAULT NULL,
  `o_date` datetime DEFAULT NULL,
  PRIMARY KEY (`o_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payment_tb`
--

DROP TABLE IF EXISTS `payment_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_tb` (
  `p_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `p_gym` bigint(20) DEFAULT NULL,
  `p_order` bigint(20) DEFAULT NULL,
  `p_membership` bigint(20) DEFAULT NULL,
  `p_cost` int(11) DEFAULT NULL,
  `p_date` datetime DEFAULT NULL,
  PRIMARY KEY (`p_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `paymentform_tb`
--

DROP TABLE IF EXISTS `paymentform_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `paymentform_tb` (
  `pf_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pf_gym` bigint(20) DEFAULT NULL,
  `pf_payment` bigint(20) DEFAULT NULL,
  `pf_type` bigint(20) DEFAULT NULL,
  `pf_cost` int(11) DEFAULT NULL,
  `pf_date` datetime DEFAULT NULL,
  PRIMARY KEY (`pf_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `paymenttype_tb`
--

DROP TABLE IF EXISTS `paymenttype_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `paymenttype_tb` (
  `pt_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pt_gym` bigint(20) DEFAULT NULL,
  `pt_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pt_date` datetime DEFAULT NULL,
  PRIMARY KEY (`pt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rocker_tb`
--

DROP TABLE IF EXISTS `rocker_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rocker_tb` (
  `r_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `r_group` bigint(20) DEFAULT NULL,
  `r_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `r_available` int(11) DEFAULT NULL,
  `r_date` datetime DEFAULT NULL,
  PRIMARY KEY (`r_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rockergroup_tb`
--

DROP TABLE IF EXISTS `rockergroup_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rockergroup_tb` (
  `rg_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rg_gym` bigint(20) DEFAULT NULL,
  `rg_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `rg_date` datetime DEFAULT NULL,
  PRIMARY KEY (`rg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role_tb`
--

DROP TABLE IF EXISTS `role_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_tb` (
  `r_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `r_gym` bigint(20) DEFAULT NULL,
  `r_role` int(11) DEFAULT NULL,
  `r_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `r_date` datetime DEFAULT NULL,
  PRIMARY KEY (`r_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `setting_tb`
--

DROP TABLE IF EXISTS `setting_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `setting_tb` (
  `se_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `se_category` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `se_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `se_key` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `se_value` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `se_remark` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `se_type` int(11) NOT NULL,
  `se_data` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `se_order` int(11) NOT NULL,
  `se_date` datetime NOT NULL,
  PRIMARY KEY (`se_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stop_tb`
--

DROP TABLE IF EXISTS `stop_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stop_tb` (
  `s_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `s_usehelth` bigint(20) DEFAULT NULL,
  `s_startday` datetime DEFAULT NULL,
  `s_endday` datetime DEFAULT NULL,
  `s_count` int(11) DEFAULT NULL,
  `s_date` datetime DEFAULT NULL,
  PRIMARY KEY (`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `systemlog_tb`
--

DROP TABLE IF EXISTS `systemlog_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `systemlog_tb` (
  `sl_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sl_type` int(11) NOT NULL,
  `sl_content` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sl_result` int(11) NOT NULL,
  `sl_date` datetime NOT NULL,
  PRIMARY KEY (`sl_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_tb`
--

DROP TABLE IF EXISTS `term_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_tb` (
  `t_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `t_gym` bigint(20) DEFAULT NULL,
  `t_daytype` bigint(20) DEFAULT NULL,
  `t_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `t_term` int(11) DEFAULT NULL,
  `t_date` datetime DEFAULT NULL,
  PRIMARY KEY (`t_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `token_tb`
--

DROP TABLE IF EXISTS `token_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `token_tb` (
  `to_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `to_user` bigint(20) NOT NULL,
  `to_token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `to_status` int(11) NOT NULL,
  `to_date` datetime NOT NULL,
  PRIMARY KEY (`to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `usehealth_tb`
--

DROP TABLE IF EXISTS `usehealth_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `usehealth_tb` (
  `uh_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uh_order` bigint(20) DEFAULT NULL,
  `uh_health` bigint(20) DEFAULT NULL,
  `uh_user` bigint(20) DEFAULT NULL,
  `uh_rocker` bigint(20) DEFAULT NULL,
  `uh_term` bigint(20) DEFAULT NULL,
  `uh_discount` bigint(20) DEFAULT NULL,
  `uh_startday` datetime DEFAULT NULL,
  `uh_endday` datetime DEFAULT NULL,
  `uh_date` datetime DEFAULT NULL,
  PRIMARY KEY (`uh_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_tb`
--

DROP TABLE IF EXISTS `user_tb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_tb` (
  `u_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `u_loginid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_passwd` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_tel` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_sex` int(11) DEFAULT NULL,
  `u_birth` datetime DEFAULT NULL,
  `u_type` int(11) NOT NULL,
  `u_connectid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `u_level` int(11) NOT NULL,
  `u_role` int(11) DEFAULT NULL,
  `u_use` int(11) DEFAULT NULL,
  `u_logindate` datetime DEFAULT NULL,
  `u_lastchangepasswddate` datetime DEFAULT NULL,
  `u_date` datetime DEFAULT NULL,
  PRIMARY KEY (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
