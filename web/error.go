package web

import (
	"html/template"
	"log"
	"net/http"
)

// Handle errors. Show error to user a log them.
type errorPage struct {
	Title   string
	Message string
}

func error502(w http.ResponseWriter, err error) {
	log.Println(err)

	title := "Error 502"
	message := "Error 502 : Internal server error"

	p := errorPage{title, message}                  // Construct page struct
	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func error404(w http.ResponseWriter) {
	title := "Error 404"
	message := "Error 404 : Page not found"

	p := errorPage{title, message}                  // Construct page struct
	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func mailFormatError(w http.ResponseWriter) {
	log.Println("Invalid email format")

	p := errorPage{"Email invalide", "Email invalide"}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func mailUniqueError(w http.ResponseWriter) {
	log.Println("Email not unique")

	p := errorPage{"Email invalide", "Email déjà utilisé"}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func voucherError(w http.ResponseWriter) {
	log.Println("Invalid voucher")

	p := errorPage{"Voucher invalide", "Voucher non valide"}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func simpleMessage(w http.ResponseWriter, msg string) {
	p := errorPage{"Debug message", msg}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func voucherExpired(w http.ResponseWriter) {
	log.Println("Voucher expired")

	p := errorPage{"Voucher invalide", "Ce voucher a expiré"}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func passwordError(w http.ResponseWriter) {
	log.Println("Connect attempt: incorrect password")

	p := errorPage{"Mot de passe invalide", "Le mot de passe ne correspond pas, veuillez réessayer."}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}

func mailError(w http.ResponseWriter) {
	log.Println("Connect attempt: incorrect mail")

	p := errorPage{"Email invalide", "Adresse mail introuvable, veuillez réessayer."}

	t, err := template.ParseFiles("html/error.hbs") // Load template
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, p) // Build and send page to user
}
