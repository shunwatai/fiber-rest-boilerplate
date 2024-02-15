CREATE TABLE IF NOT EXISTS "documents" (
  "id"	INTEGER NOT NULL,
	"user_id"	INTEGER,
	"name"	TEXT NOT NULL,
	"file_path"	TEXT NOT NULL,
	"file_type"	TEXT NOT NULL,
	"file_size"	NUMERIC NOT NULL,
	"hash"	TEXT NOT NULL,
	"public"	NUMERIC DEFAULT '1',
	"created_at"	DATETIME NOT NULL,
	"updated_at"	DATETIME NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
