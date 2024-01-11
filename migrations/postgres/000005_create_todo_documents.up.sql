-- Create TRIGGER function for auto update updated_at column after records being altered
CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create todo_documents table
DROP SEQUENCE IF EXISTS todo_documents_id_seq;
CREATE SEQUENCE todo_documents_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."todo_documents" (
    "id" integer DEFAULT nextval('todo_documents_id_seq') NOT NULL,
    "todo_id" integer NOT NULL,
    "document_id" integer NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "todo_documents_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "todo_documents_todo_id" ON "todo_documents" ("todo_id");
CREATE INDEX "todo_documents_document_id" ON "todo_documents" ("document_id");

ALTER TABLE "todo_documents"
ADD FOREIGN KEY ("todo_id") REFERENCES "todos" ("id") ON DELETE CASCADE;
ALTER TABLE "todo_documents"
ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id") ON DELETE CASCADE;

-- Apply trigger update_updated_at_column()
CREATE TRIGGER update_todo_documents_updated_at BEFORE UPDATE ON todo_documents FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
