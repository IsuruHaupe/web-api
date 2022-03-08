// Code generated by sqlc. DO NOT EDIT.
// source: contact.sql

package db

import (
	"context"
)

const createContact = `-- name: CreateContact :one
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
RETURNING id, firstname, lastname, fullname, home_address, email, phone_number
`

type CreateContactParams struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Fullname    string `json:"fullname"`
	HomeAddress string `json:"home_address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (q *Queries) CreateContact(ctx context.Context, arg CreateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, createContact,
		arg.Firstname,
		arg.Lastname,
		arg.Fullname,
		arg.HomeAddress,
		arg.Email,
		arg.PhoneNumber,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Fullname,
		&i.HomeAddress,
		&i.Email,
		&i.PhoneNumber,
	)
	return i, err
}

const deleteContact = `-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1
`

func (q *Queries) DeleteContact(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteContact, id)
	return err
}

const getContact = `-- name: GetContact :one
SELECT id, firstname, lastname, fullname, home_address, email, phone_number FROM contacts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetContact(ctx context.Context, id int64) (Contact, error) {
	row := q.db.QueryRowContext(ctx, getContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Fullname,
		&i.HomeAddress,
		&i.Email,
		&i.PhoneNumber,
	)
	return i, err
}

const getContactsWithSkill = `-- name: GetContactsWithSkill :many
SELECT id, firstname, lastname, fullname, home_address, email, phone_number FROM contacts
WHERE id IN (
    SELECT contact_id
    FROM contact_has_skill
    WHERE skill_id IN (
    	SELECT id 
    	FROM skills 
    	WHERE skill_name = $1
    )
  )
`

func (q *Queries) GetContactsWithSkill(ctx context.Context, skillName string) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, getContactsWithSkill, skillName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Fullname,
			&i.HomeAddress,
			&i.Email,
			&i.PhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getContactsWithSkillAndLevel = `-- name: GetContactsWithSkillAndLevel :many
SELECT id, firstname, lastname, fullname, home_address, email, phone_number FROM contacts
WHERE id IN (
    SELECT contact_id
    FROM contact_has_skill
    WHERE skill_id IN (
    	SELECT id 
    	FROM skills 
    	WHERE skill_name = $1 AND skill_level = $2
    )
  )
`

type GetContactsWithSkillAndLevelParams struct {
	SkillName  string `json:"skill_name"`
	SkillLevel string `json:"skill_level"`
}

func (q *Queries) GetContactsWithSkillAndLevel(ctx context.Context, arg GetContactsWithSkillAndLevelParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, getContactsWithSkillAndLevel, arg.SkillName, arg.SkillLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Fullname,
			&i.HomeAddress,
			&i.Email,
			&i.PhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listContacts = `-- name: ListContacts :many
SELECT id, firstname, lastname, fullname, home_address, email, phone_number FROM contacts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListContactsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListContacts(ctx context.Context, arg ListContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listContacts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Fullname,
			&i.HomeAddress,
			&i.Email,
			&i.PhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateContact = `-- name: UpdateContact :one
UPDATE contacts 
SET firstname = $2, 
lastname = $3, 
fullname = $4, 
home_address = $5, 
email = $6, 
phone_number = $7
WHERE id = $1
RETURNING id, firstname, lastname, fullname, home_address, email, phone_number
`

type UpdateContactParams struct {
	ID          int64  `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Fullname    string `json:"fullname"`
	HomeAddress string `json:"home_address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (q *Queries) UpdateContact(ctx context.Context, arg UpdateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, updateContact,
		arg.ID,
		arg.Firstname,
		arg.Lastname,
		arg.Fullname,
		arg.HomeAddress,
		arg.Email,
		arg.PhoneNumber,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Fullname,
		&i.HomeAddress,
		&i.Email,
		&i.PhoneNumber,
	)
	return i, err
}
