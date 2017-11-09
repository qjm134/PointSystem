CREATE TABLE `point_records_table_0` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `order_id` varchar(20) NOT NULL DEFAULT '0' COMMENT '订单ID',
  `operation` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '更新积分操作（0 减少，1 增加)',
  `points` int(11) NOT NULL DEFAULT '0' COMMENT '增加或减少的积分',
  `record_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '积分记录时间 Unix时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_userid_orderid` (`user_id`,`order_id`),
  KEY `idx_userid_recordtime` (`user_id`,`record_time`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COMMENT='用户积分流水表';