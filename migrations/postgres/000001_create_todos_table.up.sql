-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create todos table
DROP SEQUENCE IF EXISTS todos_id_seq;
CREATE SEQUENCE todos_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."todos" (
    "id" integer DEFAULT nextval('todos_id_seq') NOT NULL,
    "task" character varying(255) NOT NULL,
    "done" boolean DEFAULT false NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "todos_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_todos_updated_at BEFORE UPDATE ON todos FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
