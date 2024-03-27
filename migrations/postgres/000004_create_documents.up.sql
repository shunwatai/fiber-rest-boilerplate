-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create documents table
DROP SEQUENCE IF EXISTS documents_id_seq;
CREATE SEQUENCE documents_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."documents" (
    "id" integer DEFAULT nextval('documents_id_seq') NOT NULL,
    "user_id" integer NULL,
    "name" character varying(255) NOT NULL,
    "file_path" text NOT NULL,
    "file_type" character varying(255) NOT NULL,
    "file_size" integer NOT NULL,
    "hash" text NOT NULL,
    "public" boolean DEFAULT true NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "documents_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "documents_user_id" ON "documents" ("user_id");
ALTER TABLE "documents"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL;

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON documents FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
