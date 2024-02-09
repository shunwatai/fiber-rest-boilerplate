CREATE TABLE `documents` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` int(11) NULL,
  `name` varchar(255) NOT NULL,
  `file_path` text NOT NULL,
  `file_type` varchar(255) NOT NULL,
  `file_size` int NOT NULL,
  `hash` text NOT NULL,
  `public` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
);
