CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"password"	TEXT NOT NULL,
	"first_name"	TEXT,
	"last_name"	TEXT,
	"disabled"	NUMERIC DEFAULT '0',
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT)
);
