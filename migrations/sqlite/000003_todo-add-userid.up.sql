-- https://stackoverflow.com/a/67344658
CREATE TEMPORARY TABLE temp_todos AS
SELECT 
    id,
    task,
    done,
    created_at,
    updated_at
FROM todos;

DROP TABLE todos;

CREATE TABLE "todos" (
	"id"	INTEGER NOT NULL,
	"user_id"	INTEGER,
	"task"	TEXT,
	"done"	NUMERIC,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);

INSERT INTO todos
 (  id,
    task,
    done,
    created_at,
    updated_at
  )
SELECT
    id,
    task,
    done,
    created_at,
    updated_at
FROM temp_todos;
