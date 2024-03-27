-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create password_resets table
DROP SEQUENCE IF EXISTS password_resets_id_seq;
CREATE SEQUENCE password_resets_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."password_resets" (
    "id" integer DEFAULT nextval('password_resets_id_seq') NOT NULL,
    "user_id" integer NOT NULL,
    "token_hash" text NOT NULL,
    "expiry_date" timestamptz NOT NULL,
    "is_used" boolean NOT NULL DEFAULT false,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "password_resets_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
COMMENT ON COLUMN "password_resets"."is_used" IS 'mark as true after password reset';

CREATE INDEX "password_resets_user_id" ON "password_resets" ("user_id");

ALTER TABLE "password_resets"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_password_resets_updated_at BEFORE UPDATE ON password_resets FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
