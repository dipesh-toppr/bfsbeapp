package models

type User struct {
	ID         int
	Email      string
	Password   string
	Firstname  string
	Lastname   string
	Identity   int
	Isdisabled bool
}

var IDENTITY = map[string]string{
	"teacher": "0",
	"student": "1",
}
