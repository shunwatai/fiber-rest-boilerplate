ALTER TABLE "logs"
DROP CONSTRAINT "logs_user_id_fkey",
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;