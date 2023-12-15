ALTER TABLE `todos`
DROP FOREIGN KEY `todos_ibfk_1`;

ALTER TABLE `todos`
DROP INDEX `user_id`;

ALTER TABLE `todos`
DROP `user_id`,
CHANGE `updated_at` `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`;
