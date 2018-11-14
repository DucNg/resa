package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DucNg/resa/modele"
	"github.com/DucNg/resa/tools"
)

// Describe the admin structure page.
// This isn't a modele! It is only used to build the visual aspect of the page for the user (frontend).
type page struct {
	I                 modele.Invite
	ParrainMail       string
	VoucherCode       string
	VoucherExpiration string
	VoucherDisable    bool
}

// AdminIndex handle the /admin page and redirect the user.
// It shows the list of invite if the admin token is present and valid or the login page.
func AdminIndex(w http.ResponseWriter, r *http.Request) {
	validSession, err := verifySession(w, r)
	if validSession {
		AdminListInvite(w, r) // The token is valid show the administration page
		return
	}
	if err != nil {
		log.Println(err)
	}
	http.ServeFile(w, r, "html/adminLogin.html") // Invalid voucher redirect to login page
}

// verifySession verify the user session cookie.
// It gets the local cookie first and then check it's validity on the database
func verifySession(w http.ResponseWriter, r *http.Request) (bool, error) {
	sessionToken, err := getAdminCookie(r)
	if err != nil {
		return false, err
	}
	if sessionToken != "" {
		// Connect to database first
		db, err := tools.Connect()
		if err != nil {
			return false, err
		}
		defer tools.Disconnect(db)

		// Verify token validity
		validToken, err := tools.VerifyAdminSession(db, sessionToken) // Tell if the token is present in database
		if err != nil {
			return false, err
		}
		return validToken, err
	}
	return false, err // No valid token was found
}

// AdminConnect connect an admin using login and password.
// The connect process is to check validity of informations, if valid create a token then redirect to index see: AdminIndex()
func AdminConnect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Getting informations from POST

	user := modele.Admin{ // Get informations from form
		Login: r.FormValue("login"),
		Psw:   r.FormValue("mdp"),
	}

	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	err = tools.ConnectAdmin(db, &user) // Verify info and fill struct
	if err != nil {
		error502(w, err)
		return
	}

	// Create session
	user.Token, err = tools.CreateAdminSession(db, user.IdAdmin)
	if err != nil {
		error502(w, err) // Token creation error
		return
	}
	// We've created server side session but we still need to create the user cookie
	setAdminCookie(w, user.Token)

	// Redirect to admin page, will auto connect the useru using his token
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// AdminListInvite build a page of every Invite in database.
// It shows parrain for every user linking idParrain to corresponding email.
// It check if the invite has a voucher or not and show it's expiration date.
// It also check if the voucher is disable.
func AdminListInvite(w http.ResponseWriter, r *http.Request) {
	var listInvite []modele.Invite
	listInvite = make([]modele.Invite, 0) // Empty list of invite

	// Connect to database first
	db, err := tools.Connect()
	if err != nil {
		error502(w, err) // Show error to user and log it
		return
	}
	defer tools.Disconnect(db)

	err = tools.ListInvite(db, &listInvite)
	if err != nil {
		log.Println(err) // Error in the select won't be critical, don't need to inform user
	}

	// Getting all the vouchers
	var vouchers map[int64]modele.Voucher
	vouchers = make(map[int64]modele.Voucher)
	err = tools.GetVouchers(db, vouchers)
	if err != nil {
		log.Println(err) // Error in the select won't be critical, don't need to inform user
	}
	// Get the parrain email and the voucher if the user has one
	var p []page                         // Construct the page
	for _, element := range listInvite { // Iterate on each invite
		var tmpPage page
		tmpParrainmail := modele.GetParrainMail(element.Parrain, listInvite) // Get the corresponding parrain mail for every Invite

		if vouchers[element.Id].Code != "" {
			tmpPage = page{
				I:                 element,
				ParrainMail:       tmpParrainmail,
				VoucherCode:       vouchers[element.Id].Code,
				VoucherExpiration: vouchers[element.Id].Expiration.Format(time.RFC822),    // Get the expiration date as a string
				VoucherDisable:    vouchers[element.Id].Expiration.Equal(time.Unix(0, 0)), // Is voucher disable?
			}
		} else {
			tmpPage = page{
				I:           element,
				ParrainMail: tmpParrainmail,
				VoucherCode: vouchers[element.Id].Code,
			}
		}

		p = append(p, tmpPage) // List of invite with parrain
	}

	t, err := template.ParseFiles("html/adminListInvite.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(w, p) // Build and send page to user
	if err != nil {
		error502(w, err)
		return
	}
}

// AddVoucher is the controller to add a voucher.
// This func is used in to situations:
// * GET method: Provide the form page to enter informations on the voucher (code and expiration)
// * POST method: Insert the voucher in database using informations from the form
func AddVoucher(w http.ResponseWriter, r *http.Request) {
	validSession, err := verifySession(w, r)
	if err != nil {
		log.Println(err)
	}
	if !validSession { // This action is only available if connected as an admin
		http.Redirect(w, r, "/admin", http.StatusFound) // Invalid voucher redirect to login page
		return
	}
	if r.Method == "GET" { // Send the form to select parameters
		id, err := strconv.ParseInt(r.FormValue("id"), 10, 64) // Receive id_invite from GET
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}

		// Connect to database first
		db, err := tools.Connect()
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}
		defer tools.Disconnect(db)

		Invite, err := tools.GetInvite(db, id)
		if err != nil {
			error502(w, err)
			return
		}

		t, err := template.ParseFiles("html/addVoucher.hbs") // Load template
		if err != nil {
			log.Println(err)
		}

		err = t.Execute(w, Invite) // Build and send page to user
		if err != nil {
			error502(w, err)
			return
		}
	} else if r.Method == "POST" {
		r.ParseForm() // Getting informations from POST

		expiration, err := time.Parse("2006-01-02T15:04", r.FormValue("expiration"))
		log.Println(err)
		prop, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		log.Println(err)
		// TODO handle errors (voucher in the past)

		voucher := modele.Voucher{ // Fill the Invite struct with available informations
			Code:       r.FormValue("code"),
			Expiration: expiration,
			Prop:       prop,
		}

		// Connect to database first
		db, err := tools.Connect()
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}
		defer tools.Disconnect(db)

		err = tools.AddVoucher(db, voucher)
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}

		// Redirect to admin page
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		error404(w)
	}
}

// DisableVoucher using a user id
// Disable means set is expiration date to UNIX timestamp 0
// TODO show a confirmation page before disabling
func DisableVoucher(w http.ResponseWriter, r *http.Request) {
	validSession, err := verifySession(w, r)
	if err != nil {
		log.Println(err)
	}
	if !validSession { // This action is only available if connected as an admin
		http.Redirect(w, r, "/admin", http.StatusFound) // Invalid voucher redirect to login page
		return
	}
	if r.Method == "GET" {
		id, err := strconv.ParseInt(r.FormValue("id"), 10, 64) // Receive id_invite from GET
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}

		// Connect to database first
		db, err := tools.Connect()
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}
		defer tools.Disconnect(db)

		err = tools.DisableVoucher(db, id)
		if err != nil {
			error502(w, err) // Show error to user and log it
			return
		}

		// Redirect to admin page
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		error404(w)
	}
}
