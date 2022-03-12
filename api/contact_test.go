package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	mockdb "github.com/IsuruHaupe/web-api/db/mock"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetContactAPI(t *testing.T) {
	user, _ := randomUser(t)
	contact := randomContact(user.Username)

	testCases := []struct {
		name          string
		contactID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Pass",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
			},
		},
		{
			name:      "Unauthorized User",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "No Authorization",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(db.Contact{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			contactID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Any()).
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
			server := newTestServer(t, database)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/contacts/%d", currentTest.contactID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			// Check results.
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestGetContactsWithSkillAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 5
	contacts := make([]db.Contact, n)
	for i := 0; i < n; i++ {
		contacts[i] = randomContact(user.Username)
	}

	type Query struct {
		skillName string
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Pass",
			query: Query{
				skillName: "Go",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {

				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("Go")).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContacts(t, recorder.Body, contacts)
			},
		},
		{
			name: "Unauthorized User",
			query: Query{
				skillName: "Go",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("Go")).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			query: Query{
				skillName: "Go",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("Go")).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			query: Query{
				skillName: "10",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("10")).
					Times(1).
					Return([]db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			query: Query{
				skillName: "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("")).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			query: Query{
				skillName: "Scala",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContactsWithSkill(gomock.Any(), gomock.Eq("Scala")).
					Times(1).
					Return([]db.Contact{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := "/contacts-with-skill"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("skill_name", currentTest.query.skillName)
			request.URL.RawQuery = q.Encode()

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestGetContactsWithSkillAndLevelAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 5
	contacts := make([]db.Contact, n)
	for i := 0; i < n; i++ {
		contacts[i] = randomContact(user.Username)
	}

	type Query struct {
		skillName  string
		skillLevel string
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Pass",
			query: Query{
				skillName:  "Go",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "Go",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchContacts(t, recorder.Body, contacts)
			},
		},
		{
			name: "Unauthorized User",
			query: Query{
				skillName:  "Go",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "Go",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			query: Query{
				skillName:  "Go",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "Go",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			query: Query{
				skillName:  "10",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "10",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			query: Query{
				skillName:  "",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			query: Query{
				skillName:  "Go",
				skillLevel: "Proficient",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.GetContactsWithSkillAndLevelParams{
					SkillName:  "Go",
					SkillLevel: "Proficient",
				}
				database.EXPECT().
					GetContactsWithSkillAndLevel(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Contact{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := "/contacts-with-skill-and-level"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("skill_name", currentTest.query.skillName)
			q.Add("skill_level", currentTest.query.skillLevel)
			request.URL.RawQuery = q.Encode()

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestCreateContactAPI(t *testing.T) {
	user, _ := randomUser(t)
	contact := randomContact(user.Username)
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Pass",
			body: gin.H{
				"firstname":    contact.Firstname,
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateContactParams{
					Owner:       contact.Owner,
					Firstname:   contact.Firstname,
					Lastname:    contact.Lastname,
					Fullname:    contact.Fullname,
					HomeAddress: contact.HomeAddress,
					Email:       contact.Email,
					PhoneNumber: contact.PhoneNumber,
				}

				database.EXPECT().
					CreateContact(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contact, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
			},
		},
		{
			name: "No Authorization",
			body: gin.H{
				"firstname":    contact.Firstname,
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateContactParams{
					Owner:       contact.Owner,
					Firstname:   contact.Firstname,
					Lastname:    contact.Lastname,
					Fullname:    contact.Fullname,
					HomeAddress: contact.HomeAddress,
					Email:       contact.Email,
					PhoneNumber: contact.PhoneNumber,
				}

				database.EXPECT().
					CreateContact(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"firstname":    contact.Firstname,
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"firstname":    "",
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContact(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Contact{}, sql.ErrConnDone)
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

			url := "/contacts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestListContactsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	contacts := make([]db.Contact, n)
	for i := 0; i < n; i++ {
		contacts[i] = randomContact(user.Username)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Pass",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.ListContactsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}

				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContacts(t, recorder.Body, contacts)
			},
		},
		{
			name: "No Authorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Any()).
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

			server := newTestServer(t, database)
			recorder := httptest.NewRecorder()

			url := "/contacts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", currentTest.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", currentTest.query.pageSize))
			request.URL.RawQuery = q.Encode()

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestDeleteContactAPI(t *testing.T) {
	user, _ := randomUser(t)
	contact := randomContact(user.Username)

	testCases := []struct {
		name          string
		contactID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Pass",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					DeleteContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Unauthorized User",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					DeleteContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "No Authorization",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)

				database.EXPECT().
					DeleteContact(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			contactID: contact.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					DeleteContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			contactID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)

				database.EXPECT().
					DeleteContact(gomock.Any(), gomock.Any()).
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

			server := newTestServer(t, database)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/contacts/%d", currentTest.contactID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestUpdateContactAPI(t *testing.T) {
	user, _ := randomUser(t)
	contact := randomContact(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker auth.Maker)
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Pass",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					GetIfExistsContactID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
			},
		},
		{
			name: "Unauthorized User",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(0)

				database.EXPECT().
					GetIfExistsContactID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)

				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(0)

				database.EXPECT().
					GetIfExistsContactID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"id":        0,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrNoRows)
				database.EXPECT().
					GetIfExistsContactID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker auth.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
				database.EXPECT().
					GetIfExistsContactID(gomock.Any(), gomock.Eq(contact.ID)).
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

			server := newTestServer(t, database)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := "/contacts"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			currentTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func randomContact(owner string) db.Contact {
	return db.Contact{
		ID:          int64(randomdata.Number(1, 20)),
		Owner:       owner,
		Firstname:   randomdata.FirstName(randomdata.Female),
		Lastname:    randomdata.LastName(),
		Fullname:    randomdata.FullName(randomdata.Female),
		HomeAddress: randomdata.Address(),
		Email:       randomdata.Email(),
		PhoneNumber: randomdata.PhoneNumber(),
	}
}

func requireBodyMatchContact(t *testing.T, body *bytes.Buffer, contact db.Contact) {
	res, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotContact db.Contact
	err = json.Unmarshal(res, &gotContact)
	require.NoError(t, err)
	require.Equal(t, contact, gotContact)
}

func requireBodyMatchContacts(t *testing.T, body *bytes.Buffer, accounts []db.Contact) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Contact
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
