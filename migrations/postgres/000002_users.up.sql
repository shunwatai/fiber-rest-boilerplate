-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create users table
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "name" character varying(255) NOT NULL,
    "password" character varying(255) NOT NULL,
    "email" character varying(255) NOT NULL,
    "first_name" character varying(255) DEFAULT NULL,
    "last_name" character varying(255) DEFAULT NULL,
    "disabled" boolean DEFAULT false NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

ALTER TABLE "users"
ADD CONSTRAINT "users_name" UNIQUE ("name"),
ADD CONSTRAINT "users_email" UNIQUE ("email");

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

INSERT INTO "users" ("name", "password", "email", "first_name", "last_name", "disabled", "created_at", "updated_at") VALUES ('admin',	'$2a$04$7F9KIfLOW3O9LyZSm2IQ8uXqH0B7P3wLjYTlkaerX53muN4U1.FDq', 'admin@example.com',	NULL,	NULL,	'0',	'2023-11-27 10:35:53+00',	'2024-02-25 08:32:51.828793+00');
