// forms.go
package main

/*
import (
	"fmt"
	"log"
	"net/http"

	"github.com/IsuruHaupe/web-api/database"
	"github.com/IsuruHaupe/web-api/handlers"
)

func main() {
	database.SetBDDEnvironmentVariable()
	var db = new(database.Database)
	db.Connect()
	defer db.Close()

	//init handlers
	var ch = handlers.NewContactHandler(*db)

	http.HandleFunc("/contact-api/create", ch.CreateContact)
	http.HandleFunc("/contact-api/remove", ch.RemoveContact)
	http.HandleFunc("/contact-api/update", ch.UpdateContact)
	http.HandleFunc("/contact-api/get", ch.GetContact)
	http.HandleFunc("/", HelloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
*/
