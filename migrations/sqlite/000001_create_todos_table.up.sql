CREATE TABLE IF NOT EXISTS "todos" (
	"id"	INTEGER NOT NULL,
	"task"	TEXT,
	"done"	NUMERIC,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT)
);
