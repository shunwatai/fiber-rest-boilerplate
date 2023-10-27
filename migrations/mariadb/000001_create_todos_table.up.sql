CREATE TABLE IF NOT EXISTS  `todos` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `task` varchar(255) NOT NULL,
  `done` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp ON UPDATE CURRENT_TIMESTAMP
);
