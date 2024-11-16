-- MySQL dump 10.19  Distrib 10.3.28-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: oo
-- ------------------------------------------------------
-- Server version	10.3.28-MariaDB-log

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
-- Table structure for table `auditlog`
--

DROP TABLE IF EXISTS `auditlog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auditlog` (
  `Changer` int(9) unsigned NOT NULL,
  `Changee` int(9) unsigned NOT NULL,
  `Field` varchar(32) DEFAULT NULL,
  `Value` varchar(100) DEFAULT NULL,
  `Date` date NOT NULL,
  KEY `changee_ind` (`Changee`),
  KEY `Changer` (`Changer`),
  CONSTRAINT `auditlog_ibfk_1` FOREIGN KEY (`Changee`) REFERENCES `member` (`id`) ON DELETE CASCADE,
  CONSTRAINT `auditlog_ibfk_2` FOREIGN KEY (`Changer`) REFERENCES `member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `auth`
--

DROP TABLE IF EXISTS `auth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth` (
  `user` varchar(50) NOT NULL,
  `pwhash` varchar(61) NOT NULL,
  `level` tinyint(2) NOT NULL,
  PRIMARY KEY (`user`),
  CONSTRAINT `auth_ibfk_1` FOREIGN KEY (`user`) REFERENCES `member` (`PrimaryEmail`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chaptermembers`
--

DROP TABLE IF EXISTS `chaptermembers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chaptermembers` (
  `member` int(9) unsigned DEFAULT NULL,
  `chapter` tinyint(2) NOT NULL,
  KEY `member_ind` (`member`),
  KEY `chapter_ind` (`chapter`),
  CONSTRAINT `chaptermembers_ibfk_1` FOREIGN KEY (`chapter`) REFERENCES `chapters` (`id`) ON DELETE CASCADE,
  CONSTRAINT `chaptermembers_ibfk_2` FOREIGN KEY (`member`) REFERENCES `member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chapters`
--

DROP TABLE IF EXISTS `chapters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chapters` (
  `id` tinyint(2) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `prior` int(9) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `giving`
--

DROP TABLE IF EXISTS `giving`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `giving` (
  `entryID` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `id` int(9) unsigned NOT NULL,
  `amount` decimal(13,2) NOT NULL,
  `check` int(9) DEFAULT NULL,
  `transaction` varchar(24) DEFAULT NULL,
  `description` varchar(32) NOT NULL,
  `date` date NOT NULL DEFAULT curdate(),
  PRIMARY KEY (`entryID`),
  KEY `id` (`id`),
  CONSTRAINT `giving_ibfk_1` FOREIGN KEY (`id`) REFERENCES `member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=331 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `member`
--

DROP TABLE IF EXISTS `member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `member` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `MemberStatus` enum('Annual Vows','Life Vows','Friend','Deceased','Removed') NOT NULL DEFAULT 'Friend',
  `FirstName` varchar(50) NOT NULL,
  `MiddleName` varchar(50) DEFAULT NULL,
  `LastName` varchar(50) NOT NULL,
  `PreferredName` varchar(50) DEFAULT NULL,
  `Title` enum('Sr.','Br.','Sibling') DEFAULT NULL,
  `LifevowName` varchar(50) DEFAULT NULL,
  `Suffix` varchar(10) DEFAULT NULL,
  `Address` varchar(100) DEFAULT NULL,
  `AddressLine2` varchar(100) DEFAULT NULL,
  `City` varchar(50) DEFAULT NULL,
  `State` varchar(30) DEFAULT NULL,
  `Country` char(2) DEFAULT NULL,
  `PostalCode` varchar(20) DEFAULT NULL,
  `PrimaryPhone` char(20) DEFAULT NULL,
  `SecondaryPhone` varchar(20) DEFAULT NULL,
  `PrimaryEmail` varchar(50) DEFAULT NULL,
  `SecondaryEmail` varchar(50) DEFAULT NULL,
  `BirthDate` date DEFAULT NULL,
  `DateRecordCreated` date NOT NULL,
  `Chapter` varchar(50) DEFAULT NULL,
  `DateFirstVows` date DEFAULT NULL,
  `DateReaffirmation` date DEFAULT NULL,
  `DateRemoved` date DEFAULT NULL,
  `DateFirstProfession` date DEFAULT NULL,
  `DateDeceased` date DEFAULT NULL,
  `DateNovitiate` date DEFAULT NULL,
  `DateLifeVows` date DEFAULT NULL,
  `Status` enum('laity','clergy','student','retired laity','retired clergy') DEFAULT NULL,
  `Leadership` enum('member','prior','council','canon','elected') NOT NULL DEFAULT 'member',
  `HowJoined` varchar(50) DEFAULT NULL,
  `HowRemoved` varchar(50) DEFAULT NULL,
  `ListInDirectory` tinyint(1) unsigned DEFAULT 0,
  `ListAddress` tinyint(1) unsigned DEFAULT 0,
  `ListPrimaryPhone` tinyint(1) unsigned DEFAULT 0,
  `ListSecondaryPhone` tinyint(1) unsigned DEFAULT 0,
  `ListPrimaryEmail` tinyint(1) unsigned DEFAULT 0,
  `ListSecondaryEmail` tinyint(1) unsigned DEFAULT 0,
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  `Occupation` varchar(50) DEFAULT NULL,
  `Employer` varchar(50) DEFAULT NULL,
  `Denomination` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `email_ind` (`PrimaryEmail`)
) ENGINE=InnoDB AUTO_INCREMENT=1877 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `subscriber`
--

DROP TABLE IF EXISTS `subscriber`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `subscriber` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) NOT NULL,
  `Attn` varchar(100) DEFAULT NULL,
  `Address` varchar(100) DEFAULT NULL,
  `AddressLine2` varchar(100) DEFAULT NULL,
  `City` varchar(100) DEFAULT NULL,
  `State` varchar(20) DEFAULT NULL,
  `Country` char(2) DEFAULT NULL,
  `PostalCode` varchar(20) DEFAULT NULL,
  `PrimaryPhone` char(20) DEFAULT NULL,
  `SecondaryPhone` varchar(20) DEFAULT NULL,
  `PrimaryEmail` varchar(50) DEFAULT NULL,
  `SecondaryEmail` varchar(50) DEFAULT NULL,
  `DateRecordCreated` date NOT NULL,
  `DatePaid` date NOT NULL,
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1820 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-11-16 18:55:20
