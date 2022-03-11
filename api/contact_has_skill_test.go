package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	mockdb "github.com/IsuruHaupe/web-api/db/mock"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateContactHasSkillAPI(t *testing.T) {
	user, _ := randomUser(t)
	skill := randomSkill(user.Username)
	contact := randomContact(user.Username)
	contactHasSkill := db.ContactHasSkill{
		Owner:     user.Username,
		ContactID: int32(contact.ID),
		SkillID:   int32(skill.ID),
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test PASS",
			body: gin.H{
				"contact_id": contact.ID,
				"skill_id":   skill.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateContactHasSkillParams{
					Owner:     contact.Owner,
					ContactID: int32(contact.ID),
					SkillID:   int32(skill.ID),
				}
				database.EXPECT().
					CreateContactHasSkill(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contactHasSkill, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContactHasSkill(t, recorder.Body, contactHasSkill)
			},
		},
		{
			name: "No Authorization",
			body: gin.H{
				"contact_id": contact.ID,
				"skill_id":   skill.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateContactHasSkillParams{
					Owner:     contact.Owner,
					ContactID: int32(contact.ID),
					SkillID:   int32(skill.ID),
				}
				database.EXPECT().
					CreateContactHasSkill(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Test INTERNAL ERROR",
			body: gin.H{
				"contact_id": contact.ID,
				"skill_id":   skill.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContactHasSkill(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ContactHasSkill{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Test INVALID PARAM",
			body: gin.H{
				"contact_id": "",
				"skill_id":   skill.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContactHasSkill(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.ContactHasSkill{}, sql.ErrConnDone)
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

			server := newTestServer(t, database)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := "/add-skill"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchContactHasSkill(t *testing.T, body *bytes.Buffer, contactHasSkill db.ContactHasSkill) {
	res, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotContactHasSkill db.ContactHasSkill
	err = json.Unmarshal(res, &gotContactHasSkill)
	require.NoError(t, err)
	require.Equal(t, contactHasSkill, gotContactHasSkill)
}
