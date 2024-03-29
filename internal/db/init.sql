CREATE DATABASE IF NOT EXISTS db_user;

USE db_user;

CREATE TABLE IF NOT EXISTS users (
  user_id CHAR(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  name VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT "",
  age TINYINT(3) UNSIGNED NOT NULL DEFAULT 0,  -- Use UNSIGNED for non-negative ages
  PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;