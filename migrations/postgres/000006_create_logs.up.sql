-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create logs table
DROP SEQUENCE IF EXISTS logs_id_seq;
CREATE SEQUENCE logs_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."logs" (
    "id" integer DEFAULT nextval('logs_id_seq') NOT NULL,
    "user_id" integer NULL,
    "ip_address" inet NOT NULL,
    "http_method" character varying(255) NOT NULL,
    "route" text NOT NULL,
    "user_agent" character varying(255) NOT NULL,
    "request_header" json NULL,
    "request_body" json NULL,
    "response_body" json NULL,
    "status" integer NOT NULL,
    "duration" integer NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "logs_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
COMMENT ON COLUMN "logs"."user_id" IS 'from jwt token';
COMMENT ON COLUMN "logs"."http_method" IS 'get, post, patch, delete etc.';
COMMENT ON COLUMN "logs"."route" IS 'api endpoint';
COMMENT ON COLUMN "logs"."status" IS 'http status code';
COMMENT ON COLUMN "logs"."duration" IS 'time in ms';

-- Following is a sample for adding foreign key user_id to users table
ALTER TABLE "logs"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_logs_updated_at BEFORE UPDATE ON logs FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
