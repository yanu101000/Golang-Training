package entity

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var Users = []User{
	{ID: 1, Name: "Emma", Email: "yanu10@gmail.com"},
	{ID: 2, Name: "Bruno", Email: "yanu11@gmail.com"},
	{ID: 3, Name: "Rick", Email: "yanu12@gmail.com"},
	{ID: 4, Name: "Lena", Email: "yanu13@gmail.com"},
}
