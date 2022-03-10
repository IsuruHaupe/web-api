-- name: CreateSkill :one
INSERT INTO skills (
  owner,
  skill_name, 
  skill_level
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetIfExistsSkillID :one 
SELECT EXISTS (SELECT * FROM skills WHERE id = $1);

-- name: GetSkillName :one
SELECT skill_name FROM skills
WHERE id = $1 LIMIT 1;

-- name: GetSkillLevel :one
SELECT skill_level FROM skills
WHERE id = $1 LIMIT 1;


-- name: GetSkill :one
SELECT * FROM skills
WHERE id = $1 LIMIT 1;

-- name: ListSkills :many
SELECT * FROM skills
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateSkill :one
UPDATE skills 
SET skill_name = $2, 
skill_level = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSkill :exec
DELETE FROM skills WHERE id = $1;