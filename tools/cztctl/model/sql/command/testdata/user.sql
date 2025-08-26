-- 用户表 --
CREATE TABLE `user`
(
    `id`          bigint(10) NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',
    `password`    varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',
    `mobile`      varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
    `gender`      char(5) COLLATE utf8mb4_general_ci      NOT NULL COMMENT '男｜女｜未公开',
    `nickname`    varchar(255) COLLATE utf8mb4_general_ci          DEFAULT '' COMMENT '用户昵称',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name_index` (`name`),
    UNIQUE KEY `mobile_index` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

