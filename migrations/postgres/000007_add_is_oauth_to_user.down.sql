ALTER TABLE "users"
DROP "is_oauth",
DROP "provider";
COMMENT ON TABLE "users" IS '';
