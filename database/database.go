package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IsuruHaupe/web-api/types"
	"github.com/go-sql-driver/mysql"
)

type Database struct {
	connection *sql.DB
}

func (db *Database) Connect() {
	var err error
	// Get a database handle.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "web_api_database",
	}
	db.connection, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	// for docker env
	//db, err = sql.Open("mysql", "root:mypassword@tcp(db:3306)/history_of_message")
	if err != nil {
		log.Panic(err)
	}

	// MySQL server isn't fully active yet.
	// Block until connection is accepted. This is a docker problem with v3 & container doesn't start
	// up in time.
	for db.connection.Ping() != nil {
		fmt.Println("Attempting connection to db")
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Connected !")
}

func (db *Database) Close() {
	fmt.Println("Closing connection")
	err := db.connection.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (db *Database) CreateContactInDatabase(c types.Contact) error {
	_, err := db.connection.Exec("INSERT INTO person (firstname, lastname, fullname, home_address, email, phone_number) VALUES (?,?,?,?,?,?);", c.FirstName, c.LastName, c.Fullname, c.Address, c.Email, c.PhoneNumber)
	if err != nil {
		return fmt.Errorf("error when adding a new contact: %v", err)
	}
	return nil
}

func (db *Database) RemoveContactInDatabase(contactId int) error {
	_, err := db.connection.Exec("DELETE FROM person WHERE id = ?;", contactId)
	if err != nil {
		return fmt.Errorf("error when removing contact: %v", err)
	}
	return nil
}

// TODO: trouver un moyen de n'ajouter que les elements non vides
func (db *Database) UpdateContactInDatabase(contactId int, c types.Contact) error {
	_, err := db.connection.Exec("UPDATE person SET firstname = ?, lastname = ?, fullname = ?, home_address=?, email=?, phone_number=? WHERE id = ?;", c.FirstName, c.LastName, c.Fullname, c.Address, c.Email, c.PhoneNumber, contactId)
	if err != nil {
		return fmt.Errorf("error when updating contact: %v", err)
	}
	return nil
}

func (db *Database) GetContactInDatabase(contactId int) (types.Contact, error) {
	rows, err := db.connection.Query("SELECT * from person WHERE id = ?;", contactId)
	if err != nil {
		return types.Contact{}, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()
	var contact types.Contact
	for rows.Next() {
		if err := rows.Scan(&contact.Id, &contact.FirstName, &contact.LastName, &contact.Fullname, &contact.Address, &contact.Email, &contact.PhoneNumber); err != nil {
			return types.Contact{}, fmt.Errorf("error : %v", err)
		}
	}
	return contact, nil
}

func SetBDDEnvironmentVariable() {
	os.Setenv("DBUSER", "root")
	os.Setenv("DBPASS", "rootroot")
}

//ENUM('Familiar', 'Proficient', 'Excellent', 'Expert')
func (db *Database) CreateSkillInDatabase(s types.Skill) error {
	_, err := db.connection.Exec("INSERT INTO skills (skill_name, skill_level) VALUES (?,?);", s.Name, s.Level)
	if err != nil {
		return fmt.Errorf("error when adding a new skill: %v", err)
	}
	return nil
}

func (db *Database) RemoveSkillInDatabase(skillId int) error {
	_, err := db.connection.Exec("DELETE FROM skills WHERE id = ?;", skillId)
	if err != nil {
		return fmt.Errorf("error when removing skill: %v", err)
	}
	return nil
}

// TODO: trouver un moyen de n'ajouter que les elements non vides
func (db *Database) UpdateSkillInDatabase(skillId int, s types.Skill) error {
	_, err := db.connection.Exec("UPDATE skills SET skill_name = ?, skill_level = ?;", s.Name, s.Level, skillId)
	if err != nil {
		return fmt.Errorf("error when updating contact: %v", err)
	}
	return nil
}

func (db *Database) GetSkillInDatabase(contactId int) (types.Skill, error) {
	rows, err := db.connection.Query("SELECT * from skills WHERE id = ?;", contactId)
	if err != nil {
		return types.Skill{}, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()
	var skill types.Skill
	for rows.Next() {
		if err := rows.Scan(&skill.Id, &skill.Name, &skill.Level); err != nil {
			return types.Skill{}, fmt.Errorf("error : %v", err)
		}
	}
	return skill, nil
}
