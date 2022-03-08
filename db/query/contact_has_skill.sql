-- name: CreateContactHasSkill :one
INSERT INTO contact_has_skill (
  contact_id, 
  skill_id
) VALUES (
  $1, $2
)
RETURNING *;