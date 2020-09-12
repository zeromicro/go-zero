CREATE TABLE `book`
(
  `book` varchar(255) NOT NULL COMMENT 'book name',
  `price` int NOT NULL COMMENT 'book price',
  PRIMARY KEY(`book`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
