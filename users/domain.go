package users

type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"username"`
	Age       int
	Gender    string
	FirstName string
	LastName  string
	Email     string
	Password  string
}
