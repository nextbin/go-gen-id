-- 仅当使用到 http_gin 服务的白名单功能需要
CREATE TABLE `nb_gen_id_whitelist` (
  `id` int(11) NOT NULL,
  `ip` varchar(16) NOT NULL,
  `create_time` datetime NOT NULL,
  `status` int(11) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_IP` (`ip`),
  KEY `IDX_CRT` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='白名单';

-- 仅当 config.go:CheckMachineIdType int = CheckMachineIdTypeMysql 时需要使用MySQL
CREATE TABLE `nb_gen_id_machine` (
  `id` int(11) NOT NULL,
  `ip` varchar(16) NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_IP` (`ip`),
  KEY `IDX_CRT` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='机器列表';