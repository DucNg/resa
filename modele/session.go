package modele

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Session associate a user id to a session.
// It use hash map so using this should be fast. Use it as much as possible and avoid using database instead.
type Session map[int64]string

// NewSession create a new empty session
func NewSession() Session {
	var mySession Session
	mySession = make(map[int64]string)

	return mySession
}

// From a user id get a session, id with a unique code
func CreateSession(currentSession Session, idUser int64) Session {
	c := 20
	b := make([]byte, c)

	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	randomString := base64.StdEncoding.EncodeToString(b)

	currentSession = Session{idUser: randomString}

	return currentSession
}
