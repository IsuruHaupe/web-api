package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var randomProgLang = [...]string{"Go", "Java", "Javascript", "C++", "Python", "R", "HTML"}
var randomProfLvl = [...]string{"Familiar", "Proficient", "Excellent", "Expert"}

func CreateRandomSkill(t *testing.T, user User) Skill {
	rand.Seed(time.Now().UnixNano())
	args := CreateSkillParams{
		Owner:      user.Username,
		SkillName:  randomProgLang[rand.Intn(len(randomProgLang))],
		SkillLevel: randomProfLvl[rand.Intn(len(randomProfLvl))],
	}

	skill, err := testQueries.CreateSkill(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, skill)

	require.Equal(t, args.SkillName, skill.SkillName)
	require.Equal(t, args.SkillLevel, skill.SkillLevel)

	require.NotZero(t, skill.ID)
	return skill
}

func TestCreateSkill(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomSkill(t, user)

	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}
func TestDeleteSkill(t *testing.T) {
	user := CreateRandomUser(t)
	skill1 := CreateRandomSkill(t, user)
	err := testQueries.DeleteSkill(context.Background(), skill1.ID)
	require.NoError(t, err)

	skill2, err := testQueries.GetSkill(context.Background(), skill1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, skill2)

	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

func TestUniqueConstraintOnSkill(t *testing.T) {
	user := CreateRandomUser(t)
	args := CreateSkillParams{
		Owner:      user.Username,
		SkillName:  "same",
		SkillLevel: "same",
	}
	skill1, err := testQueries.CreateSkill(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, skill1)

	skill2, err := testQueries.CreateSkill(context.Background(), args)
	require.Error(t, err)
	require.Empty(t, skill2)

	err = testQueries.DeleteSkill(context.Background(), skill1.ID)
	require.NoError(t, err)
	err = testQueries.DeleteSkill(context.Background(), skill2.ID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

func TestGetSkill(t *testing.T) {
	user := CreateRandomUser(t)
	skill1 := CreateRandomSkill(t, user)
	skill2, err := testQueries.GetSkill(context.Background(), skill1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, skill2)

	require.Equal(t, skill1.ID, skill2.ID)
	require.Equal(t, skill1.SkillName, skill2.SkillName)
	require.Equal(t, skill1.SkillLevel, skill2.SkillLevel)

	err = testQueries.DeleteSkill(context.Background(), skill1.ID)
	require.NoError(t, err)
	err = testQueries.DeleteSkill(context.Background(), skill2.ID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

func createMultipleSkills(t *testing.T, user User) {
	randomProgLangLength := len(randomProgLang)
	for i := 0; i < randomProgLangLength; i++ {
		args := CreateSkillParams{
			Owner:      user.Username,
			SkillName:  randomProgLang[i],
			SkillLevel: randomProfLvl[rand.Intn(len(randomProfLvl))],
		}

		skill, err := testQueries.CreateSkill(context.Background(), args)
		require.NoError(t, err)
		require.NotEmpty(t, skill)

		require.Equal(t, args.SkillName, skill.SkillName)
		require.Equal(t, args.SkillLevel, skill.SkillLevel)

		require.NotZero(t, skill.ID)
	}
}

func TestListSkills(t *testing.T) {
	user := CreateRandomUser(t)
	createMultipleSkills(t, user)

	args := ListSkillsParams{
		Owner:  user.Username,
		Limit:  5,
		Offset: 1,
	}

	skills, err := testQueries.ListSkills(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, skills, 5)

	for _, skill := range skills {
		require.NotEmpty(t, skill)
	}
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

func TestUpdateSkill(t *testing.T) {
	user := CreateRandomUser(t)
	skill1 := CreateRandomSkill(t, user)

	rand.Seed(time.Now().UnixNano())
	args := UpdateSkillParams{
		ID:         skill1.ID,
		SkillName:  randomProgLang[rand.Intn(len(randomProgLang))],
		SkillLevel: randomProfLvl[rand.Intn(len(randomProfLvl))],
	}

	skill2, err := testQueries.UpdateSkill(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, skill2)

	require.Equal(t, args.ID, skill2.ID)
	require.Equal(t, args.SkillName, skill2.SkillName)
	require.Equal(t, args.SkillLevel, skill2.SkillLevel)

	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

}
