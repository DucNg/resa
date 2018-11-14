// Package main run the web server.
// It provide configuration using flags and ini files. Details in config package.
// Main shouldn't connect to the database.
package main

import (
	"bufio"
	"fmt"
	"github.com/vharitonsky/iniflags"
	"log"
	"net/http"
	"os"

	"github.com/DucNg/resa/config"
	"github.com/DucNg/resa/tools"
	"github.com/DucNg/resa/web"
)

// inita initialize database and insert needed values to start.
// It ask user to input informations from the command line.
func inita() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Création 1er admin\nLogin : ")
	reader.Scan()
	login := reader.Text()
	fmt.Print("Mot de passe : ")
	reader.Scan()
	mdp := reader.Text()
	/*	fmt.Print("1er voucher : ")
		reader.Scan()
		voucher := reader.Text()*/

	err := tools.CreateDatabase()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database created with sucess")

	err = tools.ParseAndCreateAdmin(login, mdp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Admin added with sucess")

	err = tools.CreateDefaultUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Default user added with sucess\nYou need to manualy add a voucher to this user via admin page.")
}

// main create the handle for every pages on the server. It links pages to related function.
func main() {
	/*if len(os.Args) > 1 {
		args := os.Args[1:] // Args without program
		cmdTools(args)
		return
	}*/
	iniflags.Parse() // Get the configuration

	if *config.Firstrun { // If the -init flag is set, initilize
		inita()
	}

	http.HandleFunc("/", web.Index)                        // Index and static files
	http.HandleFunc("/connect", web.Connect)               // Connection and user page
	http.HandleFunc("/register", web.Register)             // Handle the register page
	http.HandleFunc("/disconnect", web.Disconnect)         // Delete session
	http.HandleFunc("/admin", web.AdminIndex)              // Show admin page if cookie or login
	http.HandleFunc("/adminconnect", web.AdminConnect)     // Handle connect admin form
	http.HandleFunc("/addVoucher", web.AddVoucher)         // Add voucher to an invite
	http.HandleFunc("/disableVoucher", web.DisableVoucher) // Disable a voucher to an invite

	fmt.Println("Listening on " + *config.Port)
	http.ListenAndServe(":"+*config.Port, nil)
}
