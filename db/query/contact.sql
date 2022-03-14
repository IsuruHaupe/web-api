-- name: CreateContact :one
INSERT INTO contacts (
  owner,
  firstname, 
  lastname, 
  fullname, 
  home_address, 
  email, 
  phone_number
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
WHERE owner = $1
ORDER BY firstname ASC
LIMIT $2
OFFSET $3;

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

-- name: GetContactSkills :many
SELECT * FROM skills
WHERE id IN (
    SELECT skill_id
    FROM contact_has_skill
    WHERE contact_id = $1
);

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

-- name: GetIfExistsContactID :one 
SELECT EXISTS (SELECT * FROM contacts WHERE id = $1);

-- name: GetFirstname :one
SELECT firstname FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetLastname :one
SELECT lastname FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetFullname :one
SELECT fullname FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetHomeAddress :one
SELECT home_address FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetEmail :one
SELECT email FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetPhoneNumber :one
SELECT phone_number FROM contacts
WHERE id = $1 LIMIT 1;
