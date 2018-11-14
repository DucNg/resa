package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/DucNg/resa/modele"
	"github.com/DucNg/resa/tools"
)

// Connect using mail and password
// Get informations from the connection form on index page.
// Verify informations (show error), create session, redirect to /
func Connect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Getting informations from POST

	user := modele.Invite{ // Fill the Invite struct with available informations
		Mail: r.FormValue("mail"),
		Mdp:  r.FormValue("mdp"),
	}

	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	err = tools.ConnectUser(db, &user) // Fetch informations from database by reference
	if err != nil {                    // mangage different type of errors
		if err.Error() == "Connect: Incorrect password" {
			passwordError(w)
		} else if err.Error() == "Connect: Incorrect mail" {
			mailError(w)
		} else {
			error502(w, err) // Show error to user and log it
		}
		return
	}

	// Create session
	token, err := tools.CreateSession(db, user.Id)
	if err != nil { // Error generating token
		error502(w, err) // Show error to user and log it
		return
	}
	// We've created server side session but we still need to create the user cookie
	setSessionCookie(w, token)

	// Build and show user page
	/*showUserPage(w,user)*/

	// Session created, user can now connect using his token: ConnectToken()

	// Redirect to home page, will auto connect the useru using his token
	http.Redirect(w, r, "/", http.StatusFound)
}

// Connect using session token
// This func is called every time you go to / to check if you have a valid session.
// Session token (client side) needs to be pass as a parameter. It will be verified.
func ConnectToken(w http.ResponseWriter, r *http.Request, token string) {
	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	user, err := tools.VerifySession(db, token)
	if err != nil {
		if err.Error() == "Verify session: No user found" { // no user was found for this token delete the user's cookie
			// Cookie could have expired server side, user could have created a "custom" cookie, need to delete it
			deleteSessionCookie(w)
			// Redirect to home page
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			error502(w, err) // Show error to user and log it
		}
		return
	}

	// Everything went ok, show user page
	showUserPage(w, user)
}

// Build and show user page. Use invite modele to fill the informations on the page.
func showUserPage(w http.ResponseWriter, user modele.Invite) {
	// Getting the user's voucher if exist
	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	// Getting all the vouchers
	var vouchers map[int64]modele.Voucher
	vouchers = make(map[int64]modele.Voucher)
	err = tools.GetVouchers(db, vouchers)

	// Select the user's voucher
	user.Voucher = vouchers[user.Id].Code

	t, err := template.ParseFiles("html/userpage.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, user) // Build and send page to user
}

// Disconnect the user. Delete the session token, client side and server side.
func Disconnect(w http.ResponseWriter, r *http.Request) {
	token, err := getSessionCookie(r) // Get the user session token to delete
	if err != nil {
		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusFound)
		return // It's okay, cookie has probably already been deleted
	}

	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	err = tools.DeleteSession(db, token) // Delete server side session
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}

	deleteSessionCookie(w) // Delete user side session

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusFound)
}
