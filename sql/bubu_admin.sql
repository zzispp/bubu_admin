/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80033 (8.0.33)
 Source Host           : localhost:3306
 Source Schema         : bubu_admin

 Target Server Type    : MySQL
 Target Server Version : 80033 (8.0.33)
 File Encoding         : 65001

 Date: 17/10/2024 21:35:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
BEGIN;
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (102, 'g', '901673b9-8994-4a9f-bd80-afd5ad5d3631', 'cs535sp0ndhvrs7lsiqg', '', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (88, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs534eh0ndhvq18295mg', 'menu', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (87, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs534k10ndhvq18295n0', 'menu', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (89, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs58r110ndhg4fpkfp7g', 'menu', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (90, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs58ro10ndhg4fpkfp80', 'menu', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (91, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs6iqb90ndhnceoh14jg', 'menu', '', '', '');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (92, 'p', 'cs535sp0ndhvrs7lsiqg', 'cs77rvp0ndhh7fe99n2g', 'menu', '', '', '');
COMMIT;

-- ----------------------------
-- Table structure for menus
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
  `id` varchar(20) NOT NULL,
  `code` varchar(32) DEFAULT NULL,
  `name` varchar(128) DEFAULT NULL,
  `description` varchar(1024) DEFAULT NULL,
  `sequence` int DEFAULT NULL,
  `type` varchar(20) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `redirect` varchar(255) DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  `parent_id` varchar(20) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `component` varchar(255) DEFAULT NULL,
  `icon` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_menus_code` (`code`),
  KEY `idx_menus_name` (`name`),
  KEY `idx_menus_sequence` (`sequence`),
  KEY `idx_menus_type` (`type`),
  KEY `idx_menus_status` (`status`),
  KEY `idx_menus_parent_id` (`parent_id`),
  KEY `idx_menus_created_at` (`created_at`),
  KEY `idx_menus_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of menus
-- ----------------------------
BEGIN;
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs534eh0ndhvq18295mg', 'dashboard', '仪表盘', '仪表盘的菜单项', 1, 'nav', '/dashboard', '/dashboard/overview', 'enable', 'root', '2024-10-12 16:24:26.231', '2024-10-13 23:38:21.917', 'Layout', 'dashboard');
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs534k10ndhvq18295n0', 'overview', '概览', '', 1, 'nav', '/overview', '', 'enable', 'cs534eh0ndhvq18295mg', '2024-10-12 16:24:48.858', '2024-10-12 16:24:48.858', 'pages/dashboard/overview', 'squareChartGantt');
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs58r110ndhg4fpkfp7g', 'setting', '系统设置', '', 2, 'nav', '/setting', '', 'enable', 'root', '2024-10-12 22:53:56.024', '2024-10-12 22:53:56.024', 'Layout', 'settings');
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs58ro10ndhg4fpkfp80', 'menuManager', '菜单管理', '菜单管理', 1, 'nav', '/menuManager', '', 'enable', 'cs58r110ndhg4fpkfp7g', '2024-10-12 22:55:28.387', '2024-10-12 22:55:28.387', 'pages/menuManager', 'menu');
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs6iqb90ndhnceoh14jg', 'roleManager', '角色管理', '角色管理', 1, 'nav', '/roleManger', '', 'enable', 'cs58r110ndhg4fpkfp7g', '2024-10-14 22:39:41.400', '2024-10-14 22:39:41.400', 'pages/roleManager', 'menu');
INSERT INTO `menus` (`id`, `code`, `name`, `description`, `sequence`, `type`, `path`, `redirect`, `status`, `parent_id`, `created_at`, `updated_at`, `component`, `icon`) VALUES ('cs77rvp0ndhh7fe99n2g', 'userManager', '用户管理', '用户管理', 1, 'nav', '/userManager', '', 'enable', 'cs58r110ndhg4fpkfp7g', '2024-10-15 22:36:47.348', '2024-10-15 22:36:53.709', 'pages/userManager', 'menu');
COMMIT;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` varchar(20) NOT NULL,
  `code` varchar(32) DEFAULT NULL,
  `name` varchar(128) DEFAULT NULL,
  `description` varchar(1024) DEFAULT NULL,
  `sequence` int DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_roles_code` (`code`),
  KEY `idx_roles_name` (`name`),
  KEY `idx_roles_sequence` (`sequence`),
  KEY `idx_roles_status` (`status`),
  KEY `idx_roles_created_at` (`created_at`),
  KEY `idx_roles_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of roles
-- ----------------------------
BEGIN;
INSERT INTO `roles` (`id`, `code`, `name`, `description`, `sequence`, `status`, `created_at`, `updated_at`) VALUES ('cs535sp0ndhvrs7lsiqg', 'merchant', '商户', '商户', 0, 'enable', '2024-10-12 16:27:31.722', '2024-10-15 22:58:18.500');
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(191) NOT NULL,
  `name` longtext,
  `password` longtext NOT NULL,
  `email` longtext NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`id`, `name`, `password`, `email`, `created_at`, `updated_at`, `deleted_at`) VALUES ('901673b9-8994-4a9f-bd80-afd5ad5d3631', 'Zwj', '$2a$10$egPd/87fPQqN28MvVAsf0ebPfMWIYAmOEMB7CRQUTkS6DEvLxRPVy', 'a@gmail.com', '2024-10-10 12:33:26.191', '2024-10-15 23:39:06.203', NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
