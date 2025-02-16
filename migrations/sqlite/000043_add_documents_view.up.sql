DROP VIEW IF EXISTS documents_view;
CREATE OR REPLACE VIEW documents_view AS 
select 
  documents.id,
  documents.user_id,
  users.name user,
  documents.name,
  documents.file_path,
  documents.file_type,
  documents.file_size,
  documents.hash,
  documents.public,
  documents.created_at,
  documents.updated_at
from documents
left join users on users.id = documents.user_id
order by id desc;
