CREATE TABLE IF NOT EXISTS  `users` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `first_name` varchar(255) NULL,
  `last_name` varchar(255) NULL,
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE `users`
ADD UNIQUE `name` (`name`),
ADD UNIQUE `email` (`email`);

INSERT INTO `users` (`name`, `password`, `email`, `first_name`, `last_name`, `disabled`) 
VALUES ('admin',	'$2a$04$Ey.Y3FdhY5jjrdKQsTxCYOU2jieFRgZZCjM3P2yXivrj.Zmk0G3BS', 'admin@example.com',	NULL,	NULL,	0);
