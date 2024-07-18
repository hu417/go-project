CREATE TABLE `user_repository` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`parent_id` int(11) DEFAULT NULL COMMENT '父级文件层级, 0-【文件夹】',
	`user_identity` varchar(36) DEFAULT NULL COMMENT '对应用户的唯一标识',
	`repository_identity` varchar(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识',
	`ext` varchar(255) DEFAULT NULL COMMENT '文件或文件夹类型',
	`name` varchar(255) DEFAULT NULL COMMENT '用户定义的文件名',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;