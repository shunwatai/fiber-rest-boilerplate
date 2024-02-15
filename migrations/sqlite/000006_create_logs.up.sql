CREATE TABLE "logs" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"ip_address"	TEXT NOT NULL,
	"http_method"	TEXT NOT NULL,
	"route"	TEXT NOT NULL,
	"user_agent"	TEXT NOT NULL,
	"request_header"	TEXT,
	"request_body"	TEXT,
	"response_body"	TEXT,
	"status"	INTEGER NOT NULL,
	"duration"	INTEGER NOT NULL,
	"created_at"	DATETIME NOT NULL,
	"updated_at"	DATETIME NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
