package handlers

/*
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IsuruHaupe/web-api/types"
)

func TestCreateContact(t *testing.T) {
	contact := types.Contact{
		FirstName:   "Isuru",
		LastName:    "HAUPE",
		Fullname:    "Isuru HAUPE",
		Address:     "9 Allée Jean Baptiste Fourrier",
		Email:       "isuru.li@gmail.com",
		PhoneNumber: "0656759791",
	}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(contact)
	if err != nil {
		t.Fatal(err)
	}

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/post-endpoint", &b)
	CreateContact(wr, req)
	if wr.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", wr.Code)
	}
}

func TestCreateContactUnknownField(t *testing.T) {

	wrongContact := struct {
		WrongItem string
	}{
		WrongItem: "zoiefjzeio",
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(wrongContact)
	if err != nil {
		t.Fatal(err)
	}

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/post-endpoint", &b)

	CreateContact(wr, req)
	if wr.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", wr.Code)
	}

	if !strings.Contains(wr.Body.String(), "Request body contains unknown field") {
		t.Errorf(
			`response body "%s" does not contain "Request body contains unknown field"`,
			wr.Body.String(),
		)
	}
}

func TestCreateContactInvalidValue(t *testing.T) {
	wrongContact := struct {
		FirstName string
		LastName  string
		Fullname  string
		Adress    string
		Email     int
	}{
		FirstName: "Isuru",
		LastName:  "HAUPE",
		Fullname:  "Isuru HAUPE",
		Adress:    "9 Allée Jean Baptiste Fourrier",
		Email:     1234,
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(wrongContact)
	if err != nil {
		t.Fatal(err)
	}

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/post-endpoint", &b)

	CreateContact(wr, req)
	if wr.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", wr.Code)
	}

	if !strings.Contains(wr.Body.String(), "Request body contains an invalid value") {
		t.Errorf(
			`response body "%s" does not contain "Request body contains an invalid value"`,
			wr.Body.String(),
		)
	}
}

func TestCreateContactEmpty(t *testing.T) {

	var b bytes.Buffer
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/post-endpoint", &b)

	CreateContact(wr, req)
	if wr.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", wr.Code)
	}

	if !strings.Contains(wr.Body.String(), "Request body must not be empty") {
		t.Errorf(
			`response body "%s" does not contain "Request body must not be empty"`,
			wr.Body.String(),
		)
	}
}*/
