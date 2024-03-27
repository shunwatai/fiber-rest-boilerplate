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
	"task"	TEXT,
	"done"	NUMERIC,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT)
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
