/*
SQLyog Ultimate v12.3.3 (64 bit)
MySQL - 5.7.22 : Database - stock
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`stock` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `stock`;

/*Table structure for table `stock_trends` */

DROP TABLE IF EXISTS `stock_trends`;

CREATE TABLE `stock_trends` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT '' COMMENT '股票代码',
  `name` varchar(10) DEFAULT '' COMMENT '股票名称',
  `open_price` float DEFAULT '0' COMMENT '开盘价',
  `close_price` float DEFAULT '0' COMMENT '收盘价',
  `open_percent` float DEFAULT '0' COMMENT '开盘幅度',
  `close_percent` float DEFAULT '0' COMMENT '收盘幅度',
  `high_percent` float DEFAULT '0' COMMENT '最高幅度',
  `low_percent` float DEFAULT '0' COMMENT '最低幅度',
  `shock` float DEFAULT '0' COMMENT '当天最大振幅',
  `amount` bigint(20) DEFAULT '0' COMMENT '成交额',
  `amount_format` varchar(20) DEFAULT '' COMMENT '成交额(格式化)',
  `close_color` varchar(4) DEFAULT '' COMMENT '1收红，2收绿',
  `date` varchar(30) DEFAULT '',
  `created_at` int(11) DEFAULT '0',
  `updated_at` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=585 DEFAULT CHARSET=utf8;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
