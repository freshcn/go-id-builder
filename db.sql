--
-- Table structure for table `idgenerator`
--

DROP TABLE IF EXISTS `idgenerator`;
CREATE TABLE `idgenerator` (
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT 'id名',
  `id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '当前的最大id',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
LOCK TABLES `idgenerator` WRITE;
INSERT INTO `idgenerator` VALUES ('test',0,'',0);
UNLOCK TABLES;
