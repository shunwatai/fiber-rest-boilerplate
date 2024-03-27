ALTER TABLE `todos` DROP FOREIGN KEY `todos_ibfk_1`,
DROP INDEX `user_id`,
DROP `user_id`,
CHANGE `updated_at` `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`;
