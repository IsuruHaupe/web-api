package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/require"
)

func CreateRandomContact(t *testing.T) Contact {
	args := CreateContactParams{
		Firstname:   randomdata.FirstName(randomdata.Female),
		Lastname:    randomdata.LastName(),
		Fullname:    randomdata.FullName(randomdata.Female),
		HomeAddress: randomdata.Address(),
		Email:       randomdata.Email(),
		PhoneNumber: randomdata.PhoneNumber(),
	}

	contact, err := testQueries.CreateContact(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	require.Equal(t, args.Firstname, contact.Firstname)
	require.Equal(t, args.Lastname, contact.Lastname)
	require.Equal(t, args.Fullname, contact.Fullname)
	require.Equal(t, args.HomeAddress, contact.HomeAddress)
	require.Equal(t, args.Email, contact.Email)
	require.Equal(t, args.PhoneNumber, contact.PhoneNumber)

	require.NotZero(t, contact.ID)

	return contact
}
func TestCreateContact(t *testing.T) {
	CreateRandomContact(t)
}

func TestGetContact(t *testing.T) {
	contact1 := CreateRandomContact(t)
	contact2, err := testQueries.GetContact(context.Background(), contact1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ID, contact2.ID)
	require.Equal(t, contact1.Firstname, contact2.Firstname)
	require.Equal(t, contact1.Lastname, contact2.Lastname)
	require.Equal(t, contact1.Fullname, contact2.Fullname)
	require.Equal(t, contact1.HomeAddress, contact2.HomeAddress)
	require.Equal(t, contact1.Email, contact2.Email)
	require.Equal(t, contact1.PhoneNumber, contact2.PhoneNumber)
}

func TestUpdateContact(t *testing.T) {
	contact1 := CreateRandomContact(t)

	args := UpdateContactParams{
		ID:          contact1.ID,
		Firstname:   randomdata.FirstName(randomdata.Female),
		Lastname:    randomdata.LastName(),
		Fullname:    randomdata.FullName(randomdata.Female),
		HomeAddress: randomdata.Address(),
		Email:       randomdata.Email(),
		PhoneNumber: randomdata.PhoneNumber(),
	}

	contact2, err := testQueries.UpdateContact(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, args.ID, contact2.ID)
	require.Equal(t, args.Firstname, contact2.Firstname)
	require.Equal(t, args.Lastname, contact2.Lastname)
	require.Equal(t, args.Fullname, contact2.Fullname)
	require.Equal(t, args.HomeAddress, contact2.HomeAddress)
	require.Equal(t, args.Email, contact2.Email)
	require.Equal(t, args.PhoneNumber, contact2.PhoneNumber)
}

func TestDeleteContact(t *testing.T) {
	contact1 := CreateRandomContact(t)
	err := testQueries.DeleteContact(context.Background(), contact1.ID)
	require.NoError(t, err)

	contact2, err := testQueries.GetContact(context.Background(), contact1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, contact2)
}

func TestListContacts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomContact(t)
	}

	args := ListContactsParams{
		Limit:  5,
		Offset: 1,
	}

	contacts, err := testQueries.ListContacts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, contacts, 5)

	for _, contact := range contacts {
		require.NotEmpty(t, contact)
	}
}

func TestGetContactsWithSkillAndLevel(t *testing.T) {
	// Create a random contact.
	originalContact := CreateRandomContact(t)
	// Create a random skill.
	skill := CreateRandomSkill(t)

	// Bind a contact with a skill
	argsContactHasSkill := CreateContactHasSkillParams{
		ContactID: int32(originalContact.ID),
		SkillID:   int32(skill.ID),
	}

	contactHasSkill, err := testQueries.CreateContactHasSkill(context.Background(), argsContactHasSkill)
	require.NoError(t, err)
	require.NotEmpty(t, contactHasSkill)

	argsContactsWithSkillAndLevelParams := GetContactsWithSkillAndLevelParams{
		SkillName:  skill.SkillName,
		SkillLevel: skill.SkillLevel,
	}

	contacts, err := testQueries.GetContactsWithSkillAndLevel(context.Background(), argsContactsWithSkillAndLevelParams)
	require.NoError(t, err)
	require.Len(t, contacts, 1)

	for _, contact := range contacts {
		require.NotEmpty(t, contact)

		require.Equal(t, originalContact.ID, contact.ID)
		require.Equal(t, originalContact.Firstname, contact.Firstname)
		require.Equal(t, originalContact.Lastname, contact.Lastname)
		require.Equal(t, originalContact.Fullname, contact.Fullname)
		require.Equal(t, originalContact.HomeAddress, contact.HomeAddress)
		require.Equal(t, originalContact.Email, contact.Email)
		require.Equal(t, originalContact.PhoneNumber, contact.PhoneNumber)
	}

	err = testQueries.DeleteSkill(context.Background(), skill.ID)
	require.NoError(t, err)
}

func TestGetContactsWithSkill(t *testing.T) {
	// Create a random contact.
	originalContact := CreateRandomContact(t)
	// Create a random skill.
	skill := CreateRandomSkill(t)

	// Bind a contact with a skill
	argsContactHasSkill := CreateContactHasSkillParams{
		ContactID: int32(originalContact.ID),
		SkillID:   int32(skill.ID),
	}

	contactHasSkill, err := testQueries.CreateContactHasSkill(context.Background(), argsContactHasSkill)
	require.NoError(t, err)
	require.NotEmpty(t, contactHasSkill)

	contacts, err := testQueries.GetContactsWithSkill(context.Background(), skill.SkillName)
	require.NoError(t, err)
	require.Len(t, contacts, 1)

	for _, contact := range contacts {
		require.NotEmpty(t, contact)

		require.Equal(t, originalContact.ID, contact.ID)
		require.Equal(t, originalContact.Firstname, contact.Firstname)
		require.Equal(t, originalContact.Lastname, contact.Lastname)
		require.Equal(t, originalContact.Fullname, contact.Fullname)
		require.Equal(t, originalContact.HomeAddress, contact.HomeAddress)
		require.Equal(t, originalContact.Email, contact.Email)
		require.Equal(t, originalContact.PhoneNumber, contact.PhoneNumber)
	}

	err = testQueries.DeleteSkill(context.Background(), skill.ID)
	require.NoError(t, err)
}
