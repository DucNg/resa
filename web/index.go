// Package web is a controller. It links user interface to database and using modele.
// It handle all the web process except running the server.
// Create cookie, build pages, get informations from forms and verify them.
package web

import (
	"net/http"
)

// Index handle the / page. It redirect user to the login page if he has a valid token, serve the login page and all static files.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // If not requesting / serve static files in html folder
		http.ServeFile(w, r, "html/"+r.URL.Path)
		return
	}

	// Verify session cookie
	sessionToken, err := getSessionCookie(r)
	if err != nil { // Error means no cookie was found
		//error502(w,err)
		http.ServeFile(w, r, "html/index.html") // No cookie, need to connect or register
		return
	}
	// Session cookie is present, need to verify and create Invite.
	// This is a connection handled in user.go
	ConnectToken(w, r, sessionToken)
}
