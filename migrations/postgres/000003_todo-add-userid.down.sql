DROP INDEX "todos_user_id";

ALTER TABLE "todos"
DROP CONSTRAINT "todos_user_id_fkey";

ALTER TABLE "todos"
DROP "user_id";
