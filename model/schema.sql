#
# TABLE STRUCTURE FOR: Retreat
#

DROP TABLE IF EXISTS `Retreat`;

CREATE TABLE `Retreat` (
  `id` int(9) unsigned,
  `Date` date DEFAULT NULL,
  `Role` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

#
# TABLE STRUCTURE FOR: auditlog
#

DROP TABLE IF EXISTS `auditlog`;

CREATE TABLE `auditlog` (
  `Changer` int(9) NOT NULL,
  `Changee` int(9) NOT NULL,
  `Date` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

#
# TABLE STRUCTURE FOR: member
#

DROP TABLE IF EXISTS `member`;

CREATE TABLE `member` (
  `id` bigint(16) unsigned NOT NULL AUTO_INCREMENT,
  `MemberStatus` enum('Annual Vows','Life Vows','Contributor','Removed') NOT NULL,
  `FirstName` varchar(100) NOT NULL,
  `MiddleName` varchar(100),
  `LastName` varchar(100) NOT NULL,
  `PreferredName` varchar(100),
  `Title` enum('Sr.','Br.','Sibling'),
  `LifevowName` varchar(255),
  `Suffix` varchar(20),
  `Address` varchar(255),
  `AddressLine2` varchar(255),
  `City` varchar(100),
  `State` varchar(100),
  `Country` varchar(100),
  `PostalCode` varchar(20),
  `PrimaryPhone` char(20),
  `SecondaryPhone` varchar(20),
  `PrimaryEmail` varchar(255),
  `SecondaryEmail` varchar(255),
  `BirthDate` date,
  `DateRecordCreated` date NOT NULL,
  `Chapter` varchar(255),
  `DateFirstVows` date,
  `DateReaffirmation` date,
  `DateRemoved` date,
  `DateFirstProfession` date,
  `DateDeceased` date,
  `DateNovitiate` date,
  `Status` enum('laity','clergy','student','retired laity','retired clergy'),
  `HowJoined` varchar(255),
  `HowRemoved` varchar(255),
  `ListInDirectory` enum('YES','NO') NOT NULL,
  `ListAddress` enum('YES','NO') NOT NULL,
  `ListPrimaryPhone` enum('YES','NO') NOT NULL,
  `ListSecondaryPhone` enum('YES','NO') NOT NULL,
  `ListPrimaryEmail` enum('YES','NO') NOT NULL,
  `ListSecondaryEmail` enum('YES','NO') NOT NULL,
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  `Occupation` varchar(255),
  `Employeer` varchar(255),
  `Denomination` varchar(255),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

DROP TABLE IF EXISTS `subscriber`;

CREATE TABLE `subscriber` (
  `id` bigint(16) unsigned NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) NOT NULL,
  `Attn` varchar(100),
  `Address` varchar(255),
  `AddressLine2` varchar(255),
  `City` varchar(100),
  `State` varchar(100),
  `Country` varchar(100),
  `PostalCode` varchar(20),
  `PrimaryPhone` char(20),
  `SecondaryPhone` varchar(20),
  `PrimaryEmail` varchar(255),
  `SecondaryEmail` varchar(255),
  `DateRecordCreated` date NOT NULL,
  `DatePaid` date NOT NULL,
  `Doxology` enum('electronic','mailed','none') NOT NULL,
  `Newsletter` enum('electronic','mailed','none') NOT NULL,
  `Communication` enum('electronic','mailed','none') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
