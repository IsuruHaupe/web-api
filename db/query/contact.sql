-- name: CreateContact :one
INSERT INTO contacts (
  firstname, 
  lastname, 
  fullname, 
  home_address, 
  email, 
  phone_number
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateContact :one
UPDATE contacts 
SET firstname = $2, 
lastname = $3, 
fullname = $4, 
home_address = $5, 
email = $6, 
phone_number = $7
WHERE id = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;

-- name: GetContactsWithSkillAndLevel :many
SELECT * FROM contacts
WHERE id IN (
    SELECT contact_id
    FROM contact_has_skill
    WHERE skill_id IN (
    	SELECT id 
    	FROM skills 
    	WHERE skill_name = $1 AND skill_level = $2
    )
  );

  -- name: GetContactsWithSkill :many
SELECT * FROM contacts
WHERE id IN (
    SELECT contact_id
    FROM contact_has_skill
    WHERE skill_id IN (
    	SELECT id 
    	FROM skills 
    	WHERE skill_name = $1
    )
  );