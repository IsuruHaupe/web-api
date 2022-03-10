package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateContactHasSkill(t *testing.T) {
	user := CreateRandomUser(t)
	// Create a contact.
	contact := CreateRandomContact(t, user)

	// Create a skill.
	skill := CreateRandomSkill(t, user)

	// Bind a contact with a skill
	argsContactHasSkill := CreateContactHasSkillParams{
		Owner:     user.Username,
		ContactID: int32(contact.ID),
		SkillID:   int32(skill.ID),
	}

	contactHasSkill, err := testQueries.CreateContactHasSkill(context.Background(), argsContactHasSkill)
	require.NoError(t, err)
	require.NotEmpty(t, contactHasSkill)

	err = testQueries.DeleteSkill(context.Background(), skill.ID)
	require.NoError(t, err)
}
