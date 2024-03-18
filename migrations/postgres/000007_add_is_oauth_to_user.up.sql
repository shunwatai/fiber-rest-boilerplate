ALTER TABLE "users"
ADD "is_oauth" boolean NOT NULL DEFAULT false;
COMMENT ON COLUMN "users"."is_oauth" IS 'flag for oauth user';
COMMENT ON TABLE "users" IS '';
