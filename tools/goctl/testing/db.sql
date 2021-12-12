CREATE DATABASE IF NOT EXISTS user;

USE user;

CREATE TABLE `user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `age` tinyint(4) unsigned NOT NULL DEFAULT 0 COMMENT '年龄',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名字',
    `addr` varchar(128) NOT NULL DEFAULT '' COMMENT '地址',
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_mid` (`name`),
    KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;