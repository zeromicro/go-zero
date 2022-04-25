CREATE TABLE `test_user`
(
    `id`          bigint                                                 NOT NULL AUTO_INCREMENT,
    `mobile`      varchar(255) COLLATE utf8mb4_bin                       NOT NULL,
    `class`       bigint                                                 NOT NULL,
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `mobile_unique` (`mobile`),
    UNIQUE KEY `class_name_unique` (`class`,`name`),
    KEY           `create_index` (`create_time`),
    KEY           `name_index` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;