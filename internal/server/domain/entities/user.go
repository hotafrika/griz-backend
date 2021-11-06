package entities

type User struct {
	ID       uint64
	Username string
	Password string //Always empty or encrypted
	Email    string
}
