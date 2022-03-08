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

func CreateRandomSkill(t *testing.T) Skill {
	rand.Seed(time.Now().UnixNano())
	args := CreateSkillParams{
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
	skill := CreateRandomSkill(t)
	err := testQueries.DeleteSkill(context.Background(), skill.ID)
	require.NoError(t, err)
}
func TestDeleteSkill(t *testing.T) {
	skill1 := CreateRandomSkill(t)
	err := testQueries.DeleteSkill(context.Background(), skill1.ID)
	require.NoError(t, err)

	skill2, err := testQueries.GetSkill(context.Background(), skill1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, skill2)
}

func TestUniqueConstraintOnSkill(t *testing.T) {
	args := CreateSkillParams{
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
}

func TestGetSkill(t *testing.T) {
	skill1 := CreateRandomSkill(t)
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
}

func createMultipleSkills(t *testing.T) []int64 {
	randomProgLangLength := len(randomProgLang)
	skillsID := make([]int64, randomProgLangLength)
	for i := 0; i < randomProgLangLength; i++ {
		args := CreateSkillParams{
			SkillName:  randomProgLang[i],
			SkillLevel: randomProfLvl[rand.Intn(len(randomProfLvl))],
		}

		skill, err := testQueries.CreateSkill(context.Background(), args)
		require.NoError(t, err)
		require.NotEmpty(t, skill)

		require.Equal(t, args.SkillName, skill.SkillName)
		require.Equal(t, args.SkillLevel, skill.SkillLevel)

		require.NotZero(t, skill.ID)
		skillsID[i] = skill.ID
	}

	return skillsID
}

func removeMultipleSkills(t *testing.T, skillsID []int64) {
	for i := 0; i < len(skillsID); i++ {
		err := testQueries.DeleteSkill(context.Background(), skillsID[i])
		require.NoError(t, err)
	}
}

func TestListSkills(t *testing.T) {
	skillsID := createMultipleSkills(t)

	args := ListSkillsParams{
		Limit:  5,
		Offset: 1,
	}

	skills, err := testQueries.ListSkills(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, skills, 5)

	for _, skill := range skills {
		require.NotEmpty(t, skill)
	}
	removeMultipleSkills(t, skillsID)
}

func TestUpdateSkill(t *testing.T) {
	skill1 := CreateRandomSkill(t)

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

	err = testQueries.DeleteSkill(context.Background(), skill2.ID)
	require.NoError(t, err)

}
