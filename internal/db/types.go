package db

import "sync"

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RevokedTokens map[string]RevokedToken `json:"revoked_tokens"`
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}
type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
type RevokedToken struct {
	Time  string `json:"time"`
	Token string `json:"token"`
}
