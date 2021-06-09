package models

type User struct {
	ID         int
	Email      string
	Password   string
	Firstname  string
	Lastname   string
	Identity   int  //changed from string
	Isdisabled bool //changed from int
}
