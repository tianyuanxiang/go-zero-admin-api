-- 镀铬工艺平台初始化SQL
-- 数据库：plating
-- 创建时间：2024-01-01
-- 说明：包含系统管理表和业务表的完整DDL及初始数据

CREATE DATABASE IF NOT EXISTS plating DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE plating;

-- ----------------------------
-- 系统用户表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_user` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username`   VARCHAR(64)  NOT NULL                COMMENT '登录用户名',
  `password`   VARCHAR(128) NOT NULL                COMMENT '登录密码（bcrypt加密）',
  `nickname`   VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '昵称/显示名',
  `email`      VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '邮箱',
  `phone`      VARCHAR(20)  NOT NULL DEFAULT ''     COMMENT '手机号',
  `avatar`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '头像URL',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态：1=启用，0=禁用',
  `remark`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间,空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- ----------------------------
-- 系统角色表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_role` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name`       VARCHAR(64)  NOT NULL                COMMENT '角色名称',
  `code`       VARCHAR(64)  NOT NULL                COMMENT '角色编码（casbin中使用）',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态：1=启用，0=禁用',
  `remark`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `sort`       INT          NOT NULL DEFAULT 0      COMMENT '排序值，越小越靠前',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统角色表';

-- ----------------------------
-- 用户角色关联表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_user_role` (
  `id`      BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL                COMMENT '用户ID',
  `role_id` BIGINT NOT NULL                COMMENT '角色ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- ----------------------------
-- 系统菜单表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_menu` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id`  BIGINT       NOT NULL DEFAULT 0      COMMENT '父菜单ID，0表示顶级',
  `name`       VARCHAR(64)  NOT NULL                COMMENT '菜单名称',
  `path`       VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '路由路径',
  `component`  VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '前端组件路径',
  `icon`       VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '菜单图标',
  `type`       TINYINT      NOT NULL DEFAULT 0      COMMENT '菜单类型：0=目录，1=菜单，2=按钮',
  `permission` VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '权限标识符，如：system:user:list',
  `sort`       INT          NOT NULL DEFAULT 0      COMMENT '排序值，越小越靠前',
  `visible`    TINYINT      NOT NULL DEFAULT 1      COMMENT '是否可见：1=可见，0=隐藏',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态：1=启用，0=禁用',
  `remark`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表';

-- ----------------------------
-- 角色菜单关联表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_role_menu` (
  `id`      BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` BIGINT NOT NULL                COMMENT '角色ID',
  `menu_id` BIGINT NOT NULL                COMMENT '菜单ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';

-- ----------------------------
-- 系统接口权限表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_api` (
  `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '接口ID',
  `path`        VARCHAR(256) NOT NULL                COMMENT '接口路径，如：/api/system/user',
  `method`      VARCHAR(16)  NOT NULL                COMMENT 'HTTP方法：GET/POST/PUT/DELETE',
  `group`       VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '接口分组，如：系统管理',
  `description` VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '接口描述',
  `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_path_method` (`path`, `method`),
  KEY `idx_group` (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统接口权限表';

-- ----------------------------
-- 角色接口关联表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_role_api` (
  `id`     BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` BIGINT NOT NULL               COMMENT '角色ID',
  `api_id`  BIGINT NOT NULL               COMMENT '接口ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_api` (`role_id`, `api_id`),
  KEY `idx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色接口关联表';

-- ----------------------------
-- 字典类型表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_dict_type` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name`       VARCHAR(64)  NOT NULL                COMMENT '字典类型名称',
  `code`       VARCHAR(64)  NOT NULL                COMMENT '字典类型编码',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态：1=启用，0=禁用',
  `remark`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典类型表';

-- ----------------------------
-- 字典数据表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_dict_data` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `type_id`    BIGINT       NOT NULL                COMMENT '字典类型ID',
  `label`      VARCHAR(128) NOT NULL                COMMENT '字典标签（显示名称）',
  `value`      VARCHAR(128) NOT NULL                COMMENT '字典键值',
  `sort`       INT          NOT NULL DEFAULT 0      COMMENT '排序值',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态：1=启用，0=禁用',
  `remark`     VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_id` (`type_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典数据表';

-- ----------------------------
-- 登录日志表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_login_log` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id`    BIGINT       NOT NULL DEFAULT 0      COMMENT '用户ID',
  `username`   VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '用户名',
  `ip`         VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '登录IP',
  `location`   VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '登录地点',
  `browser`    VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '浏览器',
  `os`         VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '操作系统',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '登录状态：1=成功，0=失败',
  `msg`        VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '提示消息',
  `login_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- ----------------------------
-- 操作日志表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_oper_log` (
  `id`             BIGINT       NOT NULL AUTO_INCREMENT COMMENT '操作日志ID',
  `title`          VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '操作模块标题',
  `business_type`  INT          NOT NULL DEFAULT 0      COMMENT '业务类型：0=其他，1=新增，2=修改，3=删除，4=查询',
  `method`         VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '方法名称',
  `request_method` VARCHAR(16)  NOT NULL DEFAULT ''     COMMENT '请求方式：GET/POST/PUT/DELETE',
  `operator_type`  INT          NOT NULL DEFAULT 0      COMMENT '操作人类型：0=其他，1=后台用户，2=手机端用户',
  `operator_name`  VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '操作人员名称',
  `operator_id`    BIGINT       NOT NULL DEFAULT 0      COMMENT '操作人员ID',
  `dept_name`      VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '部门名称',
  `oper_url`       VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '请求URL',
  `oper_ip`        VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '操作IP',
  `oper_param`     TEXT                                 COMMENT '请求参数（JSON）',
  `json_result`    TEXT                                 COMMENT '返回结果（JSON）',
  `status`         TINYINT      NOT NULL DEFAULT 1      COMMENT '操作状态：1=成功，0=失败',
  `error_msg`      VARCHAR(2000) NOT NULL DEFAULT ''    COMMENT '错误消息',
  `oper_time`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_oper_time` (`oper_time`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- ----------------------------
-- 文件记录表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_file` (
  `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `filename`    VARCHAR(256) NOT NULL                COMMENT '存储文件名（UUID生成）',
  `origin_name` VARCHAR(256) NOT NULL DEFAULT ''     COMMENT '原始文件名',
  `file_path`   VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '文件存储路径',
  `file_url`    VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '文件访问URL',
  `file_size`   BIGINT       NOT NULL DEFAULT 0      COMMENT '文件大小（字节）',
  `file_type`   VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '文件类型/MIME类型',
  `uploader_id` BIGINT       NOT NULL DEFAULT 0      COMMENT '上传人用户ID',
  `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_uploader_id` (`uploader_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件记录表';

-- ----------------------------
-- Casbin规则表（由gorm-adapter管理，此处也提供DDL方便查看）
-- ----------------------------
CREATE TABLE IF NOT EXISTS `casbin_rule` (
  `id`    BIGINT      NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ptype` VARCHAR(100) DEFAULT ''             COMMENT '策略类型：p或g',
  `v0`    VARCHAR(100) DEFAULT ''             COMMENT '第0个字段（通常为sub/role）',
  `v1`    VARCHAR(100) DEFAULT ''             COMMENT '第1个字段（通常为obj/role）',
  `v2`    VARCHAR(100) DEFAULT ''             COMMENT '第2个字段（通常为act）',
  `v3`    VARCHAR(100) DEFAULT ''             COMMENT '第3个字段（扩展）',
  `v4`    VARCHAR(100) DEFAULT ''             COMMENT '第4个字段（扩展）',
  `v5`    VARCHAR(100) DEFAULT ''             COMMENT '第5个字段（扩展）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Casbin权限规则表';

-- ----------------------------
-- 槽体配置表
-- ----------------------------
CREATE TABLE IF NOT EXISTS `tank_config` (
  `tank_id`         VARCHAR(32)  NOT NULL                COMMENT '槽体ID（主键，如：TANK-001）',
  `tank_name`       VARCHAR(128) NOT NULL                COMMENT '槽体名称',
  `volume_liters`   DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '槽体容积（升）',
  `density_g_per_l` DECIMAL(10,4) NOT NULL DEFAULT 0.00 COMMENT '镀液密度（g/L）',
  `d_cro3_g_per_ah` DECIMAL(10,6) NOT NULL DEFAULT 0.00 COMMENT 'CrO3消耗系数（g/Ah），法拉第系数',
  `d_cr3_g_per_ah`  DECIMAL(10,6) NOT NULL DEFAULT 0.00 COMMENT 'Cr3+生成系数（g/Ah）',
  `carryover_a`     DECIMAL(10,6) NOT NULL DEFAULT 0.00 COMMENT '带出系数A（随工件面积）',
  `carryover_b`     DECIMAL(10,6) NOT NULL DEFAULT 0.00 COMMENT '带出系数B（随工件数量）',
  `carryover_c`     DECIMAL(10,6) NOT NULL DEFAULT 0.00 COMMENT '带出系数C（固定带出）',
  `is_active`       TINYINT      NOT NULL DEFAULT 1      COMMENT '是否激活：1=激活，0=停用',
  `remark`          VARCHAR(512) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`tank_id`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='槽体配置表';

-- ----------------------------
-- 模型状态表（记录每次计算后的浓度状态）
-- ----------------------------
CREATE TABLE IF NOT EXISTS `model_state` (
  `id`          BIGINT        NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `tank_id`     VARCHAR(32)   NOT NULL                COMMENT '关联槽体ID',
  `calc_time`   DATETIME      NOT NULL                COMMENT '计算/记录时间',
  `cro3_g_per_l` DECIMAL(10,4) NOT NULL DEFAULT 0.00 COMMENT '铬酸浓度（g/L，CrO3）',
  `cr3_g_per_l`  DECIMAL(10,4) NOT NULL DEFAULT 0.00 COMMENT '三价铬浓度（g/L，Cr3+）',
  `tank_volume`  DECIMAL(10,2) NOT NULL DEFAULT 0.00  COMMENT '当前槽液体积（升）',
  `source`      VARCHAR(32)   NOT NULL DEFAULT 'calc' COMMENT '数据来源：init=初始化，calc=模型计算，lab_override=化验覆盖',
  `remark`      VARCHAR(512)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tank_calc_time` (`tank_id`, `calc_time`),
  KEY `idx_calc_time` (`calc_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模型状态表，记录槽液浓度的历史计算状态';

-- ----------------------------
-- 生产事件表（记录每批次生产）
-- ----------------------------
CREATE TABLE IF NOT EXISTS `production_events` (
  `id`               BIGINT        NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id`          VARCHAR(32)   NOT NULL                COMMENT '关联槽体ID',
  `event_time`       DATETIME      NOT NULL                COMMENT '事件发生时间',
  `batch_start_time` DATETIME      NOT NULL                COMMENT '批次开始时间',
  `part_count`       INT           NOT NULL DEFAULT 0      COMMENT '本批次工件数量（件）',
  `part_area_dm2`    DECIMAL(10,4) NOT NULL DEFAULT 0.00   COMMENT '工件总面积（dm2）',
  `carryover_type`   CHAR(1)       NOT NULL DEFAULT 'A'    COMMENT '带出类型：A/B/C（对应carryover_a/b/c）',
  `ah_this_batch`    DECIMAL(12,4) NOT NULL DEFAULT 0.00   COMMENT '本批次实际Ah数（安培时）',
  `ah_calculated`    DECIMAL(12,4) NOT NULL DEFAULT 0.00   COMMENT '模型计算使用的Ah（可能来自TDEngine或手工录入）',
  `processed`        TINYINT       NOT NULL DEFAULT 0      COMMENT '是否已被模型计算处理：1=已处理，0=未处理',
  `remark`           VARCHAR(512)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tank_event_time` (`tank_id`, `event_time`),
  KEY `idx_processed` (`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='生产事件表，记录每批次镀铬生产数据';

-- ----------------------------
-- 加药事件表（记录每次加入化学品）
-- ----------------------------
CREATE TABLE IF NOT EXISTS `dosing_events` (
  `id`            BIGINT        NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id`       VARCHAR(32)   NOT NULL                COMMENT '关联槽体ID',
  `event_time`    DATETIME      NOT NULL                COMMENT '加药时间',
  `chemical`      VARCHAR(64)   NOT NULL                COMMENT '化学品名称，如：CrO3、H2SO4',
  `amount_grams`  DECIMAL(12,4) NOT NULL DEFAULT 0.00   COMMENT '加入量（克）',
  `processed`     TINYINT       NOT NULL DEFAULT 0      COMMENT '是否已被模型处理：1=已处理，0=未处理',
  `remark`        VARCHAR(512)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tank_event_time` (`tank_id`, `event_time`),
  KEY `idx_chemical` (`chemical`),
  KEY `idx_processed` (`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='加药事件表，记录每次向槽液中添加化学品的数据';

-- ----------------------------
-- 补水事件表（记录每次补充去离子水）
-- ----------------------------
CREATE TABLE IF NOT EXISTS `water_events` (
  `id`          BIGINT        NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `tank_id`     VARCHAR(32)   NOT NULL                COMMENT '关联槽体ID',
  `event_time`  DATETIME      NOT NULL                COMMENT '补水时间',
  `water_liters` DECIMAL(10,4) NOT NULL DEFAULT 0.00  COMMENT '补水量（升）',
  `processed`   TINYINT       NOT NULL DEFAULT 0      COMMENT '是否已被模型处理：1=已处理，0=未处理',
  `remark`      VARCHAR(512)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` DATETIME              DEFAULT NULL   COMMENT '软删除时间, 空-未删除，非空为已删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tank_event_time` (`tank_id`, `event_time`),
  KEY `idx_processed` (`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='补水事件表，记录每次向槽液补充去离子水的数据';

-- ============================
-- 初始化数据
-- ============================

-- 插入admin用户（密码：Admin@123，bcrypt hash）
INSERT INTO `sys_user` (`id`, `username`, `password`, `nickname`, `email`, `status`, `remark`)
VALUES (1, 'admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 'admin@example.com', 1, '系统内置超级管理员账号')
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 插入超级管理员角色
INSERT INTO `sys_role` (`id`, `name`, `code`, `status`, `remark`, `sort`)
VALUES (1, '超级管理员', 'admin', 1, '系统内置超级管理员角色，拥有所有权限', 1)
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 关联admin用户与超级管理员角色
INSERT INTO `sys_user_role` (`user_id`, `role_id`)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 初始化菜单数据
-- 系统管理目录
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `type`, `permission`, `sort`, `visible`, `status`)
VALUES
(1, 0, '系统管理', '/system', 'Layout', 'setting', 0, '', 1, 1, 1),
(2, 1, '用户管理', '/system/user', 'system/user/index', 'user', 1, 'system:user:list', 1, 1, 1),
(3, 1, '角色管理', '/system/role', 'system/role/index', 'peoples', 1, 'system:role:list', 2, 1, 1),
(4, 1, '菜单管理', '/system/menu', 'system/menu/index', 'tree-table', 1, 'system:menu:list', 3, 1, 1),
(5, 1, '接口管理', '/system/api', 'system/api/index', 'api', 1, 'system:api:list', 4, 1, 1),
(6, 1, '字典管理', '/system/dict', 'system/dict/index', 'dict', 1, 'system:dict:list', 5, 1, 1),
(7, 0, '日志管理', '/log', 'Layout', 'log', 0, '', 2, 1, 1),
(8, 7, '登录日志', '/log/login', 'log/login/index', 'logininfor', 1, 'system:loginlog:list', 1, 1, 1),
(9, 7, '操作日志', '/log/oper', 'log/oper/index', 'form', 1, 'system:operlog:list', 2, 1, 1),
(10, 0, '槽液分析', '/plating', 'Layout', 'chart', 0, '', 3, 1, 1),
(11, 10, '槽体配置', '/plating/tank', 'plating/tank/index', 'tank', 1, 'plating:tank:list', 1, 1, 1),
(12, 10, '生产事件', '/plating/production', 'plating/production/index', 'production', 1, 'plating:production:list', 2, 1, 1),
(13, 10, '加药管理', '/plating/dosing', 'plating/dosing/index', 'dosing', 1, 'plating:dosing:list', 3, 1, 1),
(14, 10, '补水管理', '/plating/water', 'plating/water/index', 'water', 1, 'plating:water:list', 4, 1, 1),
(15, 10, '模型状态', '/plating/state', 'plating/state/index', 'monitor', 1, 'plating:state:list', 5, 1, 1)
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 按钮权限菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `type`, `permission`, `sort`, `visible`, `status`)
VALUES
(20, 2, '新增用户', '', '', '', 2, 'system:user:create', 1, 1, 1),
(21, 2, '编辑用户', '', '', '', 2, 'system:user:update', 2, 1, 1),
(22, 2, '删除用户', '', '', '', 2, 'system:user:delete', 3, 1, 1),
(23, 3, '新增角色', '', '', '', 2, 'system:role:create', 1, 1, 1),
(24, 3, '编辑角色', '', '', '', 2, 'system:role:update', 2, 1, 1),
(25, 3, '删除角色', '', '', '', 2, 'system:role:delete', 3, 1, 1),
(26, 3, '分配菜单', '', '', '', 2, 'system:role:assignmenu', 4, 1, 1),
(27, 3, '分配接口', '', '', '', 2, 'system:role:assignapi', 5, 1, 1)
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 初始化字典类型
INSERT INTO `sys_dict_type` (`id`, `name`, `code`, `status`, `remark`)
VALUES
(1, '系统状态', 'sys_status', 1, '通用状态：启用/禁用'),
(2, '菜单类型', 'sys_menu_type', 1, '菜单类型：目录/菜单/按钮'),
(3, '是否显示', 'sys_visible', 1, '菜单是否显示'),
(4, '带出类型', 'plating_carryover_type', 1, '镀铬带出类型：A/B/C'),
(5, '事件来源', 'plating_state_source', 1, '模型状态数据来源类型')
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 初始化字典数据
INSERT INTO `sys_dict_data` (`type_id`, `label`, `value`, `sort`, `status`)
VALUES
(1, '启用', '1', 1, 1),
(1, '禁用', '0', 2, 1),
(2, '目录', '0', 1, 1),
(2, '菜单', '1', 2, 1),
(2, '按钮', '2', 3, 1),
(3, '显示', '1', 1, 1),
(3, '隐藏', '0', 2, 1),
(4, 'A型带出（面积相关）', 'A', 1, 1),
(4, 'B型带出（数量相关）', 'B', 2, 1),
(4, 'C型带出（固定带出）', 'C', 3, 1),
(5, '初始化', 'init', 1, 1),
(5, '模型计算', 'calc', 2, 1),
(5, '化验覆盖', 'lab_override', 3, 1)
ON DUPLICATE KEY UPDATE `id` = `id`;

-- 初始化casbin超级管理员角色拥有所有权限（通配符）
-- 实际项目中应通过API动态分配，此处仅做示例
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`)
VALUES ('p', 'admin', '/api/*', '*')
ON DUPLICATE KEY UPDATE `id` = `id`;
