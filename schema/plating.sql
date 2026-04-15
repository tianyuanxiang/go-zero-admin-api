/*
 Navicat Premium Data Transfer

 Source Server         : mysql-5.7
 Source Server Type    : MySQL
 Source Server Version : 50700
 Source Host           : 172.16.90.70:3307
 Source Schema         : plating

 Target Server Type    : MySQL
 Target Server Version : 50700
 File Encoding         : 65001

 Date: 10/04/2026 13:47:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_casbin_rule`(`ptype` ASC, `v0` ASC, `v1` ASC, `v2` ASC, `v3` ASC, `v4` ASC, `v5` ASC) USING BTREE,
  UNIQUE INDEX `idx_casbin_rule`(`ptype` ASC, `v0` ASC, `v1` ASC, `v2` ASC, `v3` ASC, `v4` ASC, `v5` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Casbin权限规则表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');

-- ----------------------------
-- Table structure for plate_dosing_events
-- ----------------------------
DROP TABLE IF EXISTS `plate_dosing_events`;
CREATE TABLE `plate_dosing_events`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联槽体ID',
  `event_time` datetime NOT NULL COMMENT '加药时间',
  `chemical` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '化学品名称，如：CrO3、H2SO4',
  `amount_grams` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '加入量（克）',
  `processed` tinyint NOT NULL DEFAULT 0 COMMENT '是否已被模型处理：1=已处理，0=未处理',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tank_event_time`(`tank_id` ASC, `event_time` ASC) USING BTREE,
  INDEX `idx_chemical`(`chemical` ASC) USING BTREE,
  INDEX `idx_processed`(`processed` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '加药事件表，记录每次向槽液中添加化学品的数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plate_dosing_events
-- ----------------------------

-- ----------------------------
-- Table structure for plate_model_state
-- ----------------------------
DROP TABLE IF EXISTS `plate_model_state`;
CREATE TABLE `plate_model_state`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `tank_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联槽体ID',
  `calc_time` datetime NOT NULL COMMENT '计算/记录时间',
  `cro3_g_per_l` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '铬酸浓度（g/L，CrO3）',
  `cr3_g_per_l` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '三价铬浓度（g/L，Cr3+）',
  `tank_volume` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '当前槽液体积（升）',
  `source` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'calc' COMMENT '数据来源：init=初始化，calc=模型计算，lab_override=化验覆盖',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tank_calc_time`(`tank_id` ASC, `calc_time` ASC) USING BTREE,
  INDEX `idx_calc_time`(`calc_time` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '模型状态表，记录槽液浓度的历史计算状态' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plate_model_state
-- ----------------------------

-- ----------------------------
-- Table structure for plate_production_events
-- ----------------------------
DROP TABLE IF EXISTS `plate_production_events`;
CREATE TABLE `plate_production_events`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联槽体ID',
  `event_time` datetime NOT NULL COMMENT '事件发生时间',
  `batch_start_time` datetime NOT NULL COMMENT '批次开始时间',
  `part_count` int NOT NULL DEFAULT 0 COMMENT '本批次工件数量（件）',
  `part_area_dm2` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '工件总面积（dm2）',
  `carryover_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'A' COMMENT '带出类型：A/B/C（对应carryover_a/b/c）',
  `ah_this_batch` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '本批次实际Ah数（安培时）',
  `ah_calculated` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '模型计算使用的Ah（可能来自TDEngine或手工录入）',
  `processed` tinyint NOT NULL DEFAULT 0 COMMENT '是否已被模型计算处理：1=已处理，0=未处理',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tank_event_time`(`tank_id` ASC, `event_time` ASC) USING BTREE,
  INDEX `idx_processed`(`processed` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '生产事件表，记录每批次镀铬生产数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plate_production_events
-- ----------------------------

-- ----------------------------
-- Table structure for plate_tank_config
-- ----------------------------
DROP TABLE IF EXISTS `plate_tank_config`;
CREATE TABLE `plate_tank_config`  (
  `tank_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '槽体ID（主键，如：TANK-001）',
  `tank_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '槽体名称',
  `volume_liters` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '槽体容积（升）',
  `density_g_per_l` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '镀液密度（g/L）',
  `d_cro3_g_per_ah` decimal(10, 6) NOT NULL DEFAULT 0.000000 COMMENT 'CrO3消耗系数（g/Ah），法拉第系数',
  `d_cr3_g_per_ah` decimal(10, 6) NOT NULL DEFAULT 0.000000 COMMENT 'Cr3+生成系数（g/Ah）',
  `carryover_a` decimal(10, 6) NOT NULL DEFAULT 0.000000 COMMENT '带出系数A（随工件面积）',
  `carryover_b` decimal(10, 6) NOT NULL DEFAULT 0.000000 COMMENT '带出系数B（随工件数量）',
  `carryover_c` decimal(10, 6) NOT NULL DEFAULT 0.000000 COMMENT '带出系数C（固定带出）',
  `is_active` tinyint NOT NULL DEFAULT 1 COMMENT '是否激活：1=激活，0=停用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`tank_id`) USING BTREE,
  INDEX `idx_is_active`(`is_active` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '槽体配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plate_tank_config
-- ----------------------------

-- ----------------------------
-- Table structure for plate_water_events
-- ----------------------------
DROP TABLE IF EXISTS `plate_water_events`;
CREATE TABLE `plate_water_events`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联槽体ID',
  `event_time` datetime NOT NULL COMMENT '补水时间',
  `water_liters` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '补水量（升）',
  `processed` tinyint NOT NULL DEFAULT 0 COMMENT '是否已被模型处理：1=已处理，0=未处理',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tank_event_time`(`tank_id` ASC, `event_time` ASC) USING BTREE,
  INDEX `idx_processed`(`processed` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '补水事件表，记录每次向槽液补充去离子水的数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plate_water_events
-- ----------------------------

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_api`;
CREATE TABLE `sys_api`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '接口ID',
  `path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '接口路径，如：/api/system/user',
  `method` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'HTTP方法：GET/POST/PUT/DELETE',
  `group` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接口分组，如：系统管理',
  `description` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接口描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_path_method`(`path` ASC, `method` ASC) USING BTREE,
  INDEX `idx_group`(`group` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统接口权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_api
-- ----------------------------

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `type_id` bigint NOT NULL COMMENT '字典类型ID',
  `label` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典标签（显示名称）',
  `value` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典键值',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序值',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_type_id`(`type_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------
INSERT INTO `sys_dict_data` VALUES (1, 1, '启用', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (2, 1, '禁用', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (3, 2, '目录', '0', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (4, 2, '菜单', '1', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (5, 2, '按钮', '2', 3, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (6, 3, '显示', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (7, 3, '隐藏', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (8, 4, 'A型带出（面积相关）', 'A', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (9, 4, 'B型带出（数量相关）', 'B', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (10, 4, 'C型带出（固定带出）', 'C', 3, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (11, 5, '初始化', 'init', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (12, 5, '模型计算', 'calc', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO `sys_dict_data` VALUES (13, 5, '化验覆盖', 'lab_override', 3, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型编码',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典类型表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dict_type
-- ----------------------------
INSERT INTO `sys_dict_type` VALUES (1, '系统状态', 'sys_status', 1, '通用状态：启用/禁用', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_dict_type` VALUES (2, '菜单类型', 'sys_menu_type', 1, '菜单类型：目录/菜单/按钮', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_dict_type` VALUES (3, '是否显示', 'sys_visible', 1, '菜单是否显示', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_dict_type` VALUES (4, '带出类型', 'plating_carryover_type', 1, '镀铬带出类型：A/B/C', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_dict_type` VALUES (5, '事件来源', 'plating_state_source', 1, '模型状态数据来源类型', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sys_file`;
CREATE TABLE `sys_file`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `filename` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储文件名（UUID生成）',
  `origin_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '原始文件名',
  `file_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件存储路径',
  `file_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件访问URL',
  `file_size` bigint NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
  `file_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件类型/MIME类型',
  `uploader_id` bigint NOT NULL DEFAULT 0 COMMENT '上传人用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_uploader_id`(`uploader_id` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_file
-- ----------------------------

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户ID',
  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录IP',
  `location` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录地点',
  `browser` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '浏览器',
  `os` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '操作系统',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '登录状态：1=成功，0=失败',
  `msg` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '提示消息',
  `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_username`(`username` ASC) USING BTREE,
  INDEX `idx_login_time`(`login_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '登录日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_login_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示顶级',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由路径',
  `component` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '前端组件路径',
  `icon` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单图标',
  `type` tinyint NOT NULL DEFAULT 0 COMMENT '菜单类型：0=目录，1=菜单，2=按钮',
  `permission` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '权限标识符，如：system:user:list',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  `visible` tinyint NOT NULL DEFAULT 1 COMMENT '是否可见：1=可见，0=隐藏',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_type`(`type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, 0, '系统管理', '/system', 'Layout', 'setting', 0, '', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (2, 1, '用户管理', '/system/user', 'system/user/index', 'user', 1, 'system:user:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (3, 1, '角色管理', '/system/role', 'system/role/index', 'peoples', 1, 'system:role:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (4, 1, '菜单管理', '/system/menu', 'system/menu/index', 'tree-table', 1, 'system:menu:list', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (5, 1, '接口管理', '/system/api', 'system/api/index', 'api', 1, 'system:api:list', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (6, 1, '字典管理', '/system/dict', 'system/dict/index', 'dict', 1, 'system:dict:list', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (7, 0, '日志管理', '/log', 'Layout', 'log', 0, '', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (8, 7, '登录日志', '/log/login', 'log/login/index', 'logininfor', 1, 'system:loginlog:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (9, 7, '操作日志', '/log/oper', 'log/oper/index', 'form', 1, 'system:operlog:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (10, 0, '槽液分析', '/plating', 'Layout', 'chart', 0, '', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (11, 10, '槽体配置', '/plating/tank', 'plating/tank/index', 'tank', 1, 'plating:tank:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (12, 10, '生产事件', '/plating/production', 'plating/production/index', 'production', 1, 'plating:production:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (13, 10, '加药管理', '/plating/dosing', 'plating/dosing/index', 'dosing', 1, 'plating:dosing:list', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (14, 10, '补水管理', '/plating/water', 'plating/water/index', 'water', 1, 'plating:water:list', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (15, 10, '模型状态', '/plating/state', 'plating/state/index', 'monitor', 1, 'plating:state:list', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (20, 2, '新增用户', '', '', '', 2, 'system:user:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (21, 2, '编辑用户', '', '', '', 2, 'system:user:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (22, 2, '删除用户', '', '', '', 2, 'system:user:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (23, 3, '新增角色', '', '', '', 2, 'system:role:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (24, 3, '编辑角色', '', '', '', 2, 'system:role:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (25, 3, '删除角色', '', '', '', 2, 'system:role:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (26, 3, '分配菜单', '', '', '', 2, 'system:role:assignmenu', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO `sys_menu` VALUES (27, 3, '分配接口', '', '', '', 2, 'system:role:assignapi', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);

-- ----------------------------
-- Table structure for sys_oper_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_oper_log`;
CREATE TABLE `sys_oper_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '操作日志ID',
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '操作模块标题',
  `business_type` int NOT NULL DEFAULT 0 COMMENT '业务类型：0=其他，1=新增，2=修改，3=删除，4=查询',
  `method` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '方法名称',
  `request_method` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求方式：GET/POST/PUT/DELETE',
  `operator_type` int NOT NULL DEFAULT 0 COMMENT '操作人类型：0=其他，1=后台用户，2=手机端用户',
  `operator_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '操作人员名称',
  `operator_id` bigint NOT NULL DEFAULT 0 COMMENT '操作人员ID',
  `dept_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '部门名称',
  `oper_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求URL',
  `oper_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '操作IP',
  `oper_param` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '请求参数（JSON）',
  `json_result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '返回结果（JSON）',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '操作状态：1=成功，0=失败',
  `error_msg` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '错误消息',
  `oper_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_oper_time`(`oper_time` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_oper_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码（casbin中使用）',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '超级管理员', 'admin', 1, '系统内置超级管理员角色，拥有所有权限', 1, '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);

-- ----------------------------
-- Table structure for sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_api`;
CREATE TABLE `sys_role_api`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `api_id` bigint NOT NULL COMMENT '接口ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_api`(`role_id` ASC, `api_id` ASC) USING BTREE,
  INDEX `idx_api_id`(`api_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色接口关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_api
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_menu`(`role_id` ASC, `menu_id` ASC) USING BTREE,
  INDEX `idx_menu_id`(`menu_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (1, 1, 1);
INSERT INTO `sys_role_menu` VALUES (2, 1, 2);

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录用户名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录密码（bcrypt加密）',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称/显示名',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像URL',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间,空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES (1, 'admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 'admin@example.com', '', '', 1, '系统内置超级管理员账号', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_role`(`user_id` ASC, `role_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO `sys_user_role` VALUES (1, 1, 1);

SET FOREIGN_KEY_CHECKS = 1;
