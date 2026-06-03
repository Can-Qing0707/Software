-- ============================================================
-- 测试数据脚本（仅用于演示和开发测试）
-- ============================================================

USE training_eval_system;

-- ============================================================
-- 1. 用户数据
-- ============================================================
-- 密码均为 123456 (bcrypt hash)
INSERT INTO users (username, password, real_name, role, email, phone, status) VALUES
('teacher01', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '张老师', 'teacher', 'zhang@edu.cn', '13800001001', 1),
('teacher02', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '李老师', 'teacher', 'li@edu.cn', '13800001002', 1),
('student01', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '王小明', 'student', 'wangxm@stu.edu.cn', '13800002001', 1),
('student02', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '赵小红', 'student', 'zhaoxh@stu.edu.cn', '13800002002', 1),
('student03', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '刘小刚', 'student', 'liuxg@stu.edu.cn', '13800002003', 1),
('student04', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '陈小美', 'student', 'chenxm@stu.edu.cn', '13800002004', 1),
('student05', '$2a$10$4YWKSJuxNdW19SpJyQQhfOuHLIxsIWlrndGPh7y0lTXt1/fvhYZC6', '孙小亮', 'student', 'sunxl@edu.cn', '13800002005', 1);

-- ============================================================
-- 2. 课程数据
-- ============================================================
INSERT INTO courses (name, description, teacher_id, semester, code, status) VALUES
('Web前端开发技术',        'Vue3 + Element Plus 企业级前端开发实训',             2, '2025-2026-2', 'WEB101', 1),
('Go语言后端开发',         'Gin框架 + GORM + MySQL 后端开发实训',               2, '2025-2026-2', 'GO201',  1),
('软件工程综合实践',       '全栈项目开发全流程实训',                             3, '2025-2026-2', 'SW301',  1);

-- ============================================================
-- 2.1 选课数据（学生加入课程）
-- ============================================================
INSERT INTO course_enrollments (course_id, student_id) VALUES
-- 学生01-05 加入 Web前端开发技术
(1, 4), (1, 5), (1, 6), (1, 7), (1, 8),
-- 学生01-04 加入 Go语言后端开发
(2, 4), (2, 5), (2, 6), (2, 7),
-- 学生03-05 加入 软件工程综合实践
(3, 6), (3, 7), (3, 8);

-- ============================================================
-- 3. 实训任务数据
-- ============================================================
INSERT INTO tasks (course_id, title, description, deadline, status) VALUES
(1, 'Vue3组件化开发实战', '使用Vue3 Composition API开发一个用户管理系统，包含用户列表、新增、编辑、删除功能。\n要求：\n1. 使用Element Plus组件库\n2. 拆分至少3个可复用组件\n3. 使用Pinia进行状态管理\n4. 添加路由守卫\n5. 代码注释规范完整', '2026-05-30 23:59:59', 1),
(1, '前端性能优化实践', '对已有项目进行性能分析和优化，提交优化前后对比报告。\n要求：\n1. 使用Lighthouse进行性能评测\n2. 至少优化3个性能瓶颈\n3. 提交优化前后数据对比\n4. 撰写优化技术文档', '2026-06-15 23:59:59', 1),
(2, 'RESTful API设计与实现', '设计并实现一个课程管理系统的RESTful API。\n要求：\n1. 遵循RESTful设计规范\n2. 使用Gin框架\n3. 实现JWT认证中间件\n4. 完整的CRUD操作\n5. 统一的错误处理与响应格式\n6. 编写API文档', '2026-06-10 23:59:59', 1),
(2, '数据库设计与GORM实践', '设计一个在线考试系统的数据库模型并使用GORM实现。\n要求：\n1. 至少5张关联表\n2. 包含多对多关系\n3. 实现事务操作\n4. 编写不少于10个常用查询\n5. 添加数据库索引优化', '2026-06-25 23:59:59', 1),
(3, '全栈项目需求分析与设计', '选择一个实际场景完成需求分析与系统设计。\n要求：\n1. 编写需求规格说明书\n2. 绘制用例图、ER图、系统架构图\n3. 设计数据库ER模型\n4. 制定项目开发计划\n5. 完成原型设计', '2026-05-20 23:59:59', 1);

-- ============================================================
-- 4. 任务-指标关联（自定义各任务权重）
-- ============================================================
-- 课程1-任务1: Vue3组件化开发实战
INSERT INTO task_indicators (task_id, indicator_id, weight) VALUES
(1, 1, 30.00),   -- 代码质量 30%
(1, 2, 15.00),   -- 文档规范性 15%
(1, 3, 35.00),   -- 功能实现度 35%
(1, 4, 10.00),   -- 创新性 10%
(1, 5, 10.00);   -- 界面与交互 10%

-- 课程2-任务1: RESTful API设计与实现
INSERT INTO task_indicators (task_id, indicator_id, weight) VALUES
(3, 1, 40.00),   -- 代码质量 40%
(3, 2, 25.00),   -- 文档规范性 25%
(3, 3, 25.00),   -- 功能实现度 25%
(3, 4, 10.00),   -- 创新性 10%
(3, 5, 0.00);    -- 界面与交互 0% (纯后端任务)

-- 课程3-任务1: 全栈项目需求分析与设计
INSERT INTO task_indicators (task_id, indicator_id, weight) VALUES
(5, 1, 10.00),   -- 代码质量 10%
(5, 2, 40.00),   -- 文档规范性 40%
(5, 3, 20.00),   -- 功能实现度 20%
(5, 4, 20.00),   -- 创新性 20%
(5, 5, 10.00);   -- 界面与交互 10%

-- ============================================================
-- 5. 提交数据（部分学生已提交）
-- ============================================================
INSERT INTO submissions (task_id, student_id, files_json, content_text, status, submit_time) VALUES
(1, 4,
 '[{"name":"user-crud.zip","url":"/uploads/task1/student00/user-crud.zip","type":"zip","size":28400}]',
 '实现了用户管理系统，包含登录注册、用户列表分页、新增编辑删除功能、角色权限控制。前端使用Vue3 Composition API + Element Plus，使用Pinia管理状态，Vue Router配置了路由守卫。',
 'uploaded', '2026-05-09 08:00:00'),
(1, 5,
 '[{"name":"UserManage.vue","url":"/uploads/task1/student01/UserManage.vue","type":"vue","size":12560},{"name":"UserForm.vue","url":"/uploads/task1/student01/UserForm.vue","type":"vue","size":8340}]',
 '使用Vue3 Composition API开发了用户管理系统，实现了用户列表展示、新增用户表单、编辑用户、删除用户等功能。项目使用Element Plus组件库，拆分出UserTable、UserForm、UserSearch三个可复用组件，使用Pinia管理用户状态，配置了路由守卫进行权限控制。',
 'uploaded', '2026-05-10 14:30:00'),
(1, 6,
 '[{"name":"user-system.zip","url":"/uploads/task1/student02/user-system.zip","type":"zip","size":45820}]',
 '完成了用户管理系统的开发，基于Vue3和Element Plus，实现了完整的用户CRUD功能，包含搜索、分页、批量删除等进阶功能。',
 'uploaded', '2026-05-11 09:15:00'),
(1, 7,
 '[{"name":"code.zip","url":"/uploads/task1/student03/code.zip","type":"zip","size":32100}]',
 '完成了Vue3组件化开发实战，包含用户管理、角色管理两个模块，使用了Element Plus和Pinia。',
 'uploaded', '2026-05-12 16:45:00'),
(3, 5,
 '[{"name":"api_doc.md","url":"/uploads/task3/student01/api_doc.md","type":"md","size":5600}]',
 '设计了课程管理系统的RESTful API，包含课程、教师、学生三个资源，使用Gin框架实现，集成了JWT认证，统一错误处理。',
 'uploaded', '2026-05-13 11:20:00'),
(5, 6,
 '[{"name":"需求分析文档.pdf","url":"/uploads/task5/student02/需求分析文档.pdf","type":"pdf","size":245800}]',
 '选择在线教育平台作为实际场景，完成了需求规格说明书，包含用例图、ER图、系统架构图、数据库设计和开发计划。',
 'uploaded', '2026-05-08 10:00:00');

-- ============================================================
-- 6. 核查结果数据（部分已核查）
-- ============================================================
INSERT INTO verification_results (submission_id, completeness, logic_issues, requirement_match, overall_pass, verified_at) VALUES
(1,
 '{"steps_total":5,"steps_completed":4,"missing_steps":["路由守卫未配置"],"completeness_ratio":0.8}',
 '{"has_logic_issues":false,"issues":[],"summary":"无明显逻辑漏洞"}',
 '{"matched_requirements":["使用ElementPlus","拆分可复用组件","使用Pinia","代码注释"],"unmatched_requirements":["路由守卫"],"match_ratio":0.8}',
 1, '2026-05-11 10:00:00'),
(3,
 '{"steps_total":6,"steps_completed":5,"missing_steps":["API文档不完整"],"completeness_ratio":0.83}',
 '{"has_logic_issues":false,"issues":[],"summary":"逻辑基本正确"}',
 '{"matched_requirements":["RESTful规范","Gin框架","JWT认证","CRUD操作","统一错误处理"],"unmatched_requirements":["API文档不够详细"],"match_ratio":0.83}',
 1, '2026-05-14 15:30:00');

-- ============================================================
-- 7. 评分数据（部分已评分）
-- ============================================================
INSERT INTO eval_scores (submission_id, indicator_id, llm_score, llm_comment, final_score) VALUES
-- 学生01 任务1 评分
(1, 1, 82.00, '代码结构清晰，命名规范，组件拆分合理，但部分逻辑可进一步抽象复用', 82.00),
(1, 2, 75.00, '文档基本完整，部分注释可更详细', 75.00),
(1, 3, 85.00, '核心功能全部实现，交互流畅', 85.00),
(1, 4, 70.00, '技术方案常规，缺乏创新亮点', 70.00),
(1, 5, 80.00, '界面美观，交互体验良好', 80.00),

-- 学生02 任务1 评分
(2, 1, 78.00, '代码可读性良好，但部分函数过长需要拆分', 78.00),
(2, 2, 70.00, '文档编写较为简单，缺少详细说明', 70.00),
(2, 3, 88.00, '功能实现完整，包含进阶功能', 88.00),
(2, 4, 75.00, '批量删除功能有一定创新', 75.00),
(2, 5, 82.00, '界面设计美观，交互流畅', 82.00);

-- ============================================================
-- 8. 报告记录数据
-- ============================================================
INSERT INTO reports (submission_id, type, format, file_url, title, generated_by) VALUES
(1, 'individual', 'pdf', '/reports/individual_001.pdf', '王小明 - Vue3组件化开发实战 评价报告', 1);
