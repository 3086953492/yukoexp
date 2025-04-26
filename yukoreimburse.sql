/*
 Navicat Premium Dump SQL

 Source Server         : 本地数据库
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : localhost:3306
 Source Schema         : yukoreimburse

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 18/01/2025 16:18:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for expense_reports
-- ----------------------------
DROP TABLE IF EXISTS `expense_reports`;
CREATE TABLE `expense_reports`  (
  `report_id` int NOT NULL AUTO_INCREMENT COMMENT '报销单 ID',
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '报销原因',
  `report_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '报销时间',
  `user_id` int NOT NULL COMMENT '用户 ID (外键)',
  `remarks` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `amount` decimal(10, 2) NOT NULL COMMENT '报销金额',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '报销单状态 (0: 待审批, 1: 通过, 2: 驳回, 3: 撤回)',
  `has_attachment` tinyint(1) NOT NULL DEFAULT 0 COMMENT '附件标志 (0: 无附件, 1: 有附件)',
  PRIMARY KEY (`report_id`) USING BTREE,
  INDEX `user_id`(`user_id` ASC) USING BTREE,
  CONSTRAINT `expense_reports_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 111 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '费用报销单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of expense_reports
-- ----------------------------
INSERT INTO `expense_reports` VALUES (1, '业务招待费', '2025-01-14 18:25:59', 1, '无备注', 500.00, 0, 0);
INSERT INTO `expense_reports` VALUES (52, '差旅费', '2025-01-15 10:48:26', 1, '国内出差费用', 1200.50, 0, 0);
INSERT INTO `expense_reports` VALUES (53, '办公用品采购', '2025-01-15 10:48:26', 2, '购买文具和打印纸', 300.00, 0, 0);
INSERT INTO `expense_reports` VALUES (54, '业务招待费', '2025-01-15 10:48:26', 3, '客户接待费用', 800.00, 2, 0);
INSERT INTO `expense_reports` VALUES (55, '差旅费', '2025-01-15 10:48:26', 4, '海外出差费用', 5200.75, 1, 0);
INSERT INTO `expense_reports` VALUES (56, '团队活动', '2025-01-15 10:48:26', 5, '部门聚餐', 1500.00, 1, 0);
INSERT INTO `expense_reports` VALUES (57, '办公设备维修', '2025-01-15 10:48:26', 1, '打印机维修费用', 200.00, 3, 0);
INSERT INTO `expense_reports` VALUES (58, '差旅费', '2025-01-15 10:48:26', 2, '会议差旅费用', 1000.00, 0, 0);
INSERT INTO `expense_reports` VALUES (59, '项目采购', '2025-01-15 10:48:26', 3, '购买服务器', 15000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (60, '快递费用', '2025-01-15 10:48:26', 4, '文件邮寄费用', 50.00, 2, 0);
INSERT INTO `expense_reports` VALUES (61, '培训费用', '2025-01-15 10:48:26', 5, '技能培训费用', 2000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (62, '业务招待费', '2025-01-15 10:48:26', 1, '商务洽谈费用', 1200.00, 1, 0);
INSERT INTO `expense_reports` VALUES (63, '差旅费', '2025-01-15 10:48:26', 2, '短期差旅费用', 900.00, 0, 0);
INSERT INTO `expense_reports` VALUES (64, '活动宣传', '2025-01-15 10:48:26', 3, '宣传物料费用', 3000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (65, '设备维护', '2025-01-15 10:48:26', 4, '办公电脑维修', 500.00, 3, 0);
INSERT INTO `expense_reports` VALUES (66, '业务招待费', '2025-01-15 10:48:26', 5, '合作客户接待', 600.00, 1, 0);
INSERT INTO `expense_reports` VALUES (67, '差旅费', '2025-01-15 10:48:26', 1, '跨省差旅费用', 1300.00, 0, 0);
INSERT INTO `expense_reports` VALUES (68, '办公用品采购', '2025-01-15 10:48:26', 2, '购买办公家具', 8000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (69, '广告投放', '2025-01-15 10:48:26', 3, '线上广告费用', 5000.00, 2, 0);
INSERT INTO `expense_reports` VALUES (70, '团队活动', '2025-01-15 10:48:26', 4, '团建活动', 2500.00, 0, 0);
INSERT INTO `expense_reports` VALUES (71, '差旅费', '2025-01-15 10:48:26', 5, '年度会议差旅', 1600.00, 1, 0);
INSERT INTO `expense_reports` VALUES (72, '会议室租赁', '2025-01-15 10:49:55', 1, '外部会议室租赁费用', 800.00, 0, 0);
INSERT INTO `expense_reports` VALUES (73, '快递费用', '2025-01-15 10:49:55', 2, '邮寄项目资料', 120.00, 1, 0);
INSERT INTO `expense_reports` VALUES (74, '业务招待费', '2025-01-15 10:49:55', 3, '客户来访餐饮费用', 950.00, 0, 0);
INSERT INTO `expense_reports` VALUES (75, '办公用品采购', '2025-01-15 10:49:55', 4, '购买文件夹和订书机', 150.00, 1, 0);
INSERT INTO `expense_reports` VALUES (76, '差旅费', '2025-01-15 10:49:55', 5, '短途高铁差旅', 500.00, 1, 0);
INSERT INTO `expense_reports` VALUES (77, '活动策划', '2025-01-15 10:49:55', 1, '公司周年庆活动费用', 3000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (78, '设备升级', '2025-01-15 10:49:55', 2, '办公电脑升级硬件', 4500.00, 2, 0);
INSERT INTO `expense_reports` VALUES (79, '广告投放', '2025-01-15 10:49:55', 3, '社交媒体广告费用', 4000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (80, '培训费用', '2025-01-15 10:49:55', 4, '线上课程订阅', 600.00, 0, 0);
INSERT INTO `expense_reports` VALUES (81, '团队活动', '2025-01-15 10:49:55', 5, '年度员工旅游', 10000.00, 3, 0);
INSERT INTO `expense_reports` VALUES (82, '业务招待费', '2025-01-15 10:49:55', 1, '与供应商洽谈费用', 1200.00, 0, 0);
INSERT INTO `expense_reports` VALUES (83, '差旅费', '2025-01-15 10:49:55', 2, '参加行业峰会费用', 2200.00, 1, 0);
INSERT INTO `expense_reports` VALUES (84, '办公用品采购', '2025-01-15 10:49:55', 3, '批量购买笔记本', 300.00, 2, 0);
INSERT INTO `expense_reports` VALUES (85, '设备维护', '2025-01-15 10:49:55', 4, '打印机清理保养', 100.00, 1, 0);
INSERT INTO `expense_reports` VALUES (86, '活动宣传', '2025-01-15 10:49:55', 5, '设计宣传海报', 800.00, 1, 0);
INSERT INTO `expense_reports` VALUES (87, '差旅费', '2025-01-15 10:49:55', 1, '国际航班费用', 12000.00, 0, 0);
INSERT INTO `expense_reports` VALUES (88, '广告投放', '2025-01-15 10:49:55', 2, 'SEO 优化服务', 7000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (89, '项目采购', '2025-01-15 10:49:55', 3, '购买专业相机', 2500.00, 0, 0);
INSERT INTO `expense_reports` VALUES (90, '业务招待费', '2025-01-15 10:49:55', 4, '客户见面餐饮', 1100.00, 2, 0);
INSERT INTO `expense_reports` VALUES (91, '团队建设', '2025-01-15 10:49:55', 5, '部门运动会', 1800.00, 1, 0);
INSERT INTO `expense_reports` VALUES (92, '差旅费', '2025-01-15 10:49:55', 1, '高铁商务座费用', 600.00, 3, 0);
INSERT INTO `expense_reports` VALUES (93, '活动策划', '2025-01-15 10:49:55', 2, '组织公益活动费用', 3500.00, 1, 0);
INSERT INTO `expense_reports` VALUES (94, '设备采购', '2025-01-15 10:49:55', 3, '更换老旧设备', 8000.00, 1, 0);
INSERT INTO `expense_reports` VALUES (95, '广告设计', '2025-01-15 10:49:55', 4, '设计线上宣传图', 1500.00, 0, 0);
INSERT INTO `expense_reports` VALUES (96, '会议安排', '2025-01-15 10:49:55', 5, '预定会议酒店', 2200.00, 1, 0);
INSERT INTO `expense_reports` VALUES (97, '快递费用', '2025-01-15 10:49:55', 1, '紧急文件快递', 60.00, 2, 0);
INSERT INTO `expense_reports` VALUES (98, '业务招待费', '2025-01-15 10:49:55', 2, '接待 VIP 客户', 3200.00, 2, 0);
INSERT INTO `expense_reports` VALUES (99, '差旅费', '2025-01-15 10:49:55', 3, '自驾差旅补贴', 800.00, 0, 0);
INSERT INTO `expense_reports` VALUES (100, '团队培训', '2025-01-15 10:49:55', 4, '技术分享会', 1200.00, 1, 0);
INSERT INTO `expense_reports` VALUES (101, '项目推进', '2025-01-15 10:49:55', 5, '项目合作费', 25000.00, 1, 1);
INSERT INTO `expense_reports` VALUES (102, '1111', '2025-01-16 18:26:21', 1, '12323', 2222.00, 0, 1);
INSERT INTO `expense_reports` VALUES (103, '1111', '2025-01-16 22:28:13', 1, '', 2222.00, 0, 1);
INSERT INTO `expense_reports` VALUES (104, '1111', '2025-01-16 18:30:08', 1, '', 2222.00, 0, 1);
INSERT INTO `expense_reports` VALUES (109, '1111', '2025-01-18 13:54:48', 1, '', 2222.00, 0, 0);
INSERT INTO `expense_reports` VALUES (110, '1111', '2025-01-18 13:54:48', 1, '', 2222.00, 0, 0);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `user_id` int NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `is_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否为管理员 (0: 否, 1: 是)',
  `lark_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Lark ID',
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', 1, 'test_admin');
INSERT INTO `users` VALUES (2, 'user1', '7c6a180b36896a0a8c02787eeafb0e4c', 0, 'lark_user_1');
INSERT INTO `users` VALUES (3, 'user2', '6cb75f652a9b52798eb6cf2201057c73', 0, 'lark_user_2');
INSERT INTO `users` VALUES (4, 'user3', '819b0643d6b89dc9b579fdfc9094f28e', 0, 'lark_user_3');
INSERT INTO `users` VALUES (5, '叶俊港', md5('6987058as'), 0, 'ou_a880482bf4fd32e04d0acda1d9e200ee');
INSERT INTO `users` VALUES (4, '苏佳轩', md5(sjx57529966@), 0, 'ou_01bcfb608879a16702ece6f9afbaee99');
SET FOREIGN_KEY_CHECKS = 1;
INSERT INTO `users` VALUES (3, 'yori', md5('Ak668899'), 0, 'ou_361a84bb42e1210aafb90cc262e6ef77');
INSERT INTO `users` VALUES (6, 'chingboa', md5('960410'), 0, 'ou_fd52f2fb10722282e0b98ca717a9ea74');
INSERT INTO `users` VALUES (7, 'linlili', md5('917728'), 0, 'ou_602e8edb480ee2190a4afff5910a777b');
INSERT INTO `users` VALUES (8, '张惠倩', md5('1314520h&'), 0, 'ou_fd538d83446f250f81d8417b01a18192');