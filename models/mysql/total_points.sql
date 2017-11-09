CREATE TABLE `total_points` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `total` int(20) NOT NULL DEFAULT '0' COMMENT '用户总积分',
  `update_time` datetime NOT NULL DEFAULT '2006-01-02 15:04:05' COMMENT '总积分更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COMMENT='用户总积分表';