-- 用户系统数据库建表脚本 (MySQL 8.0)
-- 对应接口：注册、登录、退出（Session 存 Redis）

-- CREATE DATABASE IF NOT EXISTS user_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE user_system;

CREATE TABLE IF NOT EXISTS `users` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username`      VARCHAR(64)  NOT NULL COMMENT '用户名，唯一',
  `password_hash` VARCHAR(255) NOT NULL COMMENT 'bcrypt 密码哈希',
  `email`         VARCHAR(128) DEFAULT '' COMMENT '邮箱，可选',
  `is_login`      TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否已登录',
  `login_at`      DATETIME(3)  NULL DEFAULT NULL COMMENT '最近登录时间',
  `created_at`    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at`    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`) COMMENT '用户名唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
