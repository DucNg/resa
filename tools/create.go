package tools

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"errors"
	"strings"

	"github.com/DucNg/resa/modele"
)

const request string = `
DROP TABLE Voucher;
DROP TABLE Session;
DROP TABLE AdminSession;
DROP TABLE Invite;
DROP TABLE Administrateur;

CREATE TABLE Invite (
	id_invite INTEGER PRIMARY KEY,
	nom TEXT,
	prenom TEXT,
	mail TEXT NOT NULL,
	mdp TEXT NOT NULL,
	numtel TEXT,
	parrain INTEGER REFERENCES id_invite
);

CREATE TABLE Voucher (
	id_voucher INTEGER PRIMARY KEY,
	code TEXT,
	expiration TIMESTAMP,
	proprietaire INTEGER,
	FOREIGN KEY (proprietaire) REFERENCES Invite(id_invite)
);

CREATE TABLE Administrateur (
	id_admin INTEGER PRIMARY KEY,
	login TEXT,
	mdp TEXT
);

CREATE TABLE Session (
	token TEXT NOT NULL PRIMARY KEY,
	id_user INTEGER REFERENCES Invite(id_invite)
);

CREATE TABLE AdminSession (
	token TEXT NOT NULL PRIMARY KEY,
	id_user INTEGER REFERENCES Invite(id_admin)
)
`

var slicedRequest []string = strings.Split(request, ";")

// InitDatabase use the script defined in the const to create the database.
// It keep going even if errors append, it concatenate errors and return everything at ones.
// It give errors on first run because it contrain DROP TABLE.
// It can be used to create the database the first time, update the structure and reset the database to it empty state.
func InitDatabase(db *sql.DB) error {
	var errs string

	for _, q := range slicedRequest {
		_, err := db.Exec(q)

		if err != nil {
			errs += err.Error() + "\n"
		}
	}
	if errs != "" {
		return errors.New(errs)
	}
	return nil
}

// CreateAdmin create the admin using a modele and return the inserted id.
// The password is hashed using the HashPassword func described in database.go.
func CreateAdmin(db *sql.DB, admin *modele.Admin) (int64, error) {
	var notHashedPsw string = admin.Psw
	var hashedPsw string

	hashedPsw, err := HashPassword(notHashedPsw) // Hashing the password before sending to database
	if err != nil {
		return -1, err
	}

	result, err := db.Exec("INSERT INTO Administrateur(login,mdp) VALUES (?,?)", admin.Login, hashedPsw)

	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}
