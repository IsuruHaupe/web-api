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

	mockdb "github.com/IsuruHaupe/web-api/db/mock"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetContactAPI(t *testing.T) {
	contact := randomContact()

	testCases := []struct {
		name         string
		contactID    int64
		buildStubs   func(database *mockdb.MockDatabase)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Test PASS",
			contactID: contact.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(contact, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
			},
		},
		{
			name:      "Test CONTACT NOT FOUND",
			contactID: contact.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(db.Contact{}, sql.ErrNoRows)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Test INTERNAL ERROR",
			contactID: contact.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Test INVALID PARAM",
			contactID: 0,
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					GetContact(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			url := fmt.Sprintf("/contacts/%d", currentTest.contactID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// Check results.
			currentTest.checkReponse(t, recorder)
		})
	}
}

func TestCreateContactAPI(t *testing.T) {
	contact := randomContact()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test PASS",
			body: gin.H{
				"firstname":    contact.Firstname,
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.CreateContactParams{
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
			},
		},
		{
			name: "Test INTERNAL ERROR",
			body: gin.H{
				"firstname":    contact.Firstname,
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Test INVALID PARAM",
			body: gin.H{
				"firstname":    "",
				"lastname":     contact.Lastname,
				"fullname":     contact.Fullname,
				"home_address": contact.HomeAddress,
				"email":        contact.Email,
				"phone_number": contact.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockDatabase) {
				store.EXPECT().
					CreateContact(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := "/contacts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(recorder)
		})
	}
}

func TestListContactsAPI(t *testing.T) {
	n := 5
	contacts := make([]db.Contact, n)
	for i := 0; i < n; i++ {
		contacts[i] = randomContact()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				arg := db.ListContactsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				database.EXPECT().
					ListContacts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(contacts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContacts(t, recorder.Body, contacts)
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
					ListContacts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Contact{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
					ListContacts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
					ListContacts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := "/contacts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", currentTest.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", currentTest.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(recorder)
		})
	}
}

func TestDeleteContactAPI(t *testing.T) {
	contact := randomContact()

	testCases := []struct {
		name          string
		contactID     int64
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			contactID: contact.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
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
			name:      "InternalError",
			contactID: contact.ID,
			buildStubs: func(database *mockdb.MockDatabase) {
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
			name:      "BadRequest",
			contactID: 0,
			buildStubs: func(database *mockdb.MockDatabase) {
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

			server := NewServer(database)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/contacts/%d", currentTest.contactID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

// TODO : test Get fields functions.
func TestUpdateContactAPI(t *testing.T) {
	contact := randomContact()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(database *mockdb.MockDatabase)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test PASS",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(contact, nil)

				database.EXPECT().
					GetIfExistsID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchContact(t, recorder.Body, contact)
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
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Test NOT FOUND",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrNoRows)
				database.EXPECT().
					GetIfExistsID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Test INTERNAL SERVER ERROR",
			body: gin.H{
				"id":        contact.ID,
				"firstname": "Isuru",
			},
			buildStubs: func(database *mockdb.MockDatabase) {
				database.EXPECT().
					UpdateContact(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Contact{}, sql.ErrConnDone)
				database.EXPECT().
					GetIfExistsID(gomock.Any(), gomock.Eq(contact.ID)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := "/contacts"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(recorder)
		})
	}
}

func randomContact() db.Contact {
	return db.Contact{
		ID:          int64(randomdata.Number(20)),
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
