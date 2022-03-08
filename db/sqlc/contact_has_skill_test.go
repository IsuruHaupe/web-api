package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateContactHasSkill(t *testing.T) {
	// Create a contact.
	contact := CreateRandomContact(t)

	// Create a skill.
	skill := CreateRandomSkill(t)

	// Bind a contact with a skill
	argsContactHasSkill := CreateContactHasSkillParams{
		ContactID: int32(contact.ID),
		SkillID:   int32(skill.ID),
	}

	contactHasSkill, err := testQueries.CreateContactHasSkill(context.Background(), argsContactHasSkill)
	require.NoError(t, err)
	require.NotEmpty(t, contactHasSkill)

	err = testQueries.DeleteSkill(context.Background(), skill.ID)
	require.NoError(t, err)
}
