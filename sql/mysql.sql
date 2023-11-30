-- ----------------------------
-- Records of mk_taskstatus
-- ----------------------------
INSERT INTO `mk_task_status` VALUES ('1', 'No');
INSERT INTO `mk_task_status` VALUES ('2', 'Off');
INSERT INTO `mk_task_status` VALUES ('3', 'Run');

ALTER TABLE `mk_task`
CHANGE COLUMN `campaign_id_id` `campaign_id`  bigint(20) NOT NULL AFTER `task_status_id`;

-- INSERT INTO `mk_account_role` (id,name) VALUES ('1','admin')