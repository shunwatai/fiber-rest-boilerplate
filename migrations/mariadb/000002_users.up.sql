CREATE TABLE IF NOT EXISTS  `users` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `first_name` varchar(255) NULL,
  `last_name` varchar(255) NULL,
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO `users` (`id`, `name`, `password`, `first_name`, `last_name`, `disabled`, `created_at`, `updated_at`) VALUES (1,	'admin',	'$2a$04$Ey.Y3FdhY5jjrdKQsTxCYOU2jieFRgZZCjM3P2yXivrj.Zmk0G3BS',	NULL,	NULL,	0,	'2024-02-24 17:47:29',	'2024-02-25 08:36:21');
