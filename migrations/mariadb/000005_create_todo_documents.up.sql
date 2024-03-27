CREATE TABLE `todo_documents` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `todo_id` int(11) NOT NULL,
  `document_id` int(11) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp,
  `updated_at` datetime NOT NULL DEFAULT current_timestamp,
  FOREIGN KEY (`todo_id`) REFERENCES `todos` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`document_id`) REFERENCES `documents` (`id`) ON DELETE CASCADE
);
