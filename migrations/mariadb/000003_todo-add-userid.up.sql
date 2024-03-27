ALTER TABLE `todos`
ADD `user_id` int(11) NULL AFTER `id`,
CHANGE `updated_at` `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`,
ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;
