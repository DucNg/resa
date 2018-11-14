package tools

import (
	"time"

	"github.com/DucNg/resa/modele"
)

// CreateDatabase connect to database and create empty database using script. It's a controller.
func CreateDatabase() error {
	// Connect to database first
	db, err := Connect()
	if err != nil {
		return err
	}
	defer Disconnect(db)

	err = InitDatabase(db)
	return err
}

// ParseAndCreateAdmin create the admin. It was created to avoid cycling dependencies in main.
func ParseAndCreateAdmin(login string, psw string) error {
	// Connect to database first
	db, err := Connect()
	if err != nil {
		return err
	}
	defer Disconnect(db)

	user := modele.Admin{
		Login: login,
		Psw:   psw,
	}

	user.IdAdmin, err = CreateAdmin(db, &user)
	return err
}

// CreateVoucher **Unused** create a default voucher
func CreateVoucher(code string) error {
	// Connect to database first
	db, err := Connect()
	if err != nil {
		return err
	}
	defer Disconnect(db)

	myvoucher := modele.Voucher{
		Code:       code,
		Expiration: time.Now().Add(time.Hour * 24), // Default beaviour, the first voucher expire in 24 hours
		Prop:       -2,
	}

	err = AddVoucher(db, myvoucher)
	return err
}

// CreateDefaultUser create the default user, needed to add the first voucher
func CreateDefaultUser() error {
	// Connect to database first
	db, err := Connect()
	if err != nil {
		return err
	}
	defer Disconnect(db)

	user := modele.Invite{
		Nom:     "admin",
		Prenom:  "admin",
		Mail:    "contact@resa.com",
		Numtel:  "",
		Mdp:     "root",
		Parrain: -2,
	}

	hashedPsw, err := HashPassword(user.Mdp) // Hashing the password before sending to database
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO Invite(nom,prenom,mail,mdp,numtel,parrain) VALUES(?,?,?,?,?,?)", user.Nom, user.Prenom, user.Mail, hashedPsw, user.Numtel, user.Parrain)
	return err
}
