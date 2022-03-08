package handlers

/*
import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/IsuruHaupe/web-api/database"
	"github.com/IsuruHaupe/web-api/types"
)

type SkillHandler struct {
	db database.Database
}

func NewSkillHandler(db database.Database) *SkillHandler {
	return &SkillHandler{
		db: db,
	}
}

func (ch *ContactHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var s types.Skill
	switch r.Method {
	case "POST":
		err := decodeJSONBody(w, r, &s)
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
		err = ch.db.CreateSkillInDatabase(s)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%+v", s)

	default:
		fmt.Fprintf(w, "Unrecognised Query type, expected POST request !")
		log.Printf("Unrecognised Query type, expected POST request !")
	}
}
*/
