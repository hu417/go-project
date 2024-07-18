CREATE TABLE `user_basic` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`name` varchar(60) DEFAULT NULL,
	`password` varchar(32) DEFAULT NULL,
	`email` varchar(100) DEFAULT NULL,

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;