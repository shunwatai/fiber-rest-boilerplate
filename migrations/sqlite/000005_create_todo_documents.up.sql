CREATE TABLE IF NOT EXISTS "todo_documents" (
  "id"	INTEGER NOT NULL,
	"todo_id"	INTEGER NOT NULL,
	"document_id"	INTEGER NOT NULL,
	"created_at"	DATETIME NOT NULL,
	"updated_at"	DATETIME NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("todo_id") REFERENCES "todo"("id"),
	FOREIGN KEY("document_id") REFERENCES "document"("id")
);
