CREATE TABLE `shorturl`
(
  `shorten` varchar(255) NOT NULL COMMENT 'shorten key',
  `url` varchar(255) NOT NULL COMMENT 'original url',
  PRIMARY KEY(`shorten`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
