CREATE TABLE `logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` int(11) NULL COMMENT 'from jwt token',
  `ip_address` varchar(255) NOT NULL,
  `http_method` varchar(255) NOT NULL COMMENT 'get, post, patch, delete etc.',
  `route` varchar(255) NOT NULL COMMENT 'api endpoint',
  `user_agent` varchar(255) NOT NULL,
  `request_header` json NULL,
  `request_body` json NULL,
  `response_body` json NULL,
  `status` int NOT NULL COMMENT 'http status code',
  `duration` int NOT NULL COMMENT 'time in ms',
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
);
