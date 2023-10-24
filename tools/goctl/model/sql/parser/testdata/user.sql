CREATE TABLE `test_user`
(
    `id`          bigint                                                 NOT NULL AUTO_INCREMENT,
    `mobile`      varchar(255) COLLATE utf8mb4_bin                       NOT NULL comment '手\t机  号',
    `class`       bigint                                                 NOT NULL comment '班级',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL comment '姓
  名',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP comment '创建\r时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` timestamp NULL DEFAULT NULL,
    `delete_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `mobile_unique` (`mobile`),
    UNIQUE KEY `class_name_unique` (`class`,`name`),
    KEY           `create_index` (`create_time`),
    KEY           `name_index` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;