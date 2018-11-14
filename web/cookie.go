package web

import (
	"net/http"
	"time"
)

// setSessionCookie Add a session cookie to the user.
// Client side session. Default expiration is 30 minutes.
func setSessionCookie(w http.ResponseWriter, token string) {
	expiration := time.Now().Add(time.Minute * 30) // Cookie expiration is in 30 minutes (30 * 1 minute)

	mySession := http.Cookie{
		Name:  "session",
		Value: token, // Cookie value is the unique token generated in session.go and has the same value as the one in the database

		Expires:  expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &mySession)
}

// getSessionCookie get the user cookie named session.
// It should contrain the user session but verification isn't made here.
// Value is "" if nothing was found.
func getSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

// Delete the cookie named session.
// Delete session client side.
func deleteSessionCookie(w http.ResponseWriter) {
	expiration := time.Unix(0, 0) // Set the expiration to 01 Jan 1970 00:00:00

	mySession := http.Cookie{
		Name:  "session",
		Value: "", // Value of the cookie is now empty

		Expires:  expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &mySession)
}

// Same as setSessionCookie() for admin.
func setAdminCookie(w http.ResponseWriter, token string) {
	expiration := time.Now().Add(time.Minute * 30) // Cookie expiration is in 30 minutes (30 * 1 minute)

	mySession := http.Cookie{
		Name:  "admin",
		Value: token, // Cookie value is the unique token generated in session.go and has the same value as the one in the database

		Expires:  expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &mySession)
}

// Same as getSessionCookie() for admin.
func getAdminCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("admin")

	if err != nil {
		return "", err
	}
	return cookie.Value, err
}
