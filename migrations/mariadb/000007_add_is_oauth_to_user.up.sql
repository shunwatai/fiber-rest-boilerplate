ALTER TABLE `users`
ADD `is_oauth` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'flag for oauth user' AFTER `disabled`,
ADD `provider` varchar(255) NULL COMMENT 'mark the oauth provider' AFTER `is_oauth`,
CHANGE `updated_at` `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`;
