package types

type Contact struct {
	Id          int
	FirstName   string
	LastName    string
	Fullname    string
	Address     string
	Email       string
	PhoneNumber string
	//Skills      []Skills
}

type Skill struct {
	Id    int
	Name  string
	Level string // use enum -> const
}
