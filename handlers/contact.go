package handlers

/*
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IsuruHaupe/web-api/database"
	"github.com/IsuruHaupe/web-api/types"
)

type ContactHandler struct {
	db database.Database
}

func NewContactHandler(db database.Database) *ContactHandler {
	return &ContactHandler{
		db: db,
	}
}

func (ch *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//rename var
	var c types.Contact
	switch r.Method {
	case "POST":
		err := decodeJSONBody(w, r, &c)
		if err != nil {
			var mr *malformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.msg, mr.status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		// add contact in database
		err = ch.db.CreateContactInDatabase(c)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%+v", c)

	default:
		fmt.Fprintf(w, "Unrecognised Query type, expected POST request !")
		log.Printf("Unrecognised Query type, expected POST request !")
	}
}

func (ch *ContactHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case "DELETE":
		contactId, err := strconv.Atoi(r.URL.Query().Get("contact_id"))
		if err != nil {
			fmt.Fprintf(w, "The parameter is not a valid !")
			log.Printf("The parameter is not a valid ")
			return
		}
		// add contact in database
		err = ch.db.RemoveContactInDatabase(contactId)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			log.Fatal(err)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	default:
		fmt.Fprintf(w, "Unrecognised Query type, expected DELETE request !")
		log.Printf("Unrecognised Query type, expected DELETE request !")
	}
}

func (ch *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var c types.Contact
	switch r.Method {
	case "PUT":
		contactId, err := strconv.Atoi(r.URL.Query().Get("contact_id"))
		if err != nil {
			fmt.Fprintf(w, "The parameter is not a valid !")
			log.Printf("The parameter is not a valid ")
			return
		}
		err = decodeJSONBody(w, r, &c)
		if err != nil {
			var mr *malformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.msg, mr.status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		// add contact in database
		err = ch.db.UpdateContactInDatabase(contactId, c)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%+v", c)

	default:
		fmt.Fprintf(w, "Unrecognised Query type, expected POST request !")
		log.Printf("Unrecognised Query type, expected POST request !")
	}
}

func (ch *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var contact types.Contact
	switch r.Method {
	case "GET":
		contactId, err := strconv.Atoi(r.URL.Query().Get("contact_id"))
		if err != nil {
			fmt.Fprintf(w, "The parameter is not a valid !")
			log.Printf("The parameter is not a valid ")
			return
		}
		// get contact in database
		contact, err = ch.db.GetContactInDatabase(contactId)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%+v", contact)

	default:
		fmt.Fprintf(w, "Unrecognised Query type, expected GET request !")
		log.Printf("Unrecognised Query type, expected GET request !")
	}
}
*/
