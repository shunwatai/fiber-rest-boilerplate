CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"first_name"	TEXT,
	"last_name"	TEXT,
	"disabled"	NUMERIC DEFAULT '0',
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT)
);

INSERT INTO `users` (`name`, `password`, `email`, `first_name`, `last_name`, `disabled`, `created_at`, `updated_at`) VALUES ('admin',	'$2a$04$Ey.Y3FdhY5jjrdKQsTxCYOU2jieFRgZZCjM3P2yXivrj.Zmk0G3BS', 'admin@example.com',	NULL,	NULL,	0,	'2024-02-24 17:47:29',	'2024-02-25 08:36:21');
