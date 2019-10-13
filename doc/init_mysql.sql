-- 仅当 config.go:CheckMachineIdType int = CheckMachineIdTypeMysql 时需要使用MySQL
CREATE TABLE `nb_id_gen_machine` (
  `id` int(11) NOT NULL,
  `ip` varchar(16) NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_IP` (`ip`),
  KEY `IDX_CRT` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='机器列表';