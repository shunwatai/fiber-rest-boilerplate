CREATE TABLE "password_resets" (
	"id"	INTEGER NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"token_hash"	TEXT NOT NULL,
	"expiry_date"	DATETIME NOT NULL,
	"is_used"	NUMERIC NOT NULL DEFAULT 0,
	"created_at"	DATETIME NOT NULL,
	"updated_at"	DATETIME NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);
