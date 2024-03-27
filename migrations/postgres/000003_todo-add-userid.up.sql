ALTER TABLE "todos"
ADD "user_id" integer NULL;

ALTER TABLE "todos"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL;

CREATE INDEX "todos_user_id" ON "todos" ("user_id");
