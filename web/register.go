package web

import (
	"net/http"

	"github.com/DucNg/resa/modele"
	"github.com/DucNg/resa/tools"
)

// Register get informations from a form, verify these informations and build a modele using them.
// It inserts informations into the database.
// It makes the association between invite and parrain.
// It create the user session (client side and server side).
// It redirect user to index (he will be automatically connected using the token)
func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Getting informations from POST

	user := modele.Invite{ // Create the Invite struct
		Nom:    r.FormValue("nom"),
		Prenom: r.FormValue("prenom"),
		Mail:   r.FormValue("mail"),
		Mdp:    r.FormValue("mdp"),
		Numtel: r.FormValue("numtel"),
	}

	// Check email format
	matched, err := modele.CheckMail(user.Mail)
	if err != nil { // Regexp error
		error502(w, err) // Show error to user and log it
		return
	}
	if !matched { // Format didn't match
		mailFormatError(w)
		return
	}

	// Check if email is unique in database

	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	// Check unique email then
	isUnique, err := tools.UniqueMail(db, user.Mail)
	if err != nil { // Database error
		error502(w, err) // Show error to user and log it
		return
	}
	if !isUnique { // Format didn't match
		mailUniqueError(w)
		return
	}

	// Check voucher validity and make association
	voucher := r.FormValue("voucher")
	idParrain, err := tools.GetParrain(db, voucher)
	if err != nil { // Database error
		if err.Error() == "Voucher expired" {
			voucherExpired(w)
		} else if err.Error() == "Voucher doesn't exist" {
			voucherError(w) // Show error to user
		} else {
			error502(w, err) // Show error to user and log it
		}
		return
	}
	if idParrain == -1 { // If parrain is -1 mean no voucher was found
		voucherError(w) // Show error to user
		return
	}

	user.Parrain = idParrain // User now has a parrain

	// If everything is valid, writting informations to database and get the user id
	userId, err := tools.CreateUser(db, &user) // userId will be used when session will be implemented
	//_,err = tools.CreateUser(db,&user)
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}

	// Create session
	token, err := tools.CreateSession(db, userId)
	if err != nil { // Error generating token
		error502(w, err) // Show error to user and log it
		return
	}
	// We've created server side session but we still need to create the user cookie
	setSessionCookie(w, token)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusFound)
}
