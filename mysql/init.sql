INSERT INTO `mk_taskstatus` VALUES ('1', 'No');
INSERT INTO `mk_taskstatus` VALUES ('2', 'Off');
INSERT INTO `mk_taskstatus` VALUES ('3', 'Start');
INSERT INTO `mk_task_status` VALUES ('4', 'Complete');

INSERT INTO `mk_social_platform` VALUES ('1', 'facebook');
INSERT INTO `mk_social_platform` VALUES ('2', 'youtube');

ALTER TABLE `account` ADD COLUMN `roles` VARCHAR(255) DEFAULT NULL AFTER `password`;

ALTER TABLE mk_campaign DROP COLUMN types;