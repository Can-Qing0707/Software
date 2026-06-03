


SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

DROP DATABASE IF EXISTS training_eval_system;
CREATE DATABASE IF NOT EXISTS training_eval_system
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;


USE training_eval_system;

-- ============================================================
-- 1. 用户表
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username    VARCHAR(64)  NOT NULL UNIQUE COMMENT '登录用户名',
  password    VARCHAR(256) NOT NULL        COMMENT '密码',
  real_name   VARCHAR(64)  NOT NULL        COMMENT '真实姓名',
  role        ENUM('admin','teacher','student') NOT NULL DEFAULT 'student' COMMENT '角色',
  email       VARCHAR(128) DEFAULT NULL    COMMENT '邮箱',
  phone       VARCHAR(20)  DEFAULT NULL    COMMENT '手机号',
  avatar      VARCHAR(256) DEFAULT NULL    COMMENT '头像URL',
  status      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态 1-正常 0-禁用',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_role (role),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================================
-- 2. 课程表
-- ============================================================
CREATE TABLE IF NOT EXISTS courses (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name          VARCHAR(128) NOT NULL        COMMENT '课程名称',
  description   TEXT                         COMMENT '课程描述',
  teacher_id    BIGINT UNSIGNED NOT NULL     COMMENT '授课教师ID',
  semester      VARCHAR(32)  DEFAULT NULL    COMMENT '学期',
  code          VARCHAR(32)  NOT NULL        COMMENT '课程代码',
  status        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态 1-进行中 0-已结束',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_code (code),
  INDEX idx_teacher (teacher_id),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- ============================================================
-- 2.1 课程选课关系表
-- ============================================================
CREATE TABLE IF NOT EXISTS course_enrollments (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  course_id   BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
  student_id  BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_course_student (course_id, student_id),
  INDEX idx_course (course_id),
  INDEX idx_student (student_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程选课关系表';

-- ============================================================
-- 3. 实训任务表
-- ============================================================
CREATE TABLE IF NOT EXISTS tasks (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  course_id       BIGINT UNSIGNED NOT NULL     COMMENT '所属课程ID',
  title           VARCHAR(256) NOT NULL        COMMENT '任务标题',
  description     TEXT                         COMMENT '任务描述/要求',
  attachment_url  VARCHAR(512) DEFAULT NULL    COMMENT '任务附件URL',
  deadline        DATETIME     DEFAULT NULL     COMMENT '截止时间',
  status          TINYINT      NOT NULL DEFAULT 1 COMMENT '状态 1-发布 0-草稿',
  created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_course (course_id),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='实训任务表';

-- ============================================================
-- 4. 实训成果提交表
-- ============================================================
CREATE TABLE IF NOT EXISTS submissions (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  task_id       BIGINT UNSIGNED NOT NULL     COMMENT '关联任务ID',
  student_id    BIGINT UNSIGNED NOT NULL     COMMENT '提交学生ID',
  files_json    JSON             NOT NULL     COMMENT '文件列表JSON [{name,url,type,size}]',
  content_text  LONGTEXT         DEFAULT NULL COMMENT 'LLM提取的文本内容',
  status        ENUM('uploaded','parsing','parsed','verified','evaluated')
                NOT NULL DEFAULT 'uploaded'   COMMENT '提交状态',
  submit_time   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_task_student (task_id, student_id),
  INDEX idx_task (task_id),
  INDEX idx_student (student_id),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='实训成果提交表';

-- ============================================================
-- 5. 核查结果表
-- ============================================================
CREATE TABLE IF NOT EXISTS verification_results (
  id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  submission_id     BIGINT UNSIGNED NOT NULL     COMMENT '关联提交ID',
  completeness      TEXT             COMMENT '步骤完整性核查结果(JSON)',
  logic_issues      TEXT             COMMENT '逻辑漏洞核查结果(JSON)',
  requirement_match TEXT             COMMENT '要求匹配度核查结果(JSON)',
  overall_pass      TINYINT(1)      DEFAULT NULL COMMENT '是否整体通过 1-通过 0-不通过',
  raw_llm_response  LONGTEXT         DEFAULT NULL COMMENT 'LLM原始返回',
  verified_at       DATETIME         DEFAULT NULL COMMENT '核查时间',
  created_at        DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_submission (submission_id),
  INDEX idx_pass (overall_pass)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='核查结果表';

-- ============================================================
-- 6. 评价指标表（全局指标库）
-- ============================================================
CREATE TABLE IF NOT EXISTS eval_indicators (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name          VARCHAR(128) NOT NULL        COMMENT '指标名称',
  description   VARCHAR(512) DEFAULT NULL    COMMENT '指标说明/评价标准',
  default_weight DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '默认权重(0-100)',
  sort_order    INT          NOT NULL DEFAULT 0 COMMENT '排序号',
  status        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评价指标表(全局指标库)';

-- ============================================================
-- 7. 课程指标关联表（每个课程可自定义默认权重）
-- ============================================================
CREATE TABLE IF NOT EXISTS course_indicators (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  course_id       BIGINT UNSIGNED NOT NULL     COMMENT '关联课程ID',
  indicator_id    BIGINT UNSIGNED NOT NULL     COMMENT '关联指标ID',
  weight          DECIMAL(5,2)   NOT NULL DEFAULT 0 COMMENT '本课程中该指标权重(0-100)',
  created_at      DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_course_indicator (course_id, indicator_id),
  INDEX idx_course (course_id),
  INDEX idx_indicator (indicator_id),
  FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  FOREIGN KEY (indicator_id) REFERENCES eval_indicators(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程指标关联表';

-- ============================================================
-- 8. 任务指标关联表（每个任务可自定义权重，覆盖课程默认）
-- ============================================================
-- ============================================================
CREATE TABLE IF NOT EXISTS task_indicators (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  task_id         BIGINT UNSIGNED NOT NULL     COMMENT '关联任务ID',
  indicator_id    BIGINT UNSIGNED NOT NULL     COMMENT '关联指标ID',
  weight          DECIMAL(5,2)   NOT NULL DEFAULT 0 COMMENT '本任务中该指标权重(0-100)',
  created_at      DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_task_indicator (task_id, indicator_id),
  INDEX idx_task (task_id),
  INDEX idx_indicator (indicator_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务指标关联表';

-- ============================================================
-- 9. 评分结果表
-- ============================================================
CREATE TABLE IF NOT EXISTS eval_scores (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  submission_id   BIGINT UNSIGNED NOT NULL     COMMENT '关联提交ID',
  indicator_id    BIGINT UNSIGNED NOT NULL     COMMENT '关联指标ID',
  llm_score       DECIMAL(5,2)    DEFAULT NULL COMMENT 'LLM评分(0-100)',
  llm_comment     TEXT            DEFAULT NULL COMMENT 'LLM评语',
  teacher_score   DECIMAL(5,2)    DEFAULT NULL COMMENT '教师评分(0-100)',
  teacher_comment TEXT            DEFAULT NULL COMMENT '教师评语',
  final_score     DECIMAL(5,2)    DEFAULT NULL COMMENT '最终得分(加权计算)',
  created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_submission_indicator (submission_id, indicator_id),
  INDEX idx_submission (submission_id),
  INDEX idx_indicator (indicator_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评分结果表';

-- ============================================================
-- 10. 报告记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS reports (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  submission_id   BIGINT UNSIGNED DEFAULT NULL COMMENT '关联提交ID(NULL=统计报表)',
  course_id       BIGINT UNSIGNED DEFAULT NULL COMMENT '关联课程ID(统计报表)',
  task_id         BIGINT UNSIGNED DEFAULT NULL COMMENT '关联任务ID(统计报表)',
  type            ENUM('individual','class','course') NOT NULL COMMENT '报告类型',
  format          ENUM('pdf','excel') NOT NULL         COMMENT '导出格式',
  file_url        VARCHAR(512)     NOT NULL             COMMENT '文件存储路径',
  title           VARCHAR(256)     DEFAULT NULL         COMMENT '报告标题',
  generated_by    BIGINT UNSIGNED NOT NULL              COMMENT '生成人ID',
  generated_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_submission (submission_id),
  INDEX idx_course (course_id),
  INDEX idx_task (task_id),
  INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='报告记录表';

-- ============================================================
-- 11. 系统配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS system_config (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  config_key    VARCHAR(128) NOT NULL UNIQUE COMMENT '配置键',
  config_value  TEXT         NOT NULL        COMMENT '配置值(JSON)',
  description   VARCHAR(256) DEFAULT NULL    COMMENT '配置说明',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ============================================================
-- 初始数据
-- ============================================================

-- 默认管理员 (密码: admin123)
INSERT INTO users (username, password, real_name, role, status) VALUES
('admin', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '管理员', 'admin', 1);

-- 默认评价指标
INSERT INTO eval_indicators (name, description, default_weight, sort_order) VALUES
('代码质量',      '代码结构、可读性、注释规范、命名规范',         25.00, 1),
('文档规范性',    '文档结构完整、格式规范、表达清晰',             20.00, 2),
('功能实现度',    '功能完成情况、核心需求实现完整性',             30.00, 3),
('创新性',        '技术选型、设计思路、解决方案的创新程度',        15.00, 4),
('界面与交互',    '界面美观度、交互流畅度、用户体验',             10.00, 5);

-- 默认LLM配置
INSERT INTO system_config (config_key, config_value, description) VALUES
('llm_provider',  '{"provider":"openai","api_url":"https://api.openai.com/v1","api_key":"","model":"gpt-4o","max_tokens":4096,"temperature":0.3}', 'LLM服务配置');
