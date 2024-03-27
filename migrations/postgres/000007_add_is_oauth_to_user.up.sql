ALTER TABLE "users"
ADD "is_oauth" boolean NOT NULL DEFAULT false,
ADD "provider" character varying(255) NULL;
COMMENT ON COLUMN "users"."is_oauth" IS 'flag for oauth user';
COMMENT ON COLUMN "users"."provider" IS 'mark the oauth provider';
COMMENT ON TABLE "users" IS '';
