create table if not exists `users` (`Id` varchar(128) not null primary key, `Created` bigint, `LastSeen` bigint, `DisplayName` varchar(64), `PhotoUrl` text, `Email` varchar(128))  engine=InnoDB charset=utf8;
create table if not exists `games` (`Name` varchar(64), `Id` varchar(16) not null primary key, `Version` bigint, `Winners` varchar(128), `Finished` boolean, `NumPlayers` bigint, `Agents` text)  engine=InnoDB charset=utf8;
create table if not exists `states` (`Id` bigint not null primary key auto_increment, `GameId` varchar(16), `Version` bigint, `Blob` text)  engine=InnoDB charset=utf8;
create table if not exists `cookies` (`Cookie` varchar(64) not null primary key, `UserId` varchar(128))  engine=InnoDB charset=utf8;
create table if not exists `players` (`Id` bigint not null primary key auto_increment, `GameId` varchar(16), `PlayerIndex` bigint, `UserId` varchar(128))  engine=InnoDB charset=utf8;
create table if not exists `agentstates` (`Id` bigint not null primary key auto_increment, `GameId` varchar(16), `PlayerIndex` bigint, `Blob` text)  engine=InnoDB charset=utf8;
