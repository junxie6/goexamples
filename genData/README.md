```
CREATE TABLE `Project` (
  `ProjectID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL STORAGE DISK,
  `Status` int(10) unsigned NOT NULL,
  PRIMARY KEY (`ProjectID`),
  KEY `Status` (`Status`)
) TABLESPACE ts_1 STORAGE DISK ENGINE=NDBCLUSTER DEFAULT CHARSET=utf8mb4;
```
```
CREATE TABLE `Bug` (
  `BugID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `ProjectID` int(10) unsigned NOT NULL,
  `Status` int(10) unsigned NOT NULL,
  `Amount` double unsigned NOT NULL STORAGE DISK,
  `Qty` int(10) unsigned NOT NULL STORAGE DISK,
  `Summary` varchar(255) NOT NULL STORAGE DISK,
  `Created` datetime NOT NULL,
  `Year` tinyint(3) unsigned NOT NULL,
  PRIMARY KEY (`BugID`),
  KEY `ProjectID` (`ProjectID`),
  KEY `Status` (`Status`),
  KEY `Created` (`Created`),
  KEY `Year` (`Year`),
  KEY `ProjectID_Status_Created` (`ProjectID`,`Status`,`Created`)
) TABLESPACE ts_1 STORAGE DISK ENGINE=NDBCLUSTER DEFAULT CHARSET=utf8mb4;
```
```
CREATE TABLE `BugBody` (
  `BugID` int(10) unsigned NOT NULL,
  `Body` text NOT NULL STORAGE DISK,
  PRIMARY KEY (`BugID`)
) TABLESPACE ts_1 STORAGE DISK ENGINE=NDBCLUSTER DEFAULT CHARSET=utf8mb4;
```
