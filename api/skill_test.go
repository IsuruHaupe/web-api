package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/IsuruHaupe/web-api/db/mock"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetSkillAPI(t *testing.T) {
	skill := randomSkill()

	testCases := []struct {
		name          string
		skillID       int64
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "Test PASS",
			skillID: skill.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetSkill(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1).
					Return(skill, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSkill(t, recorder.Body, skill)
			},
		},
		{
			name:    "Test SKILL NOT FOUND",
			skillID: skill.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetSkill(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1).
					Return(db.Skill{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "Test INTERNAL ERROR",
			skillID: skill.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetSkill(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1).
					Return(db.Skill{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "Test INVALID PARAM",
			skillID: 0,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetSkill(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]
		t.Run(currentTest.name, func(t *testing.T) {
			// Init mock.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mockdb.NewMockDatabase(ctrl)

			// Create a stub.
			currentTest.buildStubs(database)

			// Start server and tests.
			server := NewServer(database)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/skills/%d", currentTest.skillID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// Check results.
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestCreateSkillAPI(t *testing.T) {
	skill := randomSkill()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test PASS",
			body: gin.H{
				"skill_name":  skill.SkillName,
				"skill_level": skill.SkillLevel,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateSkillParams{
					SkillName:  skill.SkillName,
					SkillLevel: skill.SkillLevel,
				}

				database.EXPECT().
					CreateSkill(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(skill, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSkill(t, recorder.Body, skill)
			},
		},
		{
			name: "Test INTERNAL ERROR",
			body: gin.H{
				"skill_name":  skill.SkillName,
				"skill_level": skill.SkillLevel,
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Skill{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Test INVALID PARAM",
			body: gin.H{
				"skill_name":  "",
				"skill_level": skill.SkillLevel,
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateSkill(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Skill{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mockdb.NewMockDatabase(ctrl)
			currentTest.buildStubs(database)

			server := NewServer(database)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := "/skills"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestListSkillsAPI(t *testing.T) {
	n := 5
	skills := make([]db.Skill, n)
	for i := 0; i < n; i++ {
		skills[i] = randomSkill()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.ListSkillsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				database.EXPECT().
					ListSkills(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(skills, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSkills(t, recorder.Body, skills)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListSkills(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Skill{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListSkills(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListSkills(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mockdb.NewMockDatabase(ctrl)
			currentTest.buildStubs(database)

			server := NewServer(database)
			recorder := httptest.NewRecorder()

			url := "/skills"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", currentTest.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", currentTest.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestDeleteSkillAPI(t *testing.T) {
	skill := randomSkill()

	testCases := []struct {
		name          string
		skillID       int64
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			skillID: skill.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					DeleteSkill(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			skillID: skill.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					DeleteSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "BadRequest",
			skillID: 0,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					DeleteSkill(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mockdb.NewMockDatabase(ctrl)
			currentTest.buildStubs(database)

			server := NewServer(database)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/skills/%d", currentTest.skillID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSkillAPI(t *testing.T) {
	skill := randomSkill()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test PASS",
			body: gin.H{
				"id":          skill.ID,
				"skill_level": "Expert",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(skill, nil)

				database.EXPECT().
					GetIfExistsSkillID(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSkill(t, recorder.Body, skill)
			},
		},
		{
			name: "Test BAD REQUEST",
			body: gin.H{
				"id":        0,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateSkill(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Test NOT FOUND",
			body: gin.H{
				"id":        skill.ID,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Skill{}, sql.ErrNoRows)
				database.EXPECT().
					GetIfExistsSkillID(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Test INTERNAL SERVER ERROR",
			body: gin.H{
				"id":        skill.ID,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Skill{}, sql.ErrConnDone)
				database.EXPECT().
					GetIfExistsSkillID(gomock.Any(), gomock.Eq(skill.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mockdb.NewMockDatabase(ctrl)
			currentTest.buildStubs(database)

			server := NewServer(database)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := "/skills"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

var randomProgLang = [...]string{"Go", "Java", "Javascript", "C++", "Python", "R", "HTML"}
var randomProfLvl = [...]string{"Familiar", "Proficient", "Excellent", "Expert"}

func randomSkill() db.Skill {
	rand.Seed(time.Now().UnixNano())
	return db.Skill{
		ID:         int64(randomdata.Number(20)),
		SkillName:  randomProgLang[rand.Intn(len(randomProgLang))],
		SkillLevel: randomProfLvl[rand.Intn(len(randomProfLvl))],
	}
}

func requireBodyMatchSkill(t *testing.T, body *bytes.Buffer, skill db.Skill) {
	res, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotSkill db.Skill
	err = json.Unmarshal(res, &gotSkill)
	require.NoError(t, err)
	require.Equal(t, skill, gotSkill)
}

func requireBodyMatchSkills(t *testing.T, body *bytes.Buffer, skills []db.Skill) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotSkills []db.Skill
	err = json.Unmarshal(data, &gotSkills)
	require.NoError(t, err)
	require.Equal(t, skills, gotSkills)
}
