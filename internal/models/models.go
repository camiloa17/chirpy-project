package models

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}
