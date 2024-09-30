#
# TABLE STRUCTURE FOR: Retreat
#

DROP TABLE IF EXISTS `Retreat`;

CREATE TABLE `Retreat` (
  `id` int(9) unsigned,
  `Date` date DEFAULT NULL,
  `Role` varchar(30) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

#
# TABLE STRUCTURE FOR: auditlog
#

DROP TABLE IF EXISTS `auditlog`;

CREATE TABLE `auditlog` (
  `Changer` int(9) unsigned NOT NULL,
  `Changee` int(9) unsigned NOT NULL,
  `Field` varchar(32),
  `Value` varchar(100),
  `Date` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

#
# TABLE STRUCTURE FOR: member
#

DROP TABLE IF EXISTS `member`;

CREATE TABLE `member` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `MemberStatus` enum('Annual Vows','Life Vows','Friend','Benefactor','Removed') NOT NULL,
  `FirstName` varchar(50) NOT NULL,
  `MiddleName` varchar(50),
  `LastName` varchar(50) NOT NULL,
  `PreferredName` varchar(50),
  `Title` enum('Sr.','Br.','Sibling'),
  `LifevowName` varchar(50),
  `Suffix` varchar(10),
  `Address` varchar(100),
  `AddressLine2` varchar(100),
  `City` varchar(50),
  `State` varchar(30),
  `Country` varchar(20),
  `PostalCode` varchar(20),
  `PrimaryPhone` char(20),
  `SecondaryPhone` varchar(20),
  `PrimaryEmail` varchar(50),
  `SecondaryEmail` varchar(50),
  `BirthDate` date,
  `DateRecordCreated` date NOT NULL,
  `Chapter` varchar(50),
  `DateFirstVows` date,
  `DateReaffirmation` date,
  `DateRemoved` date,
  `DateFirstProfession` date,
  `DateDeceased` date,
  `DateNovitiate` date,
  `Status` enum('laity','clergy','student','retired laity','retired clergy'),
  `HowJoined` varchar(50),
  `HowRemoved` varchar(50),
  `ListInDirectory` tinyint(1) unsigned DEFAULT '0',
  `ListAddress` tinyint(1) unsigned DEFAULT '0',
  `ListPrimaryPhone` tinyint(1) unsigned DEFAULT '0',
  `ListSecondaryPhone` tinyint(1) unsigned DEFAULT '0',
  `ListPrimaryEmail` tinyint(1) unsigned DEFAULT '0',
  `ListSecondaryEmail` tinyint(1) unsigned DEFAULT '0',
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  `Occupation` varchar(50),
  `Employer` varchar(50),
  `Denomination` varchar(50),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

DROP TABLE IF EXISTS `subscriber`;

CREATE TABLE `subscriber` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) NOT NULL,
  `Attn` varchar(100),
  `Address` varchar(100),
  `AddressLine2` varchar(100),
  `City` varchar(100),
  `State` varchar(20),
  `Country` varchar(20),
  `PostalCode` varchar(20),
  `PrimaryPhone` char(20),
  `SecondaryPhone` varchar(20),
  `PrimaryEmail` varchar(50),
  `SecondaryEmail` varchar(50),
  `DateRecordCreated` date NOT NULL,
  `DatePaid` date NOT NULL,
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

DROP TABLE IF EXISTS `auth`;

CREATE TABLE `auth` (
  `user` varchar(100) NOT NULL,
  `pwhash` varchar(100) NOT NULL,
  `level` tinyint(2) NOT NULL,
  PRIMARY KEY (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

