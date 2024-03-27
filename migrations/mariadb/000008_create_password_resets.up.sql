CREATE TABLE `password_resets` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` int(11) NOT NULL,
  `token_hash` text NOT NULL,
  `expiry_date` datetime NOT NULL,
  `is_used` tinyint NOT NULL DEFAULT '0' COMMENT 'mark as true after password reset',
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
